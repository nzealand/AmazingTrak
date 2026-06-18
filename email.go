package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func (app *App) sendMail(toEmail, subject, body string) {
	if app.smtpHost == "" || toEmail == "" {
		return
	}
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		toEmail, subject, body)
	addr := app.smtpHost + ":" + app.smtpPort
	var auth smtp.Auth
	if app.smtpUser != "" {
		auth = smtp.PlainAuth("", app.smtpUser, app.smtpPass, app.smtpHost)
	}
	from := app.smtpUser
	if from == "" {
		from = "noreply@amazingtrak.com"
	}
	to := strings.Split(toEmail, ",")
	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}
	smtp.SendMail(addr, auth, from, to, []byte(msg)) //nolint:errcheck — best-effort notification
}

// sendConfirmEmail emails a registered user a link to confirm their address.
// Best-effort: confirmation is optional and never blocks account use.
func (app *App) sendConfirmEmail(toEmail, token, baseURL string) {
	confirmURL := fmt.Sprintf("%s/confirm-email?token=%s", baseURL, token)
	body := fmt.Sprintf(`Welcome to %s!

Please confirm your email address by visiting the link below:

%s

If you didn't create an account, you can ignore this message.
`, getSiteName(), confirmURL)
	app.sendMail(toEmail, "Confirm your "+getSiteName()+" account", body)
}

func (app *App) sendSuggestionEmail(toEmail string, train Train, sug Suggestion, baseURL string) {
	if app.smtpHost == "" || toEmail == "" {
		return
	}

	approveURL := fmt.Sprintf("%s%s/suggestions/%d/approve", baseURL, app.adminPrefix, sug.ID)
	trainURL := fmt.Sprintf("%s/trains/%s", baseURL, train.Slug)
	suggestionsURL := fmt.Sprintf("%s%s/suggestions", baseURL, app.adminPrefix)

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
`,
		mediaLabel, train.DisplayName, train.CorridorName,
		sug.URL,
		commentLine,
		trainURL,
		approveURL,
		suggestionsURL,
	)

	subject := fmt.Sprintf("New %s submitted for %s", mediaLabel, train.DisplayName)
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		toEmail, subject, body)

	addr := app.smtpHost + ":" + app.smtpPort
	var auth smtp.Auth
	if app.smtpUser != "" {
		auth = smtp.PlainAuth("", app.smtpUser, app.smtpPass, app.smtpHost)
	}
	from := app.smtpUser
	if from == "" {
		from = "noreply@amazingtrak.com"
	}
	to := strings.Split(toEmail, ",")
	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}
	smtp.SendMail(addr, auth, from, to, []byte(msg)) //nolint:errcheck — best-effort notification
}
