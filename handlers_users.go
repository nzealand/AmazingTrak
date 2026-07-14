package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"sort"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// maxPasswordBytes is bcrypt's hard input limit; bytes past this are silently
// ignored by bcrypt, so we reject longer passwords rather than truncate them.
const maxPasswordBytes = 72

// validUsername allows 3–30 chars of letters, digits, underscore, hyphen, dot.
func validUsername(s string) bool {
	if len(s) < 3 || len(s) > 30 {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' || r == '.') {
			return false
		}
	}
	return true
}

func validEmail(s string) bool {
	if s == "" {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil
}

func (app *App) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	if u, _ := app.getUserSession(r); u != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	setTimingCookie(w)
	app.renderPublic(w, r, "register.html", publicPage{
		Title: "Create an account — " + getSiteName(),
		Flash: getFlash(w, r),
		Data:  map[string]any{"EmailTaken": r.URL.Query().Get("emailtaken") == "1", "EmailEnabled": app.emailEnabled()},
	})
}

func (app *App) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	fail := func(msg string) {
		setFlash(w, msg)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}

	// Bot defenses: honeypot fields and minimum fill time fail silently with a
	// fake success message so bots can't tell their submission was dropped.
	if r.FormValue("website") != "" || r.FormValue("a") != "" || r.FormValue("b") != "ok" || !checkTiming(r) {
		setFlash(w, "Welcome! Your account is ready.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ipHash := hashIP(r)
	prefs, _ := getSitePrefs(app.db)
	perHour := prefs.RegisterRatePerHour
	if perHour <= 0 {
		perHour = defaultRatePerHour
	}
	perDay := prefs.RegisterRatePerDay
	if perDay <= 0 {
		perDay = defaultRatePerDay
	}
	// Registration limit is site-wide (counts all new accounts in the window,
	// not per-IP), with a wait-time message.
	if blocked, reason := app.checkRegisterRateLimit(perHour, perDay); blocked {
		fail(reason)
		return
	}

	if !validUsername(username) {
		fail("Username must be 3–30 characters (letters, digits, . _ -).")
		return
	}
	if email != "" && !validEmail(email) {
		fail("Please enter a valid email address, or leave it blank.")
		return
	}
	if len(password) < 8 {
		fail("Password must be at least 8 characters.")
		return
	}
	if len(password) > maxPasswordBytes {
		fail(fmt.Sprintf("Password must be no more than %d characters.", maxPasswordBytes))
		return
	}
	if password != confirm {
		fail("Passwords do not match.")
		return
	}
	if _, err := userByUsername(app.db, username); err == nil {
		fail("That username is already taken.")
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Database error", 500)
		return
	}
	// One account per email address. Only enforced for supplied (non-blank)
	// emails — registering without an email stays allowed.
	if email != "" {
		if used, err := emailInUse(app.db, email, 0); err != nil {
			http.Error(w, "Database error", 500)
			return
		} else if used {
			http.Redirect(w, r, "/register?emailtaken=1", http.StatusSeeOther)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	// Email confirmation is optional: only issue a token when email is enabled
	// and the user supplied an address. Its absence never blocks registration.
	confirmToken := ""
	emailConfigured := app.emailEnabled() && email != ""
	if emailConfigured {
		confirmToken = newToken()
	}

	res, err := app.db.Exec(
		`INSERT INTO users (username, email, password_hash, status, email_confirmed, confirm_token, confirm_sent_at) VALUES (?, ?, ?, 'pending', 0, ?, CASE WHEN ?='' THEN NULL ELSE CURRENT_TIMESTAMP END)`,
		username, email, string(hash), confirmToken, confirmToken,
	)
	if err != nil {
		// Backstop for the unique-email index in case two registrations race
		// past the check above between the SELECT and this INSERT.
		if strings.Contains(err.Error(), "idx_users_email_unique") {
			http.Redirect(w, r, "/register?emailtaken=1", http.StatusSeeOther)
			return
		}
		// Username uniqueness (case-sensitive UNIQUE or the case-insensitive
		// index) racing past the pre-check lands here too.
		if strings.Contains(err.Error(), "idx_users_username_ci") || strings.Contains(err.Error(), "users.username") {
			fail("That username is already taken.")
			return
		}
		// Never surface raw DB/SQL errors to the user; log for diagnosis instead.
		log.Printf("register: insert user failed: %v", err)
		fail("Could not create your account. Please try again.")
		return
	}
	userID, _ := res.LastInsertId()
	app.recordActionRateLimit("register", ipHash)

	if emailConfigured {
		go app.sendVerifyEmail(email, confirmToken, app.verifyExpiryHours())
	}
	app.maybeNotifyPending()

	// Auto-login on successful registration.
	if sid, err := app.createUserSession(userID, r); err == nil {
		app.setUserSessionCookie(w, sid)
	}

	if emailConfigured {
		setFlash(w, "Welcome! Check your email to confirm your address — you can keep using the site in the meantime.")
	} else {
		setFlash(w, "Welcome! Your account is ready. An admin will review it shortly.")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) handleUserLoginForm(w http.ResponseWriter, r *http.Request) {
	if u, _ := app.getUserSession(r); u != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.renderPublic(w, r, "login.html", publicPage{
		Title: "Log in — " + getSiteName(),
		Flash: getFlash(w, r),
		Data:  map[string]any{"EmailEnabled": app.emailEnabled()},
	})
}

func (app *App) handleUserLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	ipHash := hashIP(r)
	// Accept either username or email address in the identifier field.
	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := r.FormValue("password")

	// Resolve the canonical username for throttle/audit so email addresses are
	// never written to login_attempts or used as account keys in rate-limit checks.
	throttleUsername := identifier
	if strings.Contains(identifier, "@") {
		if resolved, err := userByEmail(app.db, identifier); err == nil {
			throttleUsername = resolved.Username
		}
	}

	// Throttle before the password is checked; blocked attempts are not recorded
	// as new failures, so the time-windowed limits drain on their own.
	if blocked, reason := app.userLoginBlocked(throttleUsername, ipHash); blocked {
		setFlash(w, reason)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	u, err := app.authenticateUser(identifier, password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if u == nil {
		app.recordLoginAttempt(ipHash, throttleUsername, false)
		// Bump the per-account sequential counter (no-op for unknown usernames);
		// hard-locks the account once it reaches the threshold.
		app.registerFailedLogin(throttleUsername)
		setFlash(w, "Invalid username or password.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.recordLoginAttempt(ipHash, u.Username, true)
	// A successful login clears the sequential failure counter.
	app.clearLoginFailures(u.ID)
	app.db.Exec(`UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id=?`, u.ID)
	sid, err := app.createUserSession(u.ID, r)
	if err != nil {
		http.Error(w, "Session error", 500)
		return
	}
	app.setUserSessionCookie(w, sid)
	setFlash(w, "Logged in as "+u.Username+".")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleUserChangePassword lets a logged-in user change their own password by
// supplying their current password plus a new one.
func (app *App) handleUserChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	user, csrf := app.getUserSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.FormValue("csrf_token") != csrf {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	back := "/users/" + user.Username

	current := r.FormValue("current_password")
	newPw := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(current)); err != nil {
		setFlash(w, "Current password is incorrect.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if len(newPw) < 8 {
		setFlash(w, "New password must be at least 8 characters.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if len(newPw) > maxPasswordBytes {
		setFlash(w, fmt.Sprintf("New password must be no more than %d characters.", maxPasswordBytes))
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if newPw != confirm {
		setFlash(w, "New passwords do not match.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if _, err := app.db.Exec(`UPDATE users SET password_hash=?, reset_token='', reset_sent_at=NULL, failed_login_count=0, login_locked=0 WHERE id=?`, string(hash), user.ID); err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	// Invalidate all *other* sessions, then re-establish this one so the user
	// stays logged in but anyone using a stale session is logged out.
	app.db.Exec(`DELETE FROM user_sessions WHERE user_id=?`, user.ID)
	if sid, err := app.createUserSession(user.ID, r); err == nil {
		app.setUserSessionCookie(w, sid)
	}
	setFlash(w, "Password updated successfully.")
	http.Redirect(w, r, back, http.StatusSeeOther)
}

func (app *App) handleUserLogout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	u, csrf := app.getUserSession(r)
	if u == nil || r.FormValue("csrf_token") != csrf {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if c, err := r.Cookie("usersession"); err == nil {
		app.deleteUserSession(c.Value)
	}
	app.clearUserSessionCookie(w)
	setFlash(w, "You have been logged out.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// badge describes an earned milestone shown on a user's public profile.
type badge struct {
	Emoji string
	Label string
}

// badgesForSubmissions computes earned milestone badges from a user's
// submissions: first video, first rarity, then 10/100 of each.
func badgesForSubmissions(subs []Suggestion) []badge {
	videoCount := 0
	rarityCount := 0
	for _, s := range subs {
		if s.MediaType == "video" {
			videoCount++
		}
		rarityCount += s.RarityCount()
	}
	var out []badge
	if videoCount >= 1 {
		out = append(out, badge{"🎬", "First video"})
	}
	if rarityCount >= 1 {
		out = append(out, badge{"✨", "First rarity"})
	}
	if videoCount >= 10 {
		out = append(out, badge{"📹", "10 videos"})
	}
	if rarityCount >= 10 {
		out = append(out, badge{"💎", "10 rarities"})
	}
	if videoCount >= 100 {
		out = append(out, badge{"🏆", "100 videos"})
	}
	if rarityCount >= 100 {
		out = append(out, badge{"👑", "100 rarities"})
	}
	return out
}

// handleUserProfile shows a registered user's public profile: join date,
// earned badges, and their submissions sorted by rarity count (most first).
func (app *App) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	profileUser, err := userByUsername(app.db, username)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	subs, err := submissionsByUserID(app.db, profileUser.ID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	sort.SliceStable(subs, func(i, j int) bool {
		return subs[i].RarityCount() > subs[j].RarityCount()
	})

	comments, err := commentsByUserID(app.db, profileUser.ID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	type profileData struct {
		ProfileUser  *User
		Submissions  []Suggestion
		Badges       []badge
		Comments     []Comment
		EmailEnabled bool
	}
	app.renderPublic(w, r, "user_profile.html", publicPage{
		Title: profileUser.Username + " — " + getSiteName(),
		Flash: getFlash(w, r),
		Data: profileData{
			ProfileUser:  &profileUser,
			Submissions:  subs,
			Badges:       badgesForSubmissions(subs),
			Comments:     comments,
			EmailEnabled: app.emailEnabled(),
		},
	})
}

func (app *App) handleConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		setFlash(w, "Invalid confirmation link.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := userByConfirmToken(app.db, token)
	if err != nil {
		setFlash(w, "That confirmation link is invalid or has already been used.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Reject expired links. Tokens older than the configured window can't be used;
	// the user can request a fresh one.
	var expired int
	app.db.QueryRow(
		`SELECT CASE WHEN confirm_sent_at IS NOT NULL AND confirm_sent_at < datetime('now', ?) THEN 1 ELSE 0 END FROM users WHERE id=?`,
		fmt.Sprintf("-%d hours", app.verifyExpiryHours()), u.ID,
	).Scan(&expired)
	if expired == 1 {
		setFlash(w, "That confirmation link has expired. Log in and request a new one.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Confirm the email and clear the token. Bump pending → confirmed, but never
	// downgrade an already approved/auto-approved account.
	newStatus := u.Status
	if u.Status == "pending" {
		newStatus = "confirmed"
	}
	// If the admin has enabled auto-approve on confirmation, skip the manual
	// review step and promote straight to auto_approved.
	if newStatus == "confirmed" {
		if prefs, err := getSitePrefs(app.db); err == nil && prefs.AutoApproveOnConfirm {
			newStatus = "auto_approved"
		}
	}
	app.db.Exec(`UPDATE users SET email_confirmed=1, confirm_token='', status=? WHERE id=?`, newStatus, u.ID)
	setFlash(w, "Thanks — your email address is confirmed.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleResendVerification re-sends the email-verification link to the logged-in
// user. Rate limited to 3/day per user and 50/day site-wide. No-op (with a clear
// message) when email is disabled or the account has no/confirmed address.
func (app *App) handleResendVerification(w http.ResponseWriter, r *http.Request) {
	user, csrf := app.getUserSession(r)
	if user == nil {
		setFlash(w, "Please log in first.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.FormValue("csrf_token") != csrf {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	back := "/users/" + user.Username

	if !app.emailEnabled() {
		setFlash(w, "Email verification isn't available right now.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if user.Email == "" {
		setFlash(w, "Your account has no email address on file.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if user.EmailConfirmed {
		setFlash(w, "Your email is already confirmed.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}

	// Rate limits: 3/day per user, 50/day across the whole site.
	var userDay, siteDay int
	app.db.QueryRow(`SELECT COUNT(*) FROM email_verifications WHERE user_id=? AND created_at > datetime('now','-1 day')`, user.ID).Scan(&userDay)
	if userDay >= 3 {
		setFlash(w, "You've requested too many verification emails today. Please try again tomorrow.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	app.db.QueryRow(`SELECT COUNT(*) FROM email_verifications WHERE created_at > datetime('now','-1 day')`).Scan(&siteDay)
	if siteDay >= 50 {
		setFlash(w, "The site has sent too many verification emails today. Please try again tomorrow.")
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}

	token := newToken()
	app.db.Exec(`UPDATE users SET confirm_token=?, confirm_sent_at=CURRENT_TIMESTAMP WHERE id=?`, token, user.ID)
	app.db.Exec(`INSERT INTO email_verifications (user_id) VALUES (?)`, user.ID)
	go app.sendVerifyEmail(user.Email, token, app.verifyExpiryHours())

	setFlash(w, "A new verification email is on its way. Please check your inbox.")
	http.Redirect(w, r, back, http.StatusSeeOther)
}

// handleForgotPasswordForm shows the "email me a reset link" form. When email is
// disabled it explains the feature is unavailable rather than offering a form
// that can't work — email-dependent features never block, they degrade.
func (app *App) handleForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	if u, _ := app.getUserSession(r); u != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.renderPublic(w, r, "forgot_password.html", publicPage{
		Title: "Reset your password — " + getSiteName(),
		Flash: getFlash(w, r),
		Data:  map[string]any{"EmailEnabled": app.emailEnabled()},
	})
}

// handleForgotPasswordPost issues a single-use reset token and emails the link.
// To avoid revealing which addresses are registered it always reports the same
// generic message, whether or not the email matched an account.
func (app *App) handleForgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	const generic = "If an account with that email exists, we've sent a password reset link. Please check your inbox."

	if !app.emailEnabled() {
		setFlash(w, "Password reset by email isn't available right now. Please contact an admin for help getting back into your account.")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}
	if !validEmail(email) {
		setFlash(w, "Please enter a valid email address.")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	// Per-IP throttle: at most 3 reset requests per IP per day. Counts every
	// well-formed request (matched or not) so the limit can't be used to probe
	// which emails exist. The same generic message is shown when throttled.
	ipHash := hashIP(r)
	if blocked, _ := app.rlExceeded(86400, 3, ` AND action='password_reset' AND ip_hash=?`, ipHash); blocked {
		setFlash(w, generic)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.recordActionRateLimit("password_reset", ipHash)

	if u, err := userByEmail(app.db, email); err == nil {
		// Per-account throttle: at most 3 reset emails per account per day,
		// independent of the requesting IP. Keyed by account id in the rate-limit
		// log. Over-limit requests silently skip sending (still generic message).
		acctKey := fmt.Sprintf("acct:%d", u.ID)
		if blocked, _ := app.rlExceeded(86400, 3, ` AND action='pwreset_acct' AND ip_hash=?`, acctKey); !blocked {
			app.recordActionRateLimit("pwreset_acct", acctKey)
			token := newToken()
			app.db.Exec(`UPDATE users SET reset_token=?, reset_sent_at=CURRENT_TIMESTAMP WHERE id=?`, token, u.ID)
			go app.sendResetEmail(u.Email, token)
		}
	}

	setFlash(w, generic)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// handleResetPasswordForm validates the token from the emailed link and shows
// the choose-a-new-password form. Invalid/expired tokens go back to the request
// page with an explanation.
func (app *App) handleResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if _, ok := app.validResetUser(token); !ok {
		setFlash(w, "That password reset link is invalid or has expired. Please request a new one.")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}
	app.renderPublic(w, r, "reset_password.html", publicPage{
		Title: "Choose a new password — " + getSiteName(),
		Flash: getFlash(w, r),
		Data:  map[string]any{"Token": token},
	})
}

// handleResetPasswordPost sets the new password when the token is still valid,
// then clears the token and logs out all existing sessions for that account.
func (app *App) handleResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	token := strings.TrimSpace(r.FormValue("token"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	u, ok := app.validResetUser(token)
	if !ok {
		setFlash(w, "That password reset link is invalid or has expired. Please request a new one.")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}
	fail := func(msg string) {
		setFlash(w, msg)
		// Hex tokens are URL-safe, so no escaping is needed.
		http.Redirect(w, r, "/reset-password?token="+token, http.StatusSeeOther)
	}
	if len(password) < 8 {
		fail("Password must be at least 8 characters.")
		return
	}
	if len(password) > maxPasswordBytes {
		fail(fmt.Sprintf("Password must be no more than %d characters.", maxPasswordBytes))
		return
	}
	if password != confirm {
		fail("Passwords do not match.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	// A successful reset also clears any brute-force lock and the failure counter.
	app.db.Exec(`UPDATE users SET password_hash=?, reset_token='', reset_sent_at=NULL, failed_login_count=0, login_locked=0 WHERE id=?`, string(hash), u.ID)
	// Invalidate existing sessions so a reset also locks out anyone who was using
	// the old password.
	app.db.Exec(`DELETE FROM user_sessions WHERE user_id=?`, u.ID)

	setFlash(w, "Your password has been reset. Please log in with your new password.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// validResetUser returns the user for a non-empty, unexpired reset token.
func (app *App) validResetUser(token string) (User, bool) {
	if token == "" {
		return User{}, false
	}
	u, err := userByResetToken(app.db, token)
	if err != nil {
		return User{}, false
	}
	var expired int
	app.db.QueryRow(
		`SELECT CASE WHEN reset_sent_at IS NULL OR reset_sent_at < datetime('now', ?) THEN 1 ELSE 0 END FROM users WHERE id=?`,
		fmt.Sprintf("-%d hours", resetExpiryHours), u.ID,
	).Scan(&expired)
	if expired == 1 {
		return User{}, false
	}
	return u, true
}
