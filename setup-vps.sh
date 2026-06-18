#!/usr/bin/env bash
# setup-vps.sh — AmazingTrak VPS bootstrap
# Tested on Ubuntu 22.04 / 24.04 LTS
# Run as root: bash setup-vps.sh
set -euo pipefail

# ── colours ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GRN='\033[0;32m'; YLW='\033[1;33m'; BLD='\033[1m'; RST='\033[0m'
info()  { echo -e "${GRN}[+]${RST} $*"; }
warn()  { echo -e "${YLW}[!]${RST} $*"; }
fatal() { echo -e "${RED}[✗]${RST} $*" >&2; exit 1; }
hr()    { echo -e "${BLD}────────────────────────────────────────────────${RST}"; }

# ── preflight ────────────────────────────────────────────────────────────────
[[ $EUID -eq 0 ]] || fatal "Run as root (sudo bash $0)"
. /etc/os-release 2>/dev/null || true
[[ "${ID:-}" == "ubuntu" || "${ID_LIKE:-}" == *"debian"* ]] \
  || warn "Untested distro '$ID'. Continuing anyway — expect apt."

GO_VERSION="1.22.5"
APP_USER="amazingtrak"
APP_DIR="/opt/amazingtrak"
DATA_DIR="/var/lib/amazingtrak"
LOG_DIR="/var/log/amazingtrak"
SERVICE="amazingtrak"

hr
echo -e "${BLD}AmazingTrak VPS Setup${RST}"
hr

# ── collect config ────────────────────────────────────────────────────────────
ask() {
  local prompt="$1" default="${2:-}" var
  if [[ -n "$default" ]]; then
    read -rp "  $prompt [$default]: " var
    echo "${var:-$default}"
  else
    read -rp "  $prompt: " var
    echo "$var"
  fi
}
ask_secret() {
  local prompt="$1" var
  read -rsp "  $prompt: " var; echo
  echo "$var"
}

echo ""
echo -e "${BLD}Configuration${RST} (press Enter to accept defaults)"
echo ""

DOMAIN=$(ask        "Domain name (e.g. amazingtrak.example.com)")
[[ -n "$DOMAIN" ]] || fatal "Domain name is required."

REPO_URL=$(ask      "Git repository URL (leave blank to copy from current dir)")
APP_PORT=$(ask      "App listen port" "8000")
BASE_URL=$(ask      "Public base URL" "https://${DOMAIN}")

echo ""
echo -e "${BLD}Admin credentials${RST}"
ADMIN_SECRET=$(ask_secret  "Admin URL secret (e.g. myrandompath/admin)")
[[ -n "$ADMIN_SECRET" ]] || fatal "Admin secret is required."
ADMIN_USER=$(ask           "Admin username" "admin")
ADMIN_PASS=$(ask_secret    "Admin password")
[[ -n "$ADMIN_PASS" ]] || fatal "Admin password is required."

echo ""
echo -e "${BLD}Email notifications (optional — press Enter to skip)${RST}"
SMTP_HOST=$(ask      "SMTP host")
SMTP_PORT=$(ask      "SMTP port" "587")
SMTP_USER=$(ask      "SMTP username")
SMTP_PASS=$(ask_secret "SMTP password")
NOTIFY_EMAIL=$(ask   "Notification email for suggestions")

echo ""
SETUP_SSL=$(ask "Set up SSL with Let's Encrypt? (y/n)" "y")
echo ""
hr

# ── system packages ───────────────────────────────────────────────────────────
info "Updating package lists…"
apt-get update -q

info "Installing system dependencies…"
apt-get install -y -q \
  curl git nginx certbot python3-certbot-nginx \
  ufw logrotate ca-certificates

# ── Go ────────────────────────────────────────────────────────────────────────
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"

