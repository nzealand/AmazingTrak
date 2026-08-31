package main

import (
	"database/sql"
	"strconv"
	"strings"

	_ "embed"

	_ "modernc.org/sqlite"
)

// station_coords.csv is a one-time reference snapshot (code,lat,lon) pulled
// from the Amtraker v3 API's /stations endpoint — the same free feed
// livetrains.go already polls for positions, keyed by the same official
// Amtrak station codes our own stops.station_code values use. It exists to
// backfill latitude/longitude on stops seeded without coordinates (see
// migrateStopCoordinates); it is not fetched at runtime.
//
//go:embed station_coords.csv
var stationCoordsCSV string

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
	// Live train positions on /map (Amtraker feed). Off by default: enabling it
	// starts polling a third-party API, so that must be an explicit admin choice.
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN live_trains_enabled INTEGER NOT NULL DEFAULT 0`)
	// How often (seconds) to poll the upstream feed; admin-adjustable between
	// 90s (matches Amtraker's own refresh cadence — no point going faster)
	// and 600s (10 min, matching other public trackers like transitdocs.com).
	// Defaults to 2 min rather than the 90s floor to go a bit easier on the
	// free upstream API out of the box.
	db.Exec(`ALTER TABLE site_preferences ADD COLUMN live_trains_poll_seconds INTEGER NOT NULL DEFAULT 120`)
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
	// Live train data sources (Amtraker for Amtrak, GTFS-RT feeds for everyone
	// else). One row per agency, keyed by a short slug, so adding a new source
	// is a data row rather than a schema change — see livetrains.go.
	db.Exec(`CREATE TABLE IF NOT EXISTS live_sources (
		source_key TEXT PRIMARY KEY,
		display_name TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 0,
		poll_seconds INTEGER NOT NULL DEFAULT 120,
		api_key TEXT NOT NULL DEFAULT '',
		last_error TEXT NOT NULL DEFAULT '',
		last_polled_at TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	) STRICT, WITHOUT ROWID`)
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
	// Re-run every startup rather than as a one-shot versioned migration: it's
	// naturally idempotent (only touches NULL coordinates). Handles the
	// upgrade case (existing stops rows already seeded long ago); main.go
	// calls this again after seedDB to also handle a brand new install,
	// since seeding runs after this point.
	if err := migrateStopCoordinates(db); err != nil {
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

	// Migration 2: Amtrak 816/817/880/806/807/887/888 were seeded into Acela
	// but are actually Northeast Regional numbers, while 2107/2117 were seeded
	// into Northeast Regional but are actually Acela numbers (Acela trains are
	// all in the 2100-2295 range). Move each by train_number, not id, since
	// ids can differ between the dev DB and production.
	if !migrationApplied(db, 2) {
		if err := fixAcelaNortheastRegionalMixup(db); err != nil {
			return err
		}
		markMigration(db, 2)
	}

	// Migration 3: reword the Acela corridor description to drop promotional
	// language ("premium", "fastest in the Americas") in favor of neutral,
	// factual wording, per trademark/branding-sensitivity review.
	if !migrationApplied(db, 3) {
		if _, err := db.Exec(`UPDATE corridors SET description=? WHERE slug='amtrak-acela' AND description=?`,
			"High-speed passenger service operated by Amtrak between Boston, New York, Philadelphia, Baltimore, and Washington, D.C.",
			"Amtrak's premium high-speed service connecting Boston, New York, Philadelphia, and Washington D.C. The fastest train in the Americas, reaching speeds up to 150 mph along the Northeast Corridor.",
		); err != nil {
			return err
		}
		markMigration(db, 3)
	}

	// Migration 4: introduce the multi-source live-tracking config table
	// (live_sources), seeded from the pre-existing Amtrak-only
	// site_preferences columns so an admin's existing enabled/interval choice
	// carries forward unchanged, plus a disabled-by-default Caltrain row
	// (needs a free 511.org API key before it can be turned on — see
	// caltrain.go). site_preferences.live_trains_* are left in place, just no
	// longer read after this point.
	if !migrationApplied(db, 4) {
		// On a brand new install this migration runs before seedDB has
		// inserted the site_preferences row (openDB → runMigrations happens
		// before seedDB — see main.go), so this query legitimately finds no
		// row yet; fall back to the same defaults getSitePrefs itself uses
		// rather than leaving the Scan destinations at their zero values.
		amtrakEnabled, amtrakPoll := 0, 120
		db.QueryRow(`SELECT COALESCE(live_trains_enabled,0), COALESCE(live_trains_poll_seconds,120) FROM site_preferences WHERE id=1`).
			Scan(&amtrakEnabled, &amtrakPoll)
		if _, err := db.Exec(`INSERT OR IGNORE INTO live_sources (source_key, display_name, enabled, poll_seconds) VALUES (?,?,?,?)`,
			"amtrak", "Amtrak (Amtraker)", amtrakEnabled, amtrakPoll); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT OR IGNORE INTO live_sources (source_key, display_name, enabled, poll_seconds) VALUES (?,?,?,?)`,
			"caltrain", "Caltrain (511.org)", 0, 90); err != nil {
			return err
		}
		markMigration(db, 4)
	}

	// Migration 5: add the Caltrain corridor to already-seeded (upgraded)
	// databases. seedDB only ever runs against an empty corridors table, so
	// appending to corridorSeeds in seed.go only reaches brand-new installs —
	// an existing production database needs this explicit insert instead. No
	// trains are added here; matching happens by whatever train_number rows
	// an admin adds under this corridor (see caltrain.go).
	//
	// Guarded on corridors already being non-empty: this migration runs
	// during openDB, before seedDB (see main.go), so on a genuinely fresh
	// install corridors is still empty here — inserting into it at this
	// point would make seedDB's own "already seeded?" check (COUNT(*) > 0)
	// see a false positive and skip ALL seeding, including the admin user.
	// A fresh install instead gets Caltrain from corridorSeeds in seed.go.
	if !migrationApplied(db, 5) {
		var corridorCount int
		db.QueryRow(`SELECT COUNT(*) FROM corridors`).Scan(&corridorCount)
		if corridorCount > 0 {
			if _, err := db.Exec(`
				INSERT INTO corridors (name, slug, region, description, sort_order)
				SELECT 'Caltrain', 'caltrain', 'California',
					'Commuter rail service along the San Francisco Peninsula and South Bay, connecting San Francisco to San Jose and Gilroy. Operated by the Peninsula Corridor Joint Powers Board, not Amtrak.',
					COALESCE((SELECT MAX(sort_order) FROM corridors), 0) + 1
				WHERE NOT EXISTS (SELECT 1 FROM corridors WHERE slug='caltrain')`); err != nil {
				return err
			}
		}
		markMigration(db, 5)
	}

	// Migration 6: seed a starter set of Caltrain trains on already-seeded
	// (upgraded) databases, mirroring the same numbers added to trainSeeds
	// in seed.go for fresh installs — empirically observed from a live
	// 511.org GTFS-RT feed sample (2026-08-31, late morning weekday), not a
	// full published timetable. A no-op on a fresh install (no caltrain
	// corridor exists yet at this point — see migration 5's comment); that
	// case is instead covered by trainSeeds.
	if !migrationApplied(db, 6) {
		for i, num := range []string{"118", "119", "120", "121", "122", "123", "124", "125", "126", "127", "128", "129"} {
			if _, err := db.Exec(`
				INSERT OR IGNORE INTO trains (corridor_id, train_number, display_name, slug, sort_order)
				SELECT id, ?, ?, ?, ? FROM corridors WHERE slug='caltrain'`,
				num, "Caltrain "+num, "caltrain-"+num, i+1); err != nil {
				return err
			}
		}
		markMigration(db, 6)
	}

	// Migration 7: register ACE, MBTA Commuter Rail, and LIRR as live
	// sources — disabled by default. ACE shares 511.org's account/key model
	// with Caltrain (needs a key); MBTA and LIRR need no key (see mbta.go,
	// mtalirr.go for why).
	if !migrationApplied(db, 7) {
		rows := []struct {
			key, name string
			poll      int
		}{
			{"ace", "ACE (511.org)", 90},
			{"mbta", "MBTA Commuter Rail", 90},
			{"lirr", "LIRR", 90},
		}
		for _, r := range rows {
			if _, err := db.Exec(`INSERT OR IGNORE INTO live_sources (source_key, display_name, enabled, poll_seconds) VALUES (?,?,0,?)`,
				r.key, r.name, r.poll); err != nil {
				return err
			}
		}
		markMigration(db, 7)
	}

	// Migration 8: add the ACE, MBTA Commuter Rail, and LIRR corridors to
	// already-seeded (upgraded) databases — same fresh-install hazard and
	// guard as migration 5 (see its comment); a fresh install instead gets
	// these from corridorSeeds in seed.go.
	if !migrationApplied(db, 8) {
		var corridorCount int
		db.QueryRow(`SELECT COUNT(*) FROM corridors`).Scan(&corridorCount)
		if corridorCount > 0 {
			corridors := []struct{ name, slug, region, desc string }{
				{"ACE", "ace", "California",
					"Altamont Corridor Express commuter rail connecting Stockton and San Jose through the Altamont Pass. Operated by the San Joaquin Regional Rail Commission, not Amtrak."},
				{"MBTA Commuter Rail", "mbta-commuter-rail", "New England",
					"Boston-area commuter rail network operated by the MBTA across its Providence, Franklin, Needham, Fairmount, Worcester, Fitchburg, Lowell, Haverhill, Newburyport, Kingston, Greenbush, and New Bedford lines. Not Amtrak."},
				{"LIRR", "lirr", "Northeast",
					"Long Island Rail Road commuter service connecting Manhattan (Penn Station/Grand Central) to Long Island. Operated by the MTA, not Amtrak."},
			}
			for _, c := range corridors {
				if _, err := db.Exec(`
					INSERT INTO corridors (name, slug, region, description, sort_order)
					SELECT ?, ?, ?, ?, COALESCE((SELECT MAX(sort_order) FROM corridors), 0) + 1
					WHERE NOT EXISTS (SELECT 1 FROM corridors WHERE slug=?)`,
					c.name, c.slug, c.region, c.desc, c.slug); err != nil {
					return err
				}
			}
		}
		markMigration(db, 8)
	}

	// Migration 9: seed a starter set of real, currently-observed train
	// numbers for ACE, MBTA Commuter Rail, and LIRR on already-seeded
	// (upgraded) databases — mirrors trainSeeds in seed.go for fresh
	// installs. Empirically observed from each source's live feed on
	// 2026-08-31; not a full published timetable for any of the three. A
	// no-op on a fresh install (no corridors exist yet — see migration 8's
	// comment).
	if !migrationApplied(db, 9) {
		trainSets := []struct {
			slug, prefix, label string
			nums                []string
		}{
			{"ace", "ace", "ACE", []string{"1", "2", "3", "4", "5", "6", "7", "8"}},
			{"mbta-commuter-rail", "mbta", "MBTA", []string{
				"45", "46", "148", "149", "247", "250", "347", "348", "421", "424",
				"545", "550", "645", "646", "743", "843", "847", "848", "852", "946",
				"1028", "1069", "1147", "1246", "1419", "1653", "1656", "1748",
				"1919", "2021", "2028", "2030",
			}},
			{"lirr", "lirr", "LIRR", []string{
				"8", "13", "38", "43", "69", "93", "150", "152", "156", "157", "159",
				"161", "163", "258", "352", "353", "354", "355", "452", "455", "552",
				"557", "652", "655", "752", "753", "755", "809", "853", "854", "855",
				"951", "1297", "1298", "1553", "1554", "1555", "1556", "1557", "1558",
				"1653", "1752", "1753", "1951", "1952", "1953", "1954", "1955", "1956",
				"1957", "1958", "1959", "2752", "2753", "2754", "2755", "2898", "2900",
				"2909", "2911", "2913", "2919",
			}},
		}
		for _, ts := range trainSets {
			for i, num := range ts.nums {
				if _, err := db.Exec(`
					INSERT OR IGNORE INTO trains (corridor_id, train_number, display_name, slug, sort_order)
					SELECT id, ?, ?, ?, ? FROM corridors WHERE slug=?`,
					num, ts.label+" "+num, ts.prefix+"-"+num, i+1, ts.slug); err != nil {
					return err
				}
			}
		}
		markMigration(db, 9)
	}

	// Migration 10: register SEPTA and Brightline as live sources — disabled
	// by default. Neither needs an API key (see septa.go, brightline.go).
	if !migrationApplied(db, 10) {
		rows := []struct {
			key, name string
			poll      int
		}{
			{"septa", "SEPTA Regional Rail", 90},
			{"brightline", "Brightline", 90},
		}
		for _, r := range rows {
			if _, err := db.Exec(`INSERT OR IGNORE INTO live_sources (source_key, display_name, enabled, poll_seconds) VALUES (?,?,0,?)`,
				r.key, r.name, r.poll); err != nil {
				return err
			}
		}
		markMigration(db, 10)
	}

	// Migration 11: add the SEPTA Regional Rail and Brightline corridors to
	// already-seeded (upgraded) databases — same fresh-install hazard and
	// guard as migration 5/8 (see migration 5's comment); a fresh install
	// instead gets these from corridorSeeds in seed.go.
	if !migrationApplied(db, 11) {
		var corridorCount int
		db.QueryRow(`SELECT COUNT(*) FROM corridors`).Scan(&corridorCount)
		if corridorCount > 0 {
			corridors := []struct{ name, slug, region, desc string }{
				{"SEPTA Regional Rail", "septa", "Mid-Atlantic",
					"Philadelphia-area commuter rail network operated by SEPTA across its Airport, Chestnut Hill East, Chestnut Hill West, Cynwyd, Fox Chase, Lansdale/Doylestown, Manayunk/Norristown, Media/Wawa, Paoli/Thorndale, Trenton, Warminster, West Trenton, and Wilmington/Newark lines. Not Amtrak."},
				{"Brightline", "brightline", "Florida",
					"Higher-speed intercity rail connecting Miami, Fort Lauderdale, Boca Raton, West Palm Beach, and Orlando. Privately operated by Brightline Trains Florida, not Amtrak."},
			}
			for _, c := range corridors {
				if _, err := db.Exec(`
					INSERT INTO corridors (name, slug, region, description, sort_order)
					SELECT ?, ?, ?, ?, COALESCE((SELECT MAX(sort_order) FROM corridors), 0) + 1
					WHERE NOT EXISTS (SELECT 1 FROM corridors WHERE slug=?)`,
					c.name, c.slug, c.region, c.desc, c.slug); err != nil {
					return err
				}
			}
		}
		markMigration(db, 11)
	}

	// Migration 12: seed a starter set of real, currently-observed train
	// numbers for SEPTA Regional Rail and Brightline on already-seeded
	// (upgraded) databases — mirrors trainSeeds in seed.go for fresh
	// installs. Empirically observed from each source's live feed on
	// 2026-08-31; not a full published timetable for either. A no-op on a
	// fresh install (no corridors exist yet — see migration 11's comment).
	if !migrationApplied(db, 12) {
		trainSets := []struct {
			slug, prefix, label string
			nums                []string
		}{
			{"septa", "septa", "SEPTA", []string{
				"452", "457", "846", "849", "850", "1085", "2388", "2579", "2591",
				"3537", "3541", "3548", "3552", "4587", "4750", "4754", "5344", "5347",
				"5349", "5351", "5355", "5496", "6228", "6229", "6256", "6313", "6342",
				"6345", "6550", "6807", "7455", "9224", "9229", "9231", "9546", "9589",
				"9593", "9748", "9750", "9757", "9759",
			}},
			{"brightline", "brightline", "Brightline", []string{
				"5150", "5151", "5152", "5340", "5347", "5348", "5355", "5356", "5363",
			}},
		}
		for _, ts := range trainSets {
			for i, num := range ts.nums {
				if _, err := db.Exec(`
					INSERT OR IGNORE INTO trains (corridor_id, train_number, display_name, slug, sort_order)
					SELECT id, ?, ?, ?, ? FROM corridors WHERE slug=?`,
					num, ts.label+" "+num, ts.prefix+"-"+num, i+1, ts.slug); err != nil {
					return err
				}
			}
		}
		markMigration(db, 12)
	}

	return nil
}

