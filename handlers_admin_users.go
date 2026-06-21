package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// markUserApproved promotes a user to the distinct "auto_approved" status when
// an admin endorses their content (adds a rarity, approves their video), unless
// the user is already manually approved.
func (app *App) markUserApproved(userID int64) {
	if prefs, err := getSitePrefs(app.db); err != nil || !prefs.AutoApproveOnVideo {
		return
	}
	app.db.Exec(
		`UPDATE users SET status='auto_approved' WHERE id=? AND status NOT IN ('approved','auto_approved')`,
		userID,
	)
}

// markUserSpammer flags a user as a spammer and removes all their comments.
func (app *App) markUserSpammer(userID int64) {
	app.db.Exec(`UPDATE users SET is_spammer=1 WHERE id=?`, userID)
	app.db.Exec(`DELETE FROM comments WHERE user_id=?`, userID)
}

// deleteUserWithSubmissionRule deletes a user, removing their non-rare media and
// suggestions (and the uploaded files) while preserving rare submissions, whose
// user_id is cleared by the ON DELETE SET NULL foreign keys.
func (app *App) deleteUserWithSubmissionRule(userID int64) error {
	// Collect this user's media first (single DB connection: fully drain before
	// issuing deletes).
	type mediaRow struct {
		id        int64
		tags      string
		localPath string
	}
	var medias []mediaRow
	rows, err := app.db.Query(`SELECT id, COALESCE(tags,''), COALESCE(local_path,'') FROM media WHERE user_id=?`, userID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var m mediaRow
		if err := rows.Scan(&m.id, &m.tags, &m.localPath); err != nil {
			rows.Close()
			return err
		}
		medias = append(medias, m)
	}
	rows.Close()

	type sugRow struct {
		id   int64
		tags string
	}
	var sugs []sugRow
	srows, err := app.db.Query(`SELECT id, COALESCE(tags,'') FROM suggestions WHERE user_id=?`, userID)
	if err != nil {
		return err
	}
	for srows.Next() {
		var s sugRow
		if err := srows.Scan(&s.id, &s.tags); err != nil {
			srows.Close()
			return err
		}
		sugs = append(sugs, s)
	}
	srows.Close()

	// Delete non-rare media (and their files); rare media is preserved.
	for _, m := range medias {
		if hasRarity(m.tags) {
			continue
		}
		if m.localPath != "" {
			deleteMediaFile(app.uploadsDir, m.localPath)
		}
		app.db.Exec(`DELETE FROM media WHERE id=?`, m.id)
	}
	// Delete non-rare suggestions; rare suggestions are preserved.
	for _, s := range sugs {
		if hasRarity(s.tags) {
			continue
		}
		app.db.Exec(`DELETE FROM suggestions WHERE id=?`, s.id)
	}

	// Removing the user clears user_id on any preserved (rare) rows via FK.
	_, err = app.db.Exec(`DELETE FROM users WHERE id=?`, userID)
	return err
}

