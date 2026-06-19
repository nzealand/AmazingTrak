# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Run (dev)
go run .

# Run with env overrides
PORT=8080 DB_PATH=amazingtrak.db ADMIN_SECRET=admin/secret go run .

# Run via start.sh (loads .env, starts with nohup)
./start.sh

# Health check
curl http://localhost:8080/healthz
```

No test files exist yet. There is no linter configured.

## Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `DB_PATH` | `amazingtrak.db` | SQLite file path |
| `UPLOADS_DIR` | `./uploads` | Directory for image uploads |
| `ADMIN_SECRET` | `admin/secret` | Secret URL path segment for admin |
| `ADMIN_USERNAME` | `admin` | Seeded admin username |
| `ADMIN_PASSWORD` | `secret` | Seeded admin password |

`.env` is loaded by `start.sh`. The admin URL is `/<ADMIN_SECRET>` — the first path segment sets the cookie path.

## Architecture

Single-binary Go app: `go run .` compiles `main.go`, `auth.go`, `db.go`, `handlers_admin.go`, `handlers_public.go`, `media.go`, `models.go`, `validate.go` into one server. Templates and static files are embedded via `//go:embed` so the binary is self-contained.

**Data layer** (`db.go`, `models.go`): SQLite (WAL mode, single connection, FK enforcement). Schema applied idempotently on startup via `applySchema`. Seed data written once on first run. Query functions live in `models.go` and use raw `database/sql` with `?` placeholders — never `fmt.Sprintf` in SQL. Key tables: `corridors → trains → media` (images/videos/websites), `suggestions` (public-submitted pending media), `sessions`, `audit_log`, `rate_limit_log`, `login_attempts`.

**Auth** (`auth.go`): bcrypt passwords, 32-byte hex session tokens stored in DB, CSRF tokens per session (synchronizer pattern). Sessions expire after 24h; a background goroutine purges them hourly. Login is throttled (5 failures in 15 min blocks the IP). `requireAdmin` middleware injects the session into context; `checkCSRF` validates form tokens.

**Handlers** (`handlers_public.go`, `handlers_admin.go`): Method+path routing via Go 1.22's `http.NewServeMux` with typed patterns (`GET /foo`, `POST /foo`). Public pages render `templates/base.html` + a page template; admin pages render `templates/admin_base.html` + an admin page template. Flash messages use a short-lived cookie.

**Media** (`media.go`): Image uploads are re-encoded as JPEG (strips EXIF metadata), GPS coordinates extracted before stripping. Uploaded files saved to `uploads/images/<slug>-<N>.jpg`. URL-based media is validated and stored as-is. Each train can have one hero image (`hero_media_id` on the train row) and one thumbnail.

**Validation** (`validate.go`): Public URL submissions are restricted to an allowlist (`youtube.com`, `youtu.be`, `vimeo.com`, `flickr.com`, `imgur.com`, `railpictures.net`, `rrpicturearchives.net`, `instagram.com`). Rate limited at 3/hour and 10/day per IP hash, server-side timing check (≥4s), honeypot field, and duplicate detection across both `suggestions` and `media` tables. Admin URLs accept any `http`/`https` URL.

**Templates**: Public templates use `templates/base.html`; admin templates use `templates/admin_base.html`. Both are parsed fresh per request from the embedded FS. Template functions are defined in `main.go` (`funcMap`) and cover YouTube/Vimeo ID extraction, null SQL types, geo formatting, and safe URL casting.

## Key Invariants

- Every admin form POST must include `csrf_token` matching the session's token.
- `media.train_id` and `media.corridor_id` may each be NULL but not both (CHECK constraint).
- Public suggestions only ever enter `suggestions` table with `status='pending'`; approval copies them into `media` with `added_by='public_approved'`.
- Hero image: only one per train (`hero_media_id` on the `trains` row); setting a new hero does NOT use a DB unique index — it's enforced in handler logic by clearing the old value first.
- Uploaded images are always re-encoded to JPEG regardless of original format.
- Ask before any destructive action, including deleting files, overwriting migrations, force-pushing, or changing production config.
- The site is live on Digital Ocean. All changes must be upgrade-safe: new DB columns must have defaults or be nullable, no column renames or drops without a migration plan, no breaking changes to existing data formats, and schema changes must apply cleanly via `applySchema` on a running instance without data loss.

After completing a significant change always update the admin version to include the latest timestamp e.g. 3.6.6 11:46AM and if the user requests to increment the version add that version to the CHANGELOG.md e.g. ===3.6.7=====