func fixAcelaNortheastRegionalMixup(db *sql.DB) error {
	var acelaID, nerID int64
	if err := db.QueryRow(`SELECT id FROM corridors WHERE slug='amtrak-acela'`).Scan(&acelaID); err != nil {
		return nil // corridor doesn't exist on this instance; nothing to fix
	}
	if err := db.QueryRow(`SELECT id FROM corridors WHERE slug='amtrak-northeast-regional'`).Scan(&nerID); err != nil {
		return nil
	}
	toNER := []string{"816", "817", "880", "806", "807", "887", "888"}
	for _, n := range toNER {
		if _, err := db.Exec(`UPDATE trains SET corridor_id=? WHERE corridor_id=? AND train_number=?`, nerID, acelaID, n); err != nil {
			return err
		}
	}
	toAcela := []string{"2107", "2117"}
	for _, n := range toAcela {
		if _, err := db.Exec(`UPDATE trains SET corridor_id=? WHERE corridor_id=? AND train_number=?`, acelaID, nerID, n); err != nil {
			return err
		}
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

// migrateStopCoordinates fills in latitude/longitude for any stop whose
// station_code matches a row in the embedded Amtrak station snapshot and
// doesn't already have coordinates set — never overwrites a value an admin
// (or a future editor) has already entered.
func migrateStopCoordinates(db *sql.DB) error {
	stmt, err := db.Prepare(`UPDATE stops SET latitude=?, longitude=? WHERE station_code=? AND latitude IS NULL`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, line := range strings.Split(stationCoordsCSV, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		lat, err1 := strconv.ParseFloat(parts[1], 64)
		lon, err2 := strconv.ParseFloat(parts[2], 64)
		if err1 != nil || err2 != nil {
			continue
		}
		if _, err := stmt.Exec(lat, lon, parts[0]); err != nil {
			return err
		}
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

CREATE TABLE IF NOT EXISTS vantage_spots (
	id TEXT PRIMARY KEY,
	latitude REAL NOT NULL,
	longitude REAL NOT NULL,
	title TEXT NOT NULL DEFAULT '',
	caption TEXT NOT NULL DEFAULT '',
	media_type TEXT NOT NULL DEFAULT 'image' CHECK (media_type IN ('image', 'video')),
	url TEXT NOT NULL DEFAULT '',
	local_path TEXT NOT NULL DEFAULT '',
	original_filename TEXT NOT NULL DEFAULT '',
	stored_filename TEXT NOT NULL DEFAULT '',
	source_domain TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	submitter_ip_hash TEXT NOT NULL DEFAULT '',
	submitter_user_agent TEXT NOT NULL DEFAULT '',
	rejection_reason TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed_at TEXT NOT NULL DEFAULT '',
	reviewed_by INTEGER REFERENCES admin_users(id)
) STRICT, WITHOUT ROWID;

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
CREATE INDEX IF NOT EXISTS idx_vantage_spots_status ON vantage_spots(status);
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
