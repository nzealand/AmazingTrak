<a id="readme-top"></a>

[![Go][Go-shield]][Go-url]
[![SQLite][SQLite-shield]][SQLite-url]

<br />
<div align="center">
  <h3 align="center">AmazingTrak</h3>

  <p align="center">
    A self-hosted Amtrak train tracker — browse corridors, trains, photos, and route schedules.
    <br />
    <a href="#getting-started">Getting Started</a>
    &middot;
    <a href="#usage">Usage</a>
    &middot;
    <a href="#roadmap">Roadmap</a>
  </p>
</div>

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li>
      <a href="#vps-deployment">VPS Deployment</a>
      <ul>
        <li><a href="#automated-setup">Automated Setup</a></li>
        <li><a href="#what-the-script-installs">What the Script Installs</a></li>
        <li><a href="#deploying-updates">Deploying Updates</a></li>
        <li><a href="#directory-layout">Directory Layout</a></li>
        <li><a href="#managing-the-service">Managing the Service</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

## About The Project

AmazingTrak is a single-binary Go web app for tracking and showcasing Amtrak train routes. It lets you organize corridors, trains, and associated media (photos, videos, links), accept public photo suggestions, and display route schedules with stop data sourced from FRA performance reports.

Features:
* Browse corridors and trains with hero images and route maps
* Admin panel to manage trains, corridors, media, and suggestions
* Paste or upload train photos directly — EXIF GPS auto-extracted
* Public photo suggestion form with multi-layer spam protection (honeypot fields, JS-verified token, timing check, configurable rate limits)
* Amtrak schedule links per corridor; live status link to TransitDocs per train
* Interactive route map with corridor page links in popups
* In-process page cache for the homepage — zero DB work on repeated loads
* Compiled template cache — templates parsed once at startup, not per request
* Long-lived browser cache for static assets (CSS/JS) via content-hash versioning
* Dark/light/auto theme toggle

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Go][Go-shield]][Go-url] — single-binary server, standard library HTTP
* [![SQLite][SQLite-shield]][SQLite-url] — embedded database, WAL mode
* [Leaflet](https://leafletjs.com/) — interactive map for route display and photo locations
* [nginx](https://nginx.org/) + [certbot](https://certbot.eff.org/) — reverse proxy and TLS termination in production

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Getting Started

### Prerequisites

* Go 1.22+
* GCC (for SQLite CGo driver)
* Optional: `caddy` for production TLS

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/your_username/amazingtrak.git
   cd amazingtrak
   ```

2. Build the binary
   ```sh
   go build -o amazingtrak .
   ```

3. Copy and edit the environment file
   ```sh
   cp .env.example .env
   # Set ADMIN_SECRET, ADMIN_USERNAME, ADMIN_PASSWORD, PORT, DB_PATH
   ```

4. Run
   ```sh
   ./start.sh
   ```

5. Verify
   ```sh
   curl http://localhost:8000/healthz
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

| URL | Description |
|-----|-------------|
| `/` | Public home — list of corridors |
| `/corridors/{slug}` | Corridor page with trains |
| `/trains/{slug}` | Train page with photos, schedule, links |
| `/trains/{slug}/suggest` | Public photo/video suggestion form |
| `/{ADMIN_SECRET}` | Admin dashboard |
| `/{ADMIN_SECRET}/trains/{id}/media` | Upload or paste train photos |
| `/healthz` | Health check |

The admin URL is `/<ADMIN_SECRET>` — keep the secret out of browser history by setting it to something non-guessable (e.g. `mytoken/dashboard`).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Configuration

Environment variables (all optional — defaults shown):

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `DB_PATH` | `amazingtrak.db` | SQLite file path |
| `UPLOADS_DIR` | `./uploads` | Directory for image uploads |
| `ADMIN_SECRET` | `admin/secret` | Secret URL path segment for admin |
| `ADMIN_USERNAME` | `admin` | Seeded admin username |
| `ADMIN_PASSWORD` | `secret` | Seeded admin password |
| `BASE_URL` | `http://localhost:PORT` | Public base URL (enables `Secure` cookies when `https://`) |
| `SMTP_HOST` | — | SMTP server for suggestion notifications |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | — | SMTP username / from address |
| `SMTP_PASS` | — | SMTP password |

Place these in `.env` — `start.sh` sources it automatically.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## VPS Deployment

### Automated Setup

`setup-vps.sh` bootstraps a fresh Ubuntu 22.04/24.04 VPS end-to-end. Run it as root:

```sh
sudo bash setup-vps.sh
```

The script is interactive — it prompts for your domain, admin credentials, optional SMTP settings, and whether to enable SSL. No flags required.

> **Prerequisites:** a DNS A record for your domain pointing to the VPS IP, and outbound internet access from the VPS.

### What the Script Installs

| Component | Detail |
|---|---|
| **Go 1.22.5** | Downloaded from go.dev; skipped if ≥ 1.22 already present |
| **nginx** | Reverse proxy on port 80/443; serves `/uploads/` directly from disk |
| **certbot** | Let's Encrypt SSL certificate with auto HTTP→HTTPS redirect |
| **UFW** | Firewall locked to SSH + Nginx only |
| **logrotate** | Daily log rotation, 14-day retention |
| **systemd service** | Hardened unit (`NoNewPrivileges`, `PrivateTmp`, `ProtectSystem`) |

The Go binary is built on the VPS from source. The SQLite database and uploads are stored in `/var/lib/amazingtrak/` — separate from the source tree, making backups straightforward.

### Deploying Updates

```sh
cd /opt/amazingtrak
git pull
go build -o amazingtrak ./...
systemctl restart amazingtrak
```

DB schema migrations run automatically on startup — no manual SQL required. Verify the deploy:

```sh
curl -s http://localhost:8000/healthz
journalctl -u amazingtrak -n 30
```

If the nginx config changed (e.g. after re-running `setup-vps.sh`):

```sh
nginx -t && systemctl reload nginx
```

### Directory Layout

```
/opt/amazingtrak/          # source + binary
  amazingtrak              # compiled binary
  .env                     # secrets (chmod 600)

/var/lib/amazingtrak/      # persistent data (back this up)
  amazingtrak.db           # SQLite database
  uploads/images/          # uploaded photos

/var/log/amazingtrak/      # application logs
  app.log
```

### Managing the Service

```sh
systemctl status amazingtrak      # check status
systemctl restart amazingtrak     # restart after update
journalctl -u amazingtrak -f      # stream logs
tail -f /var/log/amazingtrak/app.log
```

Nginx config lives at `/etc/nginx/sites-available/amazingtrak`. After editing, run `nginx -t && systemctl reload nginx`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Roadmap

- [x] Corridor and train management
- [x] Image upload with EXIF GPS extraction
- [x] Paste-to-upload in admin
- [x] Public suggestion form with rate limiting
- [x] Route stop schedules
- [x] Dark/light theme toggle
- [ ] Train-to-stop schedule editor in admin
- [ ] On-time performance stats display
- [ ] Public search

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the MIT License.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS -->
[Go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[SQLite-shield]: https://img.shields.io/badge/SQLite-07405E?style=for-the-badge&logo=sqlite&logoColor=white
[SQLite-url]: https://www.sqlite.org/