func (app *App) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	users, err := allUsers(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	admins, err := allAdminUsers(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type usersData struct {
		Users  []User
		Admins []AdminUser
	}
	app.renderAdmin(w, r, "users.html", adminPage{
		Title:     "Users",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      usersData{Users: users, Admins: admins},
	})
}

func (app *App) handleAdminUserApprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	app.db.Exec(`UPDATE users SET status='approved' WHERE id=?`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "approve_user", "user", id, "")
	setFlash(w, "User approved.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminUserUnapprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	u, err := userByID(app.db, id)
	if err != nil {
		setFlash(w, "User not found.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	// Drop back to confirmed (if email was verified) or pending.
	newStatus := "pending"
	if u.EmailConfirmed {
		newStatus = "confirmed"
	}
	app.db.Exec(`UPDATE users SET status=? WHERE id=?`, newStatus, id)
	app.db.Exec(`UPDATE comments SET status='pending', reviewed_at=NULL WHERE user_id=? AND status='approved'`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "unapprove_user", "user", id, "")
	setFlash(w, "User moved back to "+newStatus+"; their approved comments set to pending.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminUserDelete(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := app.deleteUserWithSubmissionRule(id); err != nil {
		setFlash(w, "Error deleting user: "+err.Error())
	} else {
		s := sessionFromCtx(r)
		app.logAudit(s.AdminUserID, "delete_user", "user", id, "")
		setFlash(w, "User deleted; rare submissions preserved.")
	}
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminUsersDeleteUnapproved(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	// Collect unapproved user IDs first (single connection), then delete.
	var ids []int64
	rows, err := app.db.Query(`SELECT id FROM users WHERE status IN ('pending','confirmed')`)
	if err != nil {
		setFlash(w, "Database error.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	for rows.Next() {
		var id int64
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	rows.Close()
	count := 0
	for _, id := range ids {
		if app.deleteUserWithSubmissionRule(id) == nil {
			count++
		}
	}
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "delete_unapproved_users", "user", 0, fmt.Sprintf("%d deleted", count))
	setFlash(w, fmt.Sprintf("Deleted %d unapproved user(s); rare submissions preserved.", count))
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminAddUser(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	fail := func(msg string) {
		setFlash(w, msg)
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
	}
	if !validUsername(username) {
		fail("Username must be 3–30 characters (letters, digits, . _ -).")
		return
	}
	if email != "" && !validEmail(email) {
		fail("Invalid email address.")
		return
	}
	if len(password) < 8 {
		fail("Password must be at least 8 characters.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	_, err = app.db.Exec(
		`INSERT INTO users (username, email, password_hash, status, email_confirmed) VALUES (?, ?, ?, 'approved', 0)`,
		username, email, string(hash),
	)
	if err != nil {
		fail("Could not create user (username may be taken): " + err.Error())
		return
	}
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "add_user", "user", 0, username)
	setFlash(w, "User account created and approved.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminAddAdmin(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	level, _ := strconv.Atoi(r.FormValue("permission_level"))

	fail := func(msg string) {
		setFlash(w, msg)
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
	}
	if !validUsername(username) {
		fail("Admin username must be 3–30 characters (letters, digits, . _ -).")
		return
	}
	if len(password) < 8 {
		fail("Admin password must be at least 8 characters.")
		return
	}
	if level < 1 || level > 6 {
		fail("Choose a permission level from 1 to 6.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	_, err = app.db.Exec(
		`INSERT INTO admin_users (username, password_hash, must_change_password, permission_level) VALUES (?, ?, 0, ?)`,
		username, string(hash), level,
	)
	if err != nil {
		fail("Could not create admin (username may be taken): " + err.Error())
		return
	}
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "add_admin", "admin_user", 0, fmt.Sprintf("%s L%d", username, level))
	setFlash(w, "Admin account created.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminAdminDelete(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	target, err := getAdminUser(app.db, id)
	if err != nil {
		setFlash(w, "Admin not found.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	if target.IsCentral() {
		setFlash(w, "The central admin account cannot be deleted.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	// Deleting the account cascades its sessions (sessions FK ON DELETE CASCADE).
	app.db.Exec(`DELETE FROM admin_users WHERE id=?`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "delete_admin", "admin_user", id, target.Username)
	setFlash(w, "Admin account deleted.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminUserSubmissions(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	u, err := userByID(app.db, id)
	if err != nil {
		setFlash(w, "User not found.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	subs, err := submissionsByUserID(app.db, id)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type subsData struct {
		User        *User
		Submissions []Suggestion
	}
	app.renderAdmin(w, r, "user_submissions.html", adminPage{
		Title:     "Submissions — " + u.Username,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      subsData{User: &u, Submissions: subs},
	})
}

// handleAdminUserResetPassword sets a new password for a registered user.
// The admin types the new password directly in the admin panel — no email link
// is issued, so this works regardless of whether email is configured.
func (app *App) handleAdminUserResetPassword(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", 400)
		return
	}
	password := r.FormValue("new_password")
	if len(password) < 8 {
		setFlash(w, "Password must be at least 8 characters.")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	res, err := app.db.Exec(`UPDATE users SET password_hash=?, reset_token='', reset_sent_at=NULL WHERE id=?`, string(hash), id)
	if err != nil {
		setFlash(w, "Error updating password: "+err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.NotFound(w, r)
		return
	}
	// Invalidate all sessions so the new password takes effect immediately.
	app.db.Exec(`DELETE FROM user_sessions WHERE user_id=?`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "reset_user_password", "user", id, "")
	setFlash(w, "Password reset and all sessions invalidated.")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}

func (app *App) handleAdminUserAnonymize(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	u, err := userByID(app.db, id)
	if err != nil {
		setFlash(w, "User not found.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}
	if u.IsApproved() {
		setFlash(w, "Cannot anonymize an approved user.")
		http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
		return
	}

	// Pick the next anonymous[N] ID — find max existing.
	var maxN int
	app.db.QueryRow(`SELECT COALESCE(MAX(CAST(SUBSTR(username,10) AS INTEGER)),0) FROM users WHERE username LIKE 'anonymous%'`).Scan(&maxN)
	newUsername := fmt.Sprintf("anonymous%d", maxN+1)

	tx, err := app.db.Begin()
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer tx.Rollback()
	if _, err = tx.Exec(`UPDATE users SET username=?, email='', status='pending', confirm_token=NULL, email_confirmed=0 WHERE id=?`, newUsername, id); err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if _, err = tx.Exec(`UPDATE comments SET status='rejected', rejection_reason='anonymized' WHERE user_id=? AND status != 'rejected'`, id); err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if err = tx.Commit(); err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	app.logAudit(s.AdminUserID, "anonymize", "user", id, u.Username+" → "+newUsername)
	setFlash(w, "User anonymized as "+newUsername+".")
	http.Redirect(w, r, app.adminPrefix+"/users", http.StatusSeeOther)
}
