#!/usr/bin/env bash
#
# recompile-dev.sh — local/dev recompile + restart.
#
# Rebuilds the binary and (re)starts the app on your workstation. DB schema
# migrations (new columns / tables) apply automatically on startup via
# applySchema/runMigrations — there is no manual SQL step.
#
# Usage:
#   ./recompile-dev.sh            # build, then run in the background (nohup)
#   ./recompile-dev.sh --run      # build, then run in the foreground
#   ./recompile-dev.sh --build    # build only (compile check, no run)

set -euo pipefail
cd "$(dirname "$0")"

MODE="bg"
case "${1:-}" in
  --build) MODE="build" ;;
  --run)   MODE="fg" ;;
  "" )     MODE="bg" ;;
  *) echo "unknown option: $1 (use --build | --run | none)"; exit 2 ;;
esac

# Load .env so PORT/DB_PATH/ADMIN_SECRET etc. are set, same as start.sh.
if [ -f .env ]; then
  set -a; source .env; set +a
fi
PORT="${PORT:-8080}"

echo "▶ go build ./..."
go build ./...
echo "✓ build OK"

[ "$MODE" = "build" ] && { echo "build-only; not starting."; exit 0; }

# Stop any previous dev instance on this port so the rebuild takes effect.
if command -v fuser >/dev/null 2>&1; then
  fuser -k "${PORT}/tcp" 2>/dev/null || true
  sleep 1
fi

if [ "$MODE" = "fg" ]; then
  echo "▶ running in foreground on :${PORT} (Ctrl-C to stop)"
  exec go run .
fi

echo "▶ starting in background on :${PORT}"
nohup go run . > app.log 2>&1 &
echo "  pid $!  (logs: app.log)"

# Health check.
for i in $(seq 1 15); do
  if curl -sf "http://localhost:${PORT}/healthz" >/dev/null 2>&1; then
    echo "✓ healthy: http://localhost:${PORT}/healthz"
    exit 0
  fi
  sleep 1
done
echo "✗ did not become healthy in 15s — check app.log" >&2
exit 1
