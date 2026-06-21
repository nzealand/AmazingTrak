package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const sessionCtxKey contextKey = "session"
const adminUserCtxKey contextKey = "adminUser"

// clientIP resolves the real client IP. The Go app listens only on loopback
// behind nginx, which is configured to OVERWRITE (not append) the forwarded
// headers with the true peer address. We therefore trust X-Real-IP /
// X-Forwarded-For only when the immediate peer is itself loopback (i.e. the
// local reverse proxy). A directly-connecting client cannot forge these.
func clientIP(r *http.Request) string {
	host := r.RemoteAddr
	if h, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		host = h
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
			return xrip
		}
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			// nginx sets a single trusted value; take the first token defensively.
			if comma := strings.IndexByte(fwd, ','); comma >= 0 {
				return strings.TrimSpace(fwd[:comma])
			}
			return strings.TrimSpace(fwd)
		}
	}
	return host
}

func hashIP(r *http.Request) string {
	h := sha256.Sum256([]byte(clientIP(r)))
	return fmt.Sprintf("%x", h[:8])
}

func (app *App) createSession(db *sql.DB, userID int64, r *http.Request) (Session, error) {
	id := newToken()
	csrf := newToken()
	expires := time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02T15:04:05")
	ipHash := hashIP(r)
	ua := r.UserAgent()
	if len(ua) > 200 {
		ua = ua[:200]
	}
	_, err := db.Exec(
		`INSERT INTO sessions (id, admin_user_id, csrf_token, expires_at, ip_hash, user_agent) VALUES (?, ?, ?, ?, ?, ?)`,
		id, userID, csrf, expires, ipHash, ua,
	)
	if err != nil {
		return Session{}, err
	}
	return Session{ID: id, AdminUserID: userID, CSRFToken: csrf, ExpiresAt: expires}, nil
}

func (app *App) getSession(r *http.Request) (*Session, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, nil
	}
	var s Session
	err = app.db.QueryRow(
		`SELECT id, admin_user_id, csrf_token, expires_at FROM sessions WHERE id=? AND expires_at > datetime('now')`,
		c.Value,
	).Scan(&s.ID, &s.AdminUserID, &s.CSRFToken, &s.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (app *App) deleteSession(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE id=?`, id)
	return err
}

func (app *App) setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     app.adminCookiePath,
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   app.secureCookies,
		SameSite: http.SameSiteStrictMode,
	})
}

func (app *App) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     app.adminCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   app.secureCookies,
		SameSite: http.SameSiteStrictMode,
	})
}

func setFlash(w http.ResponseWriter, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    msg,
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getFlash(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("flash")
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return c.Value
}

func (app *App) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := app.getSession(r)
		if err != nil || s == nil {
			http.Redirect(w, r, app.adminPrefix+"/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), sessionCtxKey, s)
		// Load the admin account so handlers/templates can gate on permission level.
		if au, err := getAdminUser(app.db, s.AdminUserID); err == nil {
			ctx = context.WithValue(ctx, adminUserCtxKey, &au)
		}
		next(w, r.WithContext(ctx))
	}
}

// requirePermission wraps requireAdmin and enforces a minimum cumulative
// permission level (1=submissions, 2=trains, 3=corridors, 4=settings).
// The central admin (level 0) always passes.
func (app *App) requirePermission(level int, next http.HandlerFunc) http.HandlerFunc {
	return app.requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		au := adminFromCtx(r)
		if au == nil || (au.PermissionLevel != 0 && au.PermissionLevel < level) {
			http.Error(w, "Forbidden — insufficient permissions", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

// requireCentralAdmin restricts a route to the central admin (permission_level 0).
func (app *App) requireCentralAdmin(next http.HandlerFunc) http.HandlerFunc {
	return app.requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		au := adminFromCtx(r)
		if au == nil || !au.IsCentral() {
			http.Error(w, "Forbidden — central admin only", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

func sessionFromCtx(r *http.Request) *Session {
	s, _ := r.Context().Value(sessionCtxKey).(*Session)
	return s
}

func adminFromCtx(r *http.Request) *AdminUser {
	au, _ := r.Context().Value(adminUserCtxKey).(*AdminUser)
	return au
}

func (app *App) checkCSRF(r *http.Request) bool {
	s := sessionFromCtx(r)
	if s == nil {
		return false
	}
	return r.FormValue("csrf_token") == s.CSRFToken
}

func (app *App) checkLoginThrottle(ipHash string) bool {
	var count int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM login_attempts WHERE ip_hash=? AND succeeded=0 AND created_at > datetime('now', '-15 minutes')`,
		ipHash,
	).Scan(&count)
	return count >= 5
}

