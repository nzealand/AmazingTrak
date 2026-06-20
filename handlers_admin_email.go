package main

import "net/http"

// emailErrorRow is one logged email failure.
type emailErrorRow struct {
	ID        int64
	To        string
	Subject   string
	Error     string
	CreatedAt string
}

// handleAdminEmailErrors lists recent email send failures (level 5). Useful when
// email is enabled but messages aren't arriving.
func (app *App) handleAdminEmailErrors(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	rows, err := app.db.Query(`SELECT id, to_addr, subject, error, created_at FROM email_errors ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()
	var errs []emailErrorRow
	for rows.Next() {
		var e emailErrorRow
		if err := rows.Scan(&e.ID, &e.To, &e.Subject, &e.Error, &e.CreatedAt); err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		errs = append(errs, e)
	}

	type emailErrorsData struct {
		Errors       []emailErrorRow
		EmailEnabled bool
	}
	app.renderAdmin(w, r, "email_errors.html", adminPage{
		Title:     "Email Errors",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      emailErrorsData{Errors: errs, EmailEnabled: app.emailEnabled()},
	})
}

// handleAdminEmailErrorsClear deletes all logged email errors.
func (app *App) handleAdminEmailErrorsClear(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	app.db.Exec(`DELETE FROM email_errors`)
	setFlash(w, "Email error log cleared.")
	http.Redirect(w, r, app.adminPrefix+"/email-errors", http.StatusSeeOther)
}