if command -v go &>/dev/null; then
  INSTALLED=$(go version | awk '{print $3}' | sed 's/go//')
  info "Go already installed: $INSTALLED"
  # Upgrade if older than required
  if [[ "$(printf '%s\n' "1.22" "$INSTALLED" | sort -V | head -1)" != "1.22" ]]; then
    warn "Go $INSTALLED is older than 1.22 — upgrading to $GO_VERSION"
    rm -rf /usr/local/go
  else
    GO_SKIP=1
  fi
fi

if [[ -z "${GO_SKIP:-}" ]]; then
  info "Downloading Go ${GO_VERSION}…"
  curl -fsSL "$GO_URL" -o "/tmp/${GO_TAR}"
  rm -rf /usr/local/go
  tar -C /usr/local -xzf "/tmp/${GO_TAR}"
  rm "/tmp/${GO_TAR}"
  ln -sf /usr/local/go/bin/go   /usr/local/bin/go
  ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
  info "Go $(go version) installed."
fi

# ── app user ─────────────────────────────────────────────────────────────────
if ! id "$APP_USER" &>/dev/null; then
  info "Creating system user '$APP_USER'…"
  useradd --system --shell /usr/sbin/nologin \
          --home-dir "$APP_DIR" --create-home "$APP_USER"
fi

# ── directories ───────────────────────────────────────────────────────────────
info "Creating directories…"
mkdir -p "${APP_DIR}" "${DATA_DIR}/uploads/images" "${LOG_DIR}"
chown -R "${APP_USER}:${APP_USER}" "${DATA_DIR}" "${LOG_DIR}"

