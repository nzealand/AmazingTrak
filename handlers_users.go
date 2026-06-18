package main

import (
	"database/sql"
	"net/http"
	"net/mail"
	"sort"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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
	if blocked, reason := app.checkActionRateLimit("register", ipHash, perHour, perDay); blocked {
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

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	// Email confirmation is optional: only issue a token when SMTP is configured
	// and the user supplied an address. Its absence never blocks registration.
	confirmToken := ""
	emailConfigured := app.smtpHost != "" && email != ""
	if emailConfigured {
		confirmToken = newToken()
	}

	res, err := app.db.Exec(
		`INSERT INTO users (username, email, password_hash, status, email_confirmed, confirm_token) VALUES (?, ?, ?, 'pending', 0, ?)`,
		username, email, string(hash), confirmToken,
	)
	if err != nil {
		fail("Could not create account: " + err.Error())
		return
	}
	userID, _ := res.LastInsertId()
	app.recordActionRateLimit("register", ipHash)

	if emailConfigured {
		go app.sendConfirmEmail(email, confirmToken, app.baseURL)
	}

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
	})
}

func (app *App) handleUserLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	ipHash := hashIP(r)
	if app.checkLoginThrottle(ipHash) {
		setFlash(w, "Too many failed attempts. Please wait 15 minutes.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	u, err := app.authenticateUser(username, password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if u == nil {
		app.recordLoginAttempt(ipHash, username, false)
		setFlash(w, "Invalid username or password.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.recordLoginAttempt(ipHash, username, true)
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
		ProfileUser *User
		Submissions []Suggestion
		Badges      []badge
		Comments    []Comment
	}
	app.renderPublic(w, r, "user_profile.html", publicPage{
		Title: profileUser.Username + " — " + getSiteName(),
		Flash: getFlash(w, r),
		Data: profileData{
			ProfileUser: &profileUser,
			Submissions: subs,
			Badges:      badgesForSubmissions(subs),
			Comments:    comments,
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
	// Confirm the email and clear the token. Bump pending → confirmed, but never
	// downgrade an already approved/auto-approved account.
	newStatus := u.Status
	if u.Status == "pending" {
		newStatus = "confirmed"
	}
	app.db.Exec(`UPDATE users SET email_confirmed=1, confirm_token='', status=? WHERE id=?`, newStatus, u.ID)
	setFlash(w, "Thanks — your email address is confirmed.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
