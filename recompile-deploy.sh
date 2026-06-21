#!/usr/bin/env bash
#
# recompile-deploy.sh — production recompile + restart on the DigitalOcean droplet.
#
# Run this ON the droplet (as root / via sudo) after pushing code. It:
#   1. Stops the service
#   2. Checkpoints + truncates the WAL (clean single-file DB)
#   3. Takes a full pre-deploy snapshot of the whole app dir  (../amazingtrak-backup)
#   4. Writes tiered rollback bundles — daily / weekly / monthly / yearly —
#      each a SINGLE FOLDER with everything needed to run that version:
#         amazingtrak.db  +  amazingtrak (binary)  +  .env  +  uploads/
#   5. Restarts the service (site back up), then git pull + rebuild
#   6. Restarts with the new binary and health-checks
#
# Tiers are refreshed by AGE, not calendar day, so they still roll forward
# correctly even if this script only runs occasionally:
#   daily   = every deploy
#   weekly  = when the existing one is >= 7 days old
#   monthly = >= 30 days old
#   yearly  = >= 365 days old
# Only ONE backup is kept per tier (each refresh overwrites that tier's folder).
#
# To roll back a version: stop the service, copy a tier folder's amazingtrak,
# amazingtrak.db, .env and uploads/ back into the app dir, start the service.
#
# Usage (on the droplet):
#   sudo ./recompile-deploy.sh                 # snapshot + tiers + pull + build + restart
#   sudo ./recompile-deploy.sh --no-pull       # skip git pull (build current tree)
#   sudo ./recompile-deploy.sh --no-snapshot   # skip the heavy full-dir snapshot

set -euo pipefail

# ── config ────────────────────────────────────────────────────────────────────
APP_USER="nv"
APP_DIR="/home/nv/amazingtrak"
SERVICE="amazingtrak"
GO_BIN="${GO_BIN:-/usr/local/go/bin/go}"
[ -x "$GO_BIN" ] || GO_BIN="$(command -v go || true)"

DB_FILE="${APP_DIR}/amazingtrak.db"
BACKUP_ROOT="/home/nv/backups/amazingtrak"   # tiered rollback bundles live here
FULL_SNAP="$(dirname "$APP_DIR")/$(basename "$APP_DIR")-backup"  # ../amazingtrak-backup

# Age thresholds (seconds) at which each tier is refreshed.
WEEKLY_AGE=$((7 * 86400))
MONTHLY_AGE=$((30 * 86400))
YEARLY_AGE=$((365 * 86400))

# ── args ──────────────────────────────────────────────────────────────────────
PULL=1
SNAPSHOT=1
for arg in "$@"; do
  case "$arg" in
    --no-pull)     PULL=0 ;;
    --no-snapshot) SNAPSHOT=0 ;;
    *) echo "unknown arg: $arg (use --no-pull / --no-snapshot)"; exit 1 ;;
  esac
done

# ── preflight ─────────────────────────────────────────────────────────────────
[ "$(id -u)" -eq 0 ] || { echo "run as root (sudo $0)"; exit 1; }
[ -d "$APP_DIR" ]    || { echo "missing $APP_DIR — run setup-vps.sh first"; exit 1; }
[ -x "$GO_BIN" ]     || { echo "Go not found at $GO_BIN"; exit 1; }
command -v sqlite3 >/dev/null || { echo "sqlite3 not installed (apt install sqlite3)"; exit 1; }

cd "$APP_DIR"

PORT="8000"
[ -f "${APP_DIR}/.env" ] && PORT="$(grep -E '^PORT=' "${APP_DIR}/.env" | cut -d= -f2 || echo 8000)"

# ── helpers ───────────────────────────────────────────────────────────────────
# older_than FILE SECONDS — true if FILE is missing or at least SECONDS old.
older_than() {
  local f="$1" max="$2"
  [ -f "$f" ] || return 0
  [ $(( $(date +%s) - $(stat -c %Y "$f") )) -ge "$max" ]
}

