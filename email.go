package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

// Email is entirely optional. Every send goes through sendMail, which is a
// silent no-op unless email is enabled (RESEND_API_KEY set + admin enabled +
// a sender address configured). Sends are best-effort; failures are recorded in
// the email_errors table for the admin to inspect, never surfaced to the user.

// emailEnabled reports whether outbound email is configured and turned on.
func (app *App) emailEnabled() bool {
	if app.resendKey == "" {
		return false
	}
	prefs, err := getSitePrefs(app.db)
	if err != nil {
		return false
	}
	return prefs.EmailEnabled && prefs.SenderEmail != ""
}

// logEmailError records a failed send so the admin can review it.
func (app *App) logEmailError(to, subject, errMsg string) {
	app.db.Exec(`INSERT INTO email_errors (to_addr, subject, error) VALUES (?, ?, ?)`, to, subject, errMsg)
}

// verifyExpiryHours returns the configured verification-link lifetime (default 24h).
func (app *App) verifyExpiryHours() int {
	if prefs, err := getSitePrefs(app.db); err == nil && prefs.VerifyExpiryHours > 0 {
		return prefs.VerifyExpiryHours
	}
	return 24
}

// sendMail delivers an email via the Resend HTTP API. Body is plain text; it is
// HTML-escaped and wrapped in a minimal HTML document. No-op when email is off.
func (app *App) sendMail(to, subject, body string) {
	if to == "" || !app.emailEnabled() {
		return
	}
	prefs, err := getSitePrefs(app.db)
	if err != nil {
		app.logEmailError(to, subject, "load prefs: "+err.Error())
		return
	}

	htmlBody := "<html><body style=\"font-family:sans-serif;line-height:1.5\">" +
		strings.ReplaceAll(html.EscapeString(body), "\n", "<br>") +
		"</body></html>"

	payload := map[string]interface{}{
		"from":    prefs.SenderEmail,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
		"text":    body,
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		app.logEmailError(to, subject, "marshal: "+err.Error())
		return
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(buf))
	if err != nil {
		app.logEmailError(to, subject, "request: "+err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer "+app.resendKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		app.logEmailError(to, subject, "send: "+err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		app.logEmailError(to, subject, fmt.Sprintf("resend HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody))))
	}
}

// sendVerifyEmail sends an email-address verification link that expires after the
// configured number of hours.
func (app *App) sendVerifyEmail(toEmail, token string, expiryHours int) {
	link := fmt.Sprintf("%s/confirm-email?token=%s", app.baseURL, token)
	body := fmt.Sprintf(`Welcome to %s!

Please confirm your email address by visiting the link below:

%s

This link expires in %d hours. If you didn't create an account, you can ignore this message.

Visit the site: %s
`, getSiteName(), link, expiryHours, app.baseURL)
	app.sendMail(toEmail, "Confirm your "+getSiteName()+" email address", body)
}

// sendConfirmEmail is retained for compatibility; it delegates to sendVerifyEmail
// using the configured expiry window.
func (app *App) sendConfirmEmail(toEmail, token, baseURL string) {
	hours := 24
	if prefs, err := getSitePrefs(app.db); err == nil && prefs.VerifyExpiryHours > 0 {
		hours = prefs.VerifyExpiryHours
	}
	app.sendVerifyEmail(toEmail, token, hours)
}

func (app *App) sendSuggestionEmail(toEmail string, train Train, sug Suggestion, baseURL string) {
	approveURL := fmt.Sprintf("%s%s/suggestions/%d/approve", app.baseURL, app.adminPrefix, sug.ID)
	trainURL := fmt.Sprintf("%s/trains/%s", app.baseURL, train.Slug)
	suggestionsURL := fmt.Sprintf("%s%s/suggestions", app.baseURL, app.adminPrefix)

	mediaLabel := "Photo/Video"
	switch sug.MediaType {
	case "video":
		mediaLabel = "Video"
	case "image":
		mediaLabel = "Photo"
	}
	commentLine := ""
	if sug.Caption != "" {
		commentLine = fmt.Sprintf("\nSubmitter comment: %s\n", sug.Caption)
	}
	body := fmt.Sprintf(`New %s submission for %s (%s)

Submitted URL: %s
%s
Train page: %s
Approve submission: %s
All pending submissions: %s
`, mediaLabel, train.DisplayName, train.CorridorName, sug.URL, commentLine, trainURL, approveURL, suggestionsURL)

	app.sendMail(toEmail, fmt.Sprintf("New %s submitted for %s", mediaLabel, train.DisplayName), body)
}

// sendConductorRequestEmail notifies admins that a user has requested the
// Conductor role for a corridor.
func (app *App) sendConductorRequestEmail(toEmail string, corridor Corridor, user User, baseURL string) {
	conductorsURL := fmt.Sprintf("%s%s/conductors", app.baseURL, app.adminPrefix)
	corridorURL := fmt.Sprintf("%s/corridors/%s", app.baseURL, corridor.Slug)
	body := fmt.Sprintf(`%s has requested to become the Conductor of %s.

Corridor page: %s
Review requests: %s
`, user.Username, corridor.Name, corridorURL, conductorsURL)
	app.sendMail(toEmail, fmt.Sprintf("Conductor request for %s", corridor.Name), body)
}

// sendThresholdEmail notifies the admin that the number of pending review items
// has crossed a notable threshold (1, 10, or 100).
func (app *App) sendThresholdEmail(toEmail string, threshold, total int) {
	adminURL := fmt.Sprintf("%s%s", app.baseURL, app.adminPrefix)
	body := fmt.Sprintf(`There are now %d items awaiting review on %s (crossed the %d threshold).

This includes pending suggestions, comments, new registrations, and conductor requests.

Review them here: %s
`, total, getSiteName(), threshold, adminURL)
	app.sendMail(toEmail, fmt.Sprintf("%s: %d items awaiting review", getSiteName(), total), body)
}
