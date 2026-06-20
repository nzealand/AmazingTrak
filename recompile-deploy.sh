#!/usr/bin/env bash
#
# recompile-deploy.sh — production recompile + restart on the DigitalOcean droplet.
#
# Run this ON the droplet (as root / via sudo) after pushing code. It pulls the
# latest code, backs up the live SQLite DB, rebuilds the binary as the app user,
# and restarts the systemd service. DB schema migrations (new nullable columns /
# CREATE TABLE IF NOT EXISTS) apply automatically on startup — no manual SQL.
#
# Usage (on the droplet):
#   sudo ./recompile-deploy.sh            # git pull + build + restart
#   sudo ./recompile-deploy.sh --no-pull  # skip git pull (build current tree)

set -euo pipefail

# Layout — matches setup-vps.sh.
APP_USER="nv"
APP_DIR="/home/nv/amazingtrak"
DATA_DIR="/home/nv/amazingtrak"
SERVICE="amazingtrak"
GO_BIN="/usr/local/go/bin/go"
DB_FILE="${DATA_DIR}/amazingtrak.db"

PULL=1
[ "${1:-}" = "--no-pull" ] && PULL=0

[ "$(id -u)" -eq 0 ] || { echo "run as root (sudo $0)"; exit 1; }
[ -d "$APP_DIR" ]   || { echo "missing $APP_DIR — run setup-vps.sh first"; exit 1; }
[ -x "$GO_BIN" ]    || { echo "Go not found at $GO_BIN"; exit 1; }

cd "$APP_DIR"

# Read PORT from the app's .env for the post-deploy health check.
PORT="8000"
[ -f "${APP_DIR}/.env" ] && PORT="$(grep -E '^PORT=' "${APP_DIR}/.env" | cut -d= -f2 || echo 8000)"

if [ "$PULL" -eq 1 ]; then
  echo "▶ git pull"
  sudo -u "$APP_USER" git -C "$APP_DIR" pull --ff-only
fi

# Back up the live DB before restart (cheap insurance; migrations are non-destructive).
if [ -f "$DB_FILE" ]; then
  BAK="${DB_FILE}.bak-$(date +%Y%m%d-%H%M%S)"
  echo "▶ backing up DB → ${BAK}"
  cp -p "$DB_FILE" "$BAK"
  # Keep the 10 most recent backups.
  ls -1t "${DB_FILE}".bak-* 2>/dev/null | tail -n +11 | xargs -r rm -f
fi

echo "▶ building binary"
sudo -u "$APP_USER" "$GO_BIN" build -o "${APP_DIR}/amazingtrak" ./...
echo "✓ built ${APP_DIR}/amazingtrak"

echo "▶ restarting ${SERVICE}"
systemctl restart "$SERVICE"

# Health check.
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