# write_bundle TIER — clean single-folder rollback bundle for that tier.
write_bundle() {
  local dir="${BACKUP_ROOT}/$1"
  echo "  ▸ ${1} bundle → ${dir}"
  mkdir -p "$dir"
  sqlite3 "$DB_FILE" ".backup '${dir}/amazingtrak.db'"        # consistent DB copy
  [ -f "${APP_DIR}/amazingtrak" ] && cp -p  "${APP_DIR}/amazingtrak" "${dir}/amazingtrak"  # pre-rebuild binary
  [ -f "${APP_DIR}/.env" ]        && cp -p  "${APP_DIR}/.env"        "${dir}/.env"
  [ -d "${APP_DIR}/uploads" ]     && rm -rf "${dir}/uploads" && cp -a "${APP_DIR}/uploads" "${dir}/uploads"
}

# ── stop service (everything below is copied while the DB is quiescent) ────────
echo "▶ stopping ${SERVICE}"
systemctl stop "$SERVICE"

# ── WAL checkpoint (flush + truncate so the DB is a clean single file) ────────
if [ -f "$DB_FILE" ]; then
  echo "▶ WAL checkpoint (TRUNCATE)"
  sqlite3 "$DB_FILE" "PRAGMA wal_checkpoint(TRUNCATE);"
fi

# ── full pre-deploy snapshot of the whole app dir (source, db, uploads, .env) ─
if [ "$SNAPSHOT" -eq 1 ]; then
  echo "▶ full snapshot → ${FULL_SNAP}"
  rm -rf "${FULL_SNAP}-old"
  [ -e "$FULL_SNAP" ] && mv "$FULL_SNAP" "${FULL_SNAP}-old"   # keep old until new copy succeeds
  cp -a "$APP_DIR" "$FULL_SNAP"
  rm -rf "${FULL_SNAP}-old"
fi

# ── tiered rollback bundles (db + binary + .env + uploads in one folder each) ─
if [ -f "$DB_FILE" ]; then
  echo "▶ tiered backups (one folder each: db + executable + .env + uploads)"
  write_bundle daily
  if older_than "${BACKUP_ROOT}/weekly/amazingtrak.db"  "$WEEKLY_AGE";  then write_bundle weekly;  fi
  if older_than "${BACKUP_ROOT}/monthly/amazingtrak.db" "$MONTHLY_AGE"; then write_bundle monthly; fi
  if older_than "${BACKUP_ROOT}/yearly/amazingtrak.db"  "$YEARLY_AGE";  then write_bundle yearly;  fi
fi

# ── restart service so the site is back up during the (slower) build ──────────
echo "▶ starting ${SERVICE}"
systemctl start "$SERVICE"

# ── git pull ──────────────────────────────────────────────────────────────────
if [ "$PULL" -eq 1 ]; then
  echo "▶ git pull"
  sudo -u "$APP_USER" git -C "$APP_DIR" pull --ff-only
fi

# ── build ─────────────────────────────────────────────────────────────────────
echo "▶ building binary"
sudo -u "$APP_USER" "$GO_BIN" build -o "${APP_DIR}/amazingtrak" ./...
echo "✓ built ${APP_DIR}/amazingtrak"

# ── restart with the new binary ───────────────────────────────────────────────
echo "▶ restarting ${SERVICE}"
systemctl restart "$SERVICE"

# ── health check ──────────────────────────────────────────────────────────────
for i in $(seq 1 20); do
  if curl -sf "http://localhost:${PORT}/healthz" >/dev/null 2>&1; then
    echo "✓ healthy: http://localhost:${PORT}/healthz"
    echo "  recent logs: journalctl -u ${SERVICE} -n 20 --no-pager"
    exit 0
  fi
  sleep 1
done

echo "✗ ${SERVICE} did not become healthy in 20s" >&2
echo "  check: journalctl -u ${SERVICE} -n 50 --no-pager" >&2
exit 1