func (app *App) recordLoginAttempt(ipHash, username string, succeeded bool) {
	s := 0
	if succeeded {
		s = 1
	}
	app.db.Exec(`INSERT INTO login_attempts (ip_hash, username, succeeded) VALUES (?, ?, ?)`, ipHash, username, s)
}

// loginHardLockThreshold is the number of sequential failed logins (since the
// account's last success or password reset) that permanently locks the account
// until an admin unlocks it or the user resets their password.
const loginHardLockThreshold = 20

// countFailedLoginsByUser counts failed login attempts for a username (matched
// case-insensitively) within the given SQLite time modifier, e.g. "-5 minutes".
func (app *App) countFailedLoginsByUser(username, since string) int {
	var n int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM login_attempts WHERE lower(username)=lower(?) AND succeeded=0 AND created_at > datetime('now', ?)`,
		username, since,
	).Scan(&n)
	return n
}

// countFailedLoginsByIP counts failed login attempts from an IP within the
// given SQLite time modifier.
func (app *App) countFailedLoginsByIP(ipHash, since string) int {
	var n int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM login_attempts WHERE ip_hash=? AND succeeded=0 AND created_at > datetime('now', ?)`,
		ipHash, since,
	).Scan(&n)
	return n
}

// userLoginBlocked enforces user-account login throttling and returns a
// human-readable reason when a login attempt should be refused before the
// password is even checked. Rules (most restrictive wins):
//
//	per IP:      10 fails / 60 min, 50 fails / 24 h
//	per account: 3 fails / 5 min, 5 fails / 24 h
//	hard lock:   20 sequential fails since last success/reset (needs admin/reset)
//
// Blocked attempts are NOT recorded as new failures, so the time-windowed
// limits drain on their own once attempts stop.
func (app *App) userLoginBlocked(username, ipHash string) (bool, string) {
	// IP-based limits apply regardless of whether the username exists.
	if app.countFailedLoginsByIP(ipHash, "-60 minutes") >= 10 {
		return true, "Too many failed login attempts from your network. Please wait up to an hour and try again."
	}
	if app.countFailedLoginsByIP(ipHash, "-1 day") >= 50 {
		return true, "Too many failed login attempts from your network. Please try again in 24 hours."
	}
	// Account-based limits only apply to real accounts.
	u, err := userByUsername(app.db, username)
	if err != nil {
		return false, ""
	}
	if u.LoginLocked {
		return true, "This account is locked after too many failed login attempts. Reset your password, or contact an admin to unlock it."
	}
	if app.countFailedLoginsByUser(username, "-5 minutes") >= 3 {
		return true, "Too many failed attempts for this account. Please wait a few minutes and try again."
	}
	if app.countFailedLoginsByUser(username, "-1 day") >= 5 {
		return true, "Too many failed attempts for this account. Please wait 24 hours, or reset your password to regain access sooner."
	}
	return false, ""
}

// registerFailedLogin bumps the sequential failed-login counter for an existing
// account (matched case-insensitively) and hard-locks it once the counter
// reaches the threshold. No-op for unknown usernames.
func (app *App) registerFailedLogin(username string) {
	app.db.Exec(
		`UPDATE users SET failed_login_count = failed_login_count + 1,
		   login_locked = CASE WHEN failed_login_count + 1 >= ? THEN 1 ELSE login_locked END
		 WHERE lower(username)=lower(?)`,
		loginHardLockThreshold, username,
	)
}

