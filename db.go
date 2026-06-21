package main

import (
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_journal=WAL&_timeout=5000&_fk=true&_synchronous=NORMAL")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := applySchema(db); err != nil {
		return nil, err
	}
	if err := runMigrations(db); err != nil {
		return nil, err
	}
	return db, nil
}

// migrationApplied reports whether a numbered migration has already run.
func migrationApplied(db *sql.DB, version int) bool {
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version=?`, version).Scan(&n)
	return n > 0
}

// markMigration records that a numbered migration has run.
func markMigration(db *sql.DB, version int) {
	db.Exec(`INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)`, version)
}

func runMigrations(db *sql.DB) error {
	// Historic idempotent ADD COLUMN migrations — safe to re-run every startup.
	// SQLite ignores "duplicate column name" errors, so these are a no-op on
	// databases that already have the column.
	db.Exec(`ALTER TABLE trains ADD COLUMN map_media_id INTEGER`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN notification_email TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE corridors ADD COLUMN schedule_url TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN rate_per_minute INTEGER NOT NULL DEFAULT 1`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN rate_per_hour INTEGER NOT NULL DEFAULT 5`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN rate_per_day INTEGER NOT NULL DEFAULT 20`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN site_name TEXT NOT NULL DEFAULT 'AmazingTrak'`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN favicon_path TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE admin_users ADD COLUMN notification_email TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE stops ADD COLUMN slug TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN register_rate_per_hour INTEGER NOT NULL DEFAULT 5`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN register_rate_per_day INTEGER NOT NULL DEFAULT 20`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN comment_rate_per_hour INTEGER NOT NULL DEFAULT 10`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN comment_rate_per_day INTEGER NOT NULL DEFAULT 50`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN admin_theme TEXT NOT NULL DEFAULT 'default'`)
	// "Spam" is tracked as a flag layered on top of the existing pending/rejected
	// statuses rather than a new CHECK-constrained enum value (SQLite can't ALTER
	// a CHECK constraint without a full table rebuild). This also means spam rows
	// are automatically included by any existing "status=pending"/"status=rejected"
	// filtering, satisfying "treat pending-spam the same as pending" for free.
	db.Exec(`ALTER TABLE suggestions ADD COLUMN is_spam INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE users ADD COLUMN is_spammer INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE suggestions ADD COLUMN caption TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE media ADD COLUMN tags TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE suggestions ADD COLUMN tags TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE suggestions ADD COLUMN auto_approved INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE media ADD COLUMN is_best INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE train_stops ADD COLUMN runs_weekday INTEGER NOT NULL DEFAULT 1`)
	db.Exec(`ALTER TABLE train_stops ADD COLUMN runs_weekend INTEGER NOT NULL DEFAULT 1`)
	// user accounts & permission levels
	db.Exec(`ALTER TABLE admin_users ADD COLUMN permission_level INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE suggestions ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE SET NULL`)
	db.Exec(`ALTER TABLE media ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE SET NULL`)
	// Corridor conductor: a registered user assigned to maintain the corridor's trains.
	db.Exec(`ALTER TABLE corridors ADD COLUMN conductor_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL`)
	// Email (Resend) settings — all optional; email is off unless enabled.
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN sender_email TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN email_enabled INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN verify_expiry_hours INTEGER NOT NULL DEFAULT 24`)
	// Trusted-tier (approved/auto_approved users) rate limits; anon limits stay in the existing columns.
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN trusted_rate_per_hour INTEGER NOT NULL DEFAULT 30`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN trusted_rate_per_day INTEGER NOT NULL DEFAULT 100`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN trusted_comment_rate_per_hour INTEGER NOT NULL DEFAULT 30`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN trusted_comment_rate_per_day INTEGER NOT NULL DEFAULT 100`)
	// Highest pending-items threshold (1/10/100) the admin has already been emailed about (hysteresis state).
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN pending_notify_level INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN admin_compact INTEGER NOT NULL DEFAULT 0`)
	// When the current email-verification token was last sent (for expiry).
	db.Exec(`ALTER TABLE users ADD COLUMN confirm_sent_at TEXT`)
	// User auto-approval policy toggles.
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN auto_approve_on_confirm INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN auto_approve_on_video INTEGER NOT NULL DEFAULT 1`)
	// Password reset: a single-use token + timestamp, kept separate from the
	// email-confirmation token so resetting a password never disturbs verification.
	db.Exec(`ALTER TABLE users ADD COLUMN reset_token TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE users ADD COLUMN reset_sent_at TEXT`)
	// Brute-force protection: a sequential failed-login counter (reset on any
	// successful login or password reset) and a hard-lock flag set when the
	// counter crosses the lockout threshold. Cleared by an admin unlock or a
	// password reset.
	db.Exec(`ALTER TABLE users ADD COLUMN failed_login_count INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`ALTER TABLE users ADD COLUMN login_locked INTEGER NOT NULL DEFAULT 0`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_reset_token ON users(reset_token)`)
	// Case-insensitive username uniqueness, preserving the user's chosen display
	// casing (the original case-sensitive UNIQUE on username stays). Best-effort:
	// if legacy case-duplicate usernames predate this index it simply won't be
	// created, and the handler-level case-insensitive check still blocks new ones.
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_ci ON users(lower(username))`)
	// Speed up the per-account / per-IP login throttle lookups.
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_login_attempts_user_time ON login_attempts(username, created_at)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_login_attempts_ip_time ON login_attempts(ip_hash, created_at)`)
	// One account per email address (case-insensitive), for non-empty addresses
	// only — blank emails stay allowed and unconstrained. Best-effort: if legacy
	// duplicate emails predate this index it simply won't be created, and the
	// registration handler's own check still blocks new duplicates.
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(lower(email)) WHERE email != ''`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_user ON suggestions(user_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_train_type ON media(train_id, media_type)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_corridor_type ON media(corridor_id, media_type)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_train_best ON media(train_id, is_best)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_train_status ON suggestions(train_id, status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_status ON suggestions(status)`)
	if err := migrateStopSlugs(db); err != nil {
		return err
	}

	// ── Versioned migrations (non-idempotent) ──────────────────────────────
	// Add new migrations below using the next integer. Each runs exactly once
	// and is recorded in schema_migrations.

	// Migration 1: generalize comments to reference a train OR a corridor.
	// SQLite can't drop the NOT NULL on train_id in place, so rebuild the table
	// (make train_id nullable, add corridor_id + CHECK) and copy existing rows.
	if !migrationApplied(db, 1) {
		if err := migrateCommentsTrainOrCorridor(db); err != nil {
			return err
		}
		markMigration(db, 1)
	}

	return nil
}

// migrateCommentsTrainOrCorridor rebuilds the comments table so a comment can
// belong to a train or a corridor (mirrors the media table's pattern).
func migrateCommentsTrainOrCorridor(db *sql.DB) error {
	// Skip if already generalized (e.g. fresh DB created by applySchema).
	var hasCorridor int
	db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('comments') WHERE name='corridor_id'`).Scan(&hasCorridor)
	if hasCorridor > 0 {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmts := []string{
		`CREATE TABLE comments_new (
			id INTEGER PRIMARY KEY,
			train_id INTEGER REFERENCES trains(id) ON DELETE CASCADE,
			corridor_id INTEGER REFERENCES corridors(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			body TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
			submitter_ip_hash TEXT,
			rejection_reason TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reviewed_at TEXT,
			reviewed_by INTEGER REFERENCES admin_users(id),
			CHECK (train_id IS NOT NULL OR corridor_id IS NOT NULL)
		)`,
		`INSERT INTO comments_new (id, train_id, corridor_id, user_id, body, status, submitter_ip_hash, rejection_reason, created_at, reviewed_at, reviewed_by)
			SELECT id, train_id, NULL, user_id, body, status, submitter_ip_hash, rejection_reason, created_at, reviewed_at, reviewed_by FROM comments`,
		`DROP TABLE comments`,
		`ALTER TABLE comments_new RENAME TO comments`,
		`CREATE INDEX IF NOT EXISTS idx_comments_train_status ON comments(train_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_corridor_status ON comments(corridor_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id)`,
	}
	for _, s := range stmts {
		if _, err := tx.Exec(s); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func migrateStopSlugs(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, name, COALESCE(station_code,'') FROM stops WHERE slug='' OR slug IS NULL`)
	if err != nil {
		return err
	}
	type stopRow struct {
		id   int64
		name string
		code string
	}
	var stops []stopRow
	for rows.Next() {
		var s stopRow
		if err := rows.Scan(&s.id, &s.name, &s.code); err != nil {
			rows.Close()
			return err
		}
		stops = append(stops, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, s := range stops {
		var slug string
		if s.code != "" {
			slug = strings.ToLower(s.code)
		} else {
			slug = slugify(s.name)
		}
		db.Exec(`UPDATE stops SET slug=? WHERE id=?`, slug, s.id)
	}
	return nil
}

func applySchema(db *sql.DB) error {
	_, err := db.Exec(`
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

CREATE TABLE IF NOT EXISTS schema_migrations (
	version INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS corridors (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	slug TEXT NOT NULL UNIQUE,
	region TEXT,
	description TEXT,
	on_time_percent REAL,
	service_quality_summary TEXT,
	hero_train_id INTEGER,
	hero_media_id INTEGER,
	thumbnail_media_id INTEGER,
	is_active INTEGER NOT NULL DEFAULT 1,
	sort_order INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS trains (
	id INTEGER PRIMARY KEY,
	corridor_id INTEGER NOT NULL REFERENCES corridors(id) ON DELETE CASCADE,
	train_number TEXT NOT NULL,
	display_name TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	direction TEXT,
	notes TEXT,
	hero_media_id INTEGER,
	thumbnail_media_id INTEGER,
	map_media_id INTEGER,
	is_active INTEGER NOT NULL DEFAULT 1,
	sort_order INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(corridor_id, train_number)
);

CREATE TABLE IF NOT EXISTS stops (
	id INTEGER PRIMARY KEY,
	corridor_id INTEGER NOT NULL REFERENCES corridors(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	station_code TEXT,
	latitude REAL,
	longitude REAL,
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS train_stops (
	id INTEGER PRIMARY KEY,
	train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
	stop_id INTEGER NOT NULL REFERENCES stops(id) ON DELETE CASCADE,
	sort_order INTEGER NOT NULL DEFAULT 0,
	scheduled_arrival TEXT,
	scheduled_departure TEXT,
	UNIQUE(train_id, stop_id)
);

CREATE TABLE IF NOT EXISTS media (
	id INTEGER PRIMARY KEY,
	train_id INTEGER REFERENCES trains(id) ON DELETE CASCADE,
	corridor_id INTEGER REFERENCES corridors(id) ON DELETE CASCADE,
	media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video', 'website')),
	source_type TEXT NOT NULL CHECK (source_type IN ('url', 'upload', 'paste', 'seed')),
	url TEXT,
	local_path TEXT,
	original_filename TEXT,
	stored_filename TEXT,
	title TEXT,
	caption TEXT,
	source_domain TEXT,
	latitude REAL,
	longitude REAL,
	location_name TEXT,
	location_source TEXT DEFAULT 'unknown',
	is_published INTEGER NOT NULL DEFAULT 1,
	added_by TEXT NOT NULL DEFAULT 'admin',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CHECK (train_id IS NOT NULL OR corridor_id IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS suggestions (
	id INTEGER PRIMARY KEY,
	train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
	url TEXT NOT NULL,
	title TEXT,
	media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video')),
	source_domain TEXT,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	submitter_ip_hash TEXT,
	submitter_user_agent TEXT,
	rejection_reason TEXT,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed_at TEXT,
	reviewed_by INTEGER REFERENCES admin_users(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	train_id INTEGER REFERENCES trains(id) ON DELETE CASCADE,
	corridor_id INTEGER REFERENCES corridors(id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	body TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	submitter_ip_hash TEXT,
	rejection_reason TEXT,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed_at TEXT,
	reviewed_by INTEGER REFERENCES admin_users(id),
	CHECK (train_id IS NOT NULL OR corridor_id IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS email_errors (
	id INTEGER PRIMARY KEY,
	to_addr TEXT NOT NULL,
	subject TEXT NOT NULL,
	error TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS email_verifications (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS conductor_requests (
	id INTEGER PRIMARY KEY,
	corridor_id INTEGER NOT NULL REFERENCES corridors(id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	message TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed_at TEXT,
	reviewed_by INTEGER REFERENCES admin_users(id)
);

CREATE TABLE IF NOT EXISTS admin_users (
	id INTEGER PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	must_change_password INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_login_at TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	admin_user_id INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
	csrf_token TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TEXT NOT NULL,
	ip_hash TEXT,
	user_agent TEXT
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL DEFAULT '',
	password_hash TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending'
		CHECK (status IN ('pending', 'confirmed', 'approved', 'auto_approved')),
	email_confirmed INTEGER NOT NULL DEFAULT 0,
	confirm_token TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_login_at TEXT
);

CREATE TABLE IF NOT EXISTS user_sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	csrf_token TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TEXT NOT NULL,
	ip_hash TEXT,
	user_agent TEXT
);

CREATE TABLE IF NOT EXISTS audit_log (
	id INTEGER PRIMARY KEY,
	admin_user_id INTEGER REFERENCES admin_users(id),
	action TEXT NOT NULL,
	entity_type TEXT NOT NULL,
	entity_id INTEGER NOT NULL,
	detail TEXT,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rate_limit_log (
	id INTEGER PRIMARY KEY,
	ip_hash TEXT NOT NULL,
	action TEXT NOT NULL DEFAULT 'suggest',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS login_attempts (
	id INTEGER PRIMARY KEY,
	ip_hash TEXT NOT NULL,
	username TEXT,
	succeeded INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS site_preferences (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	default_theme TEXT NOT NULL DEFAULT 'auto' CHECK (default_theme IN ('light', 'dark', 'auto')),
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trains_corridor ON trains(corridor_id);
CREATE INDEX IF NOT EXISTS idx_media_train ON media(train_id);
CREATE INDEX IF NOT EXISTS idx_media_corridor ON media(corridor_id);
CREATE INDEX IF NOT EXISTS idx_suggestions_train_status ON suggestions(train_id, status);
CREATE INDEX IF NOT EXISTS idx_rate_limit_ip_time ON rate_limit_log(ip_hash, created_at);
CREATE INDEX IF NOT EXISTS idx_login_attempts_ip_time ON login_attempts(ip_hash, created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_comments_train_status ON comments(train_id, status);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_conductor_req_status ON conductor_requests(status);
CREATE INDEX IF NOT EXISTS idx_email_verifications_user_time ON email_verifications(user_id, created_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_conductor_req_pending ON conductor_requests(corridor_id, user_id) WHERE status='pending';

CREATE TRIGGER IF NOT EXISTS corridors_updated_at AFTER UPDATE ON corridors BEGIN
	UPDATE corridors SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
CREATE TRIGGER IF NOT EXISTS trains_updated_at AFTER UPDATE ON trains BEGIN
	UPDATE trains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
CREATE TRIGGER IF NOT EXISTS media_updated_at AFTER UPDATE ON media BEGIN
	UPDATE media SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
`)
	return err
}
