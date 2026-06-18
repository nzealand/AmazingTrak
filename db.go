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
	// migrations: add columns if the DB predates them
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
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_user ON suggestions(user_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_train_type ON media(train_id, media_type)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_corridor_type ON media(corridor_id, media_type)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_media_train_best ON media(train_id, is_best)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_train_status ON suggestions(train_id, status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_suggestions_status ON suggestions(status)`)
	if err := migrateStopSlugs(db); err != nil {
		return nil, err
	}
	return db, nil
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
	train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	body TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	submitter_ip_hash TEXT,
	rejection_reason TEXT,
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