# ── source code ───────────────────────────────────────────────────────────────
if [[ -n "$REPO_URL" ]]; then
  info "Cloning repository…"
  if [[ -d "${APP_DIR}/.git" ]]; then
    sudo -u "$APP_USER" git -C "$APP_DIR" pull
  else
    rm -rf "${APP_DIR:?}"/*
    sudo -u "$APP_USER" git clone "$REPO_URL" "$APP_DIR"
  fi
else
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  info "Copying source from ${SCRIPT_DIR}…"
  rsync -a --exclude='.env' --exclude='amazingtrak.db*' \
           --exclude='uploads/' --exclude='nohup.out' \
    "${SCRIPT_DIR}/" "${APP_DIR}/"
fi

chown -R "${APP_USER}:${APP_USER}" "${APP_DIR}"

# Ensure data.sql is present (required for go:embed at build time)
if [[ -f "${APP_DIR}/old/data.sql" && ! -f "${APP_DIR}/data.sql" ]]; then
  cp "${APP_DIR}/old/data.sql" "${APP_DIR}/data.sql"
fi
[[ -f "${APP_DIR}/data.sql" ]] || fatal "data.sql missing — place it in ${APP_DIR}/ before building."

# ── .env ─────────────────────────────────────────────────────────────────────
info "Writing .env…"
cat > "${APP_DIR}/.env" <<ENV
PORT=${APP_PORT}
DB_PATH=${DATA_DIR}/amazingtrak.db
UPLOADS_DIR=${DATA_DIR}/uploads
BASE_URL=${BASE_URL}
ADMIN_SECRET=${ADMIN_SECRET}
ADMIN_USERNAME=${ADMIN_USER}
ADMIN_PASSWORD=${ADMIN_PASS}
SMTP_HOST=${SMTP_HOST}
SMTP_PORT=${SMTP_PORT}
SMTP_USER=${SMTP_USER}
SMTP_PASS=${SMTP_PASS}
ENV
chmod 600 "${APP_DIR}/.env"
chown "${APP_USER}:${APP_USER}" "${APP_DIR}/.env"

# Write notification email to DB prefs after first run — stored in settings table,
# not .env, but we note it for the operator here.
[[ -n "$NOTIFY_EMAIL" ]] && \
  warn "Set notification email in Admin → Settings after first login (not stored in .env)."

# ── build ─────────────────────────────────────────────────────────────────────
info "Building AmazingTrak binary…"
cd "$APP_DIR"
sudo -u "$APP_USER" /usr/local/go/bin/go build -o "${APP_DIR}/amazingtrak" ./...
info "Binary built at ${APP_DIR}/amazingtrak"

# ── systemd service ───────────────────────────────────────────────────────────
info "Installing systemd service…"
cat > "/etc/systemd/system/${SERVICE}.service" <<UNIT
[Unit]
Description=AmazingTrak Amtrak tracker
After=network.target
Wants=network.target

[Service]
Type=simple
User=${APP_USER}
WorkingDirectory=${APP_DIR}
EnvironmentFile=${APP_DIR}/.env
ExecStart=${APP_DIR}/amazingtrak
Restart=on-failure
RestartSec=5s
StandardOutput=append:${LOG_DIR}/app.log
StandardError=append:${LOG_DIR}/app.log

# Hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=${DATA_DIR} ${LOG_DIR}
CapabilityBoundingSet=

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable "${SERVICE}"
systemctl restart "${SERVICE}"
info "Service started."

# ── logrotate ────────────────────────────────────────────────────────────────
cat > "/etc/logrotate.d/${SERVICE}" <<LOGROTATE
${LOG_DIR}/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    postrotate
        systemctl kill -s USR1 ${SERVICE} 2>/dev/null || true
    endscript
}
LOGROTATE

# ── nginx ─────────────────────────────────────────────────────────────────────
info "Configuring nginx…"
cat > "/etc/nginx/sites-available/${SERVICE}" <<NGINX
server {
    listen 80;
    server_name ${DOMAIN};

    # Max upload size (matches app's 32 MB body limit)
    client_max_body_size 32M;

    # Static uploads — serve directly from disk without hitting Go
    location /uploads/ {
        alias ${DATA_DIR}/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    location / {
        proxy_pass         http://127.0.0.1:${APP_PORT};
        proxy_http_version 1.1;
        proxy_set_header   Host              \$host;
        proxy_set_header   X-Real-IP         \$remote_addr;
        proxy_set_header   X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 120s;
    }
}
NGINX

ln -sf "/etc/nginx/sites-available/${SERVICE}" \
       "/etc/nginx/sites-enabled/${SERVICE}"

# Remove default site if still present
rm -f /etc/nginx/sites-enabled/default

nginx -t
systemctl enable nginx
systemctl reload nginx

# ── SSL ───────────────────────────────────────────────────────────────────────
if [[ "${SETUP_SSL,,}" == "y" ]]; then
  info "Obtaining SSL certificate via Let's Encrypt…"
  certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos \
    --email "${SMTP_USER:-admin@${DOMAIN}}" --redirect \
    || warn "certbot failed — run 'certbot --nginx -d ${DOMAIN}' manually after DNS propagates."
fi

# ── firewall ─────────────────────────────────────────────────────────────────
info "Configuring UFW firewall…"
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow "Nginx Full"
ufw --force enable

# ── health check ─────────────────────────────────────────────────────────────
sleep 3
if curl -fsS "http://127.0.0.1:${APP_PORT}/healthz" &>/dev/null; then
  info "Health check passed — app is running."
else
  warn "Health check failed. Check logs: journalctl -u ${SERVICE} -n 50"
fi

# ── summary ───────────────────────────────────────────────────────────────────
hr
echo -e "${GRN}${BLD}Setup complete!${RST}"
hr
echo ""
echo -e "  Site:        ${BLD}${BASE_URL}${RST}"
echo -e "  Admin:       ${BLD}${BASE_URL}/${ADMIN_SECRET}${RST}"
echo -e "  App logs:    journalctl -u ${SERVICE} -f"
echo -e "  Tail log:    tail -f ${LOG_DIR}/app.log"
echo -e "  DB:          ${DATA_DIR}/amazingtrak.db"
echo -e "  Uploads:     ${DATA_DIR}/uploads/"
echo ""
echo -e "  Service:     systemctl {status|restart|stop} ${SERVICE}"
echo ""
warn "Store your Admin URL secret safely — it is not recoverable from this script."
hr