// clearLoginFailures resets the sequential counter and lock for an account
// (after a successful login, password reset, or admin unlock).
func (app *App) clearLoginFailures(userID int64) {
	app.db.Exec(`UPDATE users SET failed_login_count=0, login_locked=0 WHERE id=?`, userID)
}

func (app *App) authenticateAdmin(username, password string) (*AdminUser, error) {
	var u AdminUser
	err := app.db.QueryRow(
		`SELECT id, username, password_hash, must_change_password FROM admin_users WHERE username=?`, username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.MustChangePassword)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, nil
	}
	return &u, nil
}

// ----- Registered user (public account) sessions -----

func (app *App) createUserSession(userID int64, r *http.Request) (string, error) {
	id := newToken()
	csrf := newToken()
	expires := time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02T15:04:05")
	ua := r.UserAgent()
	if len(ua) > 200 {
		ua = ua[:200]
	}
	_, err := app.db.Exec(
		`INSERT INTO user_sessions (id, user_id, csrf_token, expires_at, ip_hash, user_agent) VALUES (?, ?, ?, ?, ?, ?)`,
		id, userID, csrf, expires, hashIP(r), ua,
	)
	return id, err
}

// getUserSession returns the logged-in user for the request, or (nil, "") when
// there is no valid user session. The second return is the session's CSRF token.
func (app *App) getUserSession(r *http.Request) (*User, string) {
	c, err := r.Cookie("usersession")
	if err != nil {
		return nil, ""
	}
	var userID int64
	var csrf string
	err = app.db.QueryRow(
		`SELECT user_id, csrf_token FROM user_sessions WHERE id=? AND expires_at > datetime('now')`,
		c.Value,
	).Scan(&userID, &csrf)
	if err != nil {
		return nil, ""
	}
	u, err := userByID(app.db, userID)
	if err != nil {
		return nil, ""
	}
	return &u, csrf
}

// requireUser gates a public route behind having a logged-in user session,
// redirecting anonymous visitors to /login instead of 403ing.
func (app *App) requireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := app.getUserSession(r)
		if u == nil {
			setFlash(w, "Please log in to view that page.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (app *App) deleteUserSession(id string) {
	app.db.Exec(`DELETE FROM user_sessions WHERE id=?`, id)
}

// conductorGuard returns the logged-in user, their session CSRF token, and
// whether that user is the conductor of corridorID. Anonymous visitors and
// non-conductors get (user-or-nil, csrf, false). Conductor management routes use
// this to gate access without ever exposing the admin area.
func (app *App) conductorGuard(r *http.Request, corridorID int64) (*User, string, bool) {
	u, csrf := app.getUserSession(r)
	if u == nil || u.IsSpammer {
		return u, csrf, false
	}
	ok, err := isConductorOf(app.db, u.ID, corridorID)
	if err != nil {
		return u, csrf, false
	}
	return u, csrf, ok
}

func (app *App) setUserSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "usersession",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   app.secureCookies,
		SameSite: http.SameSiteStrictMode,
	})
}

func (app *App) clearUserSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "usersession",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   app.secureCookies,
		SameSite: http.SameSiteStrictMode,
	})
}

func (app *App) authenticateUser(username, password string) (*User, error) {
	u, err := userByUsername(app.db, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, nil
	}
	return &u, nil
}

func (app *App) purgeExpiredSessions() {
	app.db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
	app.db.Exec(`DELETE FROM user_sessions WHERE expires_at <= datetime('now')`)
}

func (app *App) purgeOldRateLimitLogs() {
	app.db.Exec(`DELETE FROM rate_limit_log WHERE created_at < datetime('now', '-25 hours')`)
	app.db.Exec(`DELETE FROM login_attempts WHERE created_at < datetime('now', '-30 days')`)
	app.db.Exec(`DELETE FROM audit_log WHERE created_at < datetime('now', '-30 days')`)
}

func (app *App) logAudit(adminUserID int64, action, entityType string, entityID int64, detail string) {
	app.db.Exec(
		`INSERT INTO audit_log (admin_user_id, action, entity_type, entity_id, detail) VALUES (?, ?, ?, ?, ?)`,
		adminUserID, action, entityType, entityID, detail,
	)
}
