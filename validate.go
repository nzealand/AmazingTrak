package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// titleMatchesTrainNumber reports whether a video title contains the train's
// number as a standalone 3-digit token (so "Amtrak 742" matches train 742, but
// "1742" or "7420" does not). Only applies to exactly-3-digit train numbers.
func titleMatchesTrainNumber(title, trainNumber string) bool {
	if len(trainNumber) != 3 {
		return false
	}
	for _, r := range trainNumber {
		if r < '0' || r > '9' {
			return false
		}
	}
	runes := []rune(title)
	for i := 0; i+3 <= len(runes); i++ {
		if string(runes[i:i+3]) != trainNumber {
			continue
		}
		isDigit := func(r rune) bool { return r >= '0' && r <= '9' }
		if i > 0 && isDigit(runes[i-1]) {
			continue
		}
		if i+3 < len(runes) && isDigit(runes[i+3]) {
			continue
		}
		return true
	}
	return false
}

var allowedDomains = map[string]string{
	"youtube.com":           "video",
	"youtu.be":              "video",
	"vimeo.com":             "video",
	"flickr.com":            "image",
	"imgur.com":             "image",
	"railpictures.net":      "image",
	"rrpicturearchives.net": "image",
	"instagram.com":         "image",
	"commons.wikimedia.org": "image",
}

const (
	maxURLLen          = 500
	maxTitleLen        = 120
	maxCommentLen      = 2000
	defaultRatePerMin  = 1
	defaultRatePerHour = 5
	defaultRatePerDay  = 20
)

// classifyPublicURL validates and classifies a public suggestion URL.
// Returns domain, mediaType, normalizedURL, ok.
func classifyPublicURL(raw string) (domain, mediaType, normalized string, ok bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) > maxURLLen {
		return "", "", "", false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", "", false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", "", false
	}
	if u.User != nil {
		return "", "", "", false
	}
	host := strings.ToLower(u.Hostname())
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")
	mt, found := allowedDomains[host]
	if !found {
		return "", "", "", false
	}
	// Convert youtu.be short links to youtube.com watch URLs
	if host == "youtu.be" {
		vid := strings.TrimPrefix(u.Path, "/")
		if i := strings.IndexAny(vid, "/?&#"); i >= 0 {
			vid = vid[:i]
		}
		host = "youtube.com"
		u = &url.URL{Scheme: "https", Host: "www.youtube.com", Path: "/watch", RawQuery: "v=" + vid}
	}
	// Convert YouTube Shorts to standard watch URLs
	if host == "youtube.com" && strings.HasPrefix(u.Path, "/shorts/") {
		vid := strings.TrimPrefix(u.Path, "/shorts/")
		if i := strings.IndexAny(vid, "/?&#"); i >= 0 {
			vid = vid[:i]
		}
		u = &url.URL{Scheme: "https", Host: "www.youtube.com", Path: "/watch", RawQuery: "v=" + vid}
	}
	u.Fragment = ""
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = cleanQueryParams(u.Query(), host)
	return host, mt, u.String(), true
}

// cleanQueryParams strips unnecessary query parameters based on domain.
// For YouTube watch URLs, only the 'v' parameter is kept.
// For Vimeo, all query params are stripped (video ID is in path).
// For all others, common tracking params are removed.
func cleanQueryParams(q url.Values, host string) string {
	switch host {
	case "youtube.com":
		if vid := q.Get("v"); vid != "" {
			return "v=" + vid
		}
		// Non-watch URL (playlist, channel, etc): strip only tracking
	case "vimeo.com":
		return ""
	}
	for _, k := range []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content", "fbclid", "gclid", "ref"} {
		q.Del(k)
	}
	return q.Encode()
}

// validateAdminURL validates any http/https URL submitted by admin.
// Also normalizes YouTube URLs so they match the format from classifyPublicURL.
func validateAdminURL(raw string) (domain, normalized string, ok bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) > maxURLLen {
		return "", "", false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", false
	}
	if u.User != nil {
		return "", "", false
	}
	host := strings.ToLower(u.Hostname())
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")
	if host == "" {
		return "", "", false
	}
	// Normalize YouTube short links and Shorts
	if host == "youtu.be" {
		vid := strings.TrimPrefix(u.Path, "/")
		if i := strings.IndexAny(vid, "/?&#"); i >= 0 {
			vid = vid[:i]
		}
		host = "youtube.com"
		u = &url.URL{Scheme: "https", Host: "www.youtube.com", Path: "/watch", RawQuery: "v=" + vid}
	}
	if host == "youtube.com" && strings.HasPrefix(u.Path, "/shorts/") {
		vid := strings.TrimPrefix(u.Path, "/shorts/")
		if i := strings.IndexAny(vid, "/?&#"); i >= 0 {
			vid = vid[:i]
		}
		u = &url.URL{Scheme: "https", Host: "www.youtube.com", Path: "/watch", RawQuery: "v=" + vid}
	}
	if host == "youtube.com" {
		u.RawQuery = cleanQueryParams(u.Query(), host)
	}
	u.Fragment = ""
	u.Path = strings.TrimRight(u.Path, "/")
	return host, u.String(), true
}

// waitMessage renders a human "try again" message from how long until a slot
// frees up.
func waitMessage(d time.Duration) string {
	if d < time.Minute {
		s := int(d.Seconds()) + 1
		return fmt.Sprintf("Please wait %d second%s before trying again.", s, plural(s))
	}
	if d < time.Hour {
		m := int(d.Minutes()) + 1
		return fmt.Sprintf("Please wait %d minute%s before trying again.", m, plural(m))
	}
	h := int(d.Hours()) + 1
	return fmt.Sprintf("Please try again in about %d hour%s.", h, plural(h))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// rlExceeded reports whether the number of rate_limit_log rows matching `cond`
// within windowSecs is at/over limit, and how long until the oldest such row
// ages out of the window (i.e. when the caller may retry).
func (app *App) rlExceeded(windowSecs, limit int, cond string, args ...interface{}) (bool, time.Duration) {
	window := fmt.Sprintf("-%d seconds", windowSecs)
	qargs := append([]interface{}{window}, args...)
	var count int
	app.db.QueryRow(`SELECT COUNT(*) FROM rate_limit_log WHERE created_at > datetime('now', ?)`+cond, qargs...).Scan(&count)
	if count < limit {
		return false, 0
	}
	var ageSecs float64
	app.db.QueryRow(`SELECT COALESCE((julianday('now') - julianday(MIN(created_at))) * 86400, 0) FROM rate_limit_log WHERE created_at > datetime('now', ?)`+cond, qargs...).Scan(&ageSecs)
	wait := time.Duration(float64(windowSecs)-ageSecs) * time.Second
	if wait < time.Second {
		wait = time.Second
	}
	return true, wait
}

// checkSuggestRateLimit enforces the per-minute/hour/day caps for media
// suggestions (site-wide, action='suggest'), returning a wait-time message.
func (app *App) checkSuggestRateLimit(perMinute, perHour, perDay int) (blocked bool, reason string) {
	if b, d := app.rlExceeded(60, perMinute, ` AND action='suggest'`); b {
		return true, waitMessage(d)
	}
	if b, d := app.rlExceeded(3600, perHour, ` AND action='suggest'`); b {
		return true, waitMessage(d)
	}
	if b, d := app.rlExceeded(86400, perDay, ` AND action='suggest'`); b {
		return true, waitMessage(d)
	}
	return false, ""
}

func (app *App) recordRateLimit(ipHash string) {
	app.db.Exec(`INSERT INTO rate_limit_log (ip_hash, action) VALUES (?, 'suggest')`, ipHash)
}

// checkRegisterRateLimit enforces the site-wide new-account caps (counts all
// registrations across the site within the window, not per-IP).
func (app *App) checkRegisterRateLimit(perHour, perDay int) (blocked bool, reason string) {
	if b, d := app.rlExceeded(3600, perHour, ` AND action='register'`); b {
		return true, "Too many new accounts have been created across the site recently. " + waitMessage(d)
	}
	if b, d := app.rlExceeded(86400, perDay, ` AND action='register'`); b {
		return true, "The daily new-account limit for the site has been reached. " + waitMessage(d)
	}
	return false, ""
}

func (app *App) recordActionRateLimit(action, ipHash string) {
	app.db.Exec(`INSERT INTO rate_limit_log (ip_hash, action) VALUES (?, ?)`, ipHash, action)
}

// checkUserCommentRateLimit enforces per-user comment caps (hour + day), with a
// wait-time message derived from the comments table.
func (app *App) checkUserCommentRateLimit(userID int64, perHour, perDay int) (blocked bool, reason string) {
	if b, d := app.commentWindowExceeded(userID, 3600, perHour); b {
		return true, waitMessage(d)
	}
	if b, d := app.commentWindowExceeded(userID, 86400, perDay); b {
		return true, waitMessage(d)
	}
	return false, ""
}

func (app *App) commentWindowExceeded(userID int64, windowSecs, limit int) (bool, time.Duration) {
	window := fmt.Sprintf("-%d seconds", windowSecs)
	var count int
	app.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE user_id=? AND created_at > datetime('now', ?)`, userID, window).Scan(&count)
	if count < limit {
		return false, 0
	}
	var ageSecs float64
	app.db.QueryRow(`SELECT COALESCE((julianday('now') - julianday(MIN(created_at))) * 86400, 0) FROM comments WHERE user_id=? AND created_at > datetime('now', ?)`, userID, window).Scan(&ageSecs)
	wait := time.Duration(float64(windowSecs)-ageSecs) * time.Second
	if wait < time.Second {
		wait = time.Second
	}
	return true, wait
}

func (app *App) checkDuplicateSuggestion(trainID int64, normURL string) bool {
	var count int
	app.db.QueryRow(
		`SELECT (SELECT COUNT(*) FROM suggestions WHERE train_id=? AND url=? AND status='pending') +
		        (SELECT COUNT(*) FROM media WHERE train_id=? AND url=?)`,
		trainID, normURL, trainID, normURL,
	).Scan(&count)
	return count > 0
}

func setTimingCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "form_start",
		Value:    time.Now().Format(time.RFC3339),
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func checkTiming(r *http.Request) bool {
	c, err := r.Cookie("form_start")
	if err != nil {
		return false
	}
	t, err := time.Parse(time.RFC3339, c.Value)
	if err != nil {
		return false
	}
	return time.Since(t) >= 2*time.Second
}

var validVideoTags = map[string]bool{
	"long_consist": true, "doubleheader": true, "sandwich_set": true, "reverse_set": true,
	"blow_overs": true, "horn_show": true, "doppler": true,
	"scenic": true, "environment": true, "historic": true, "special_event": true,
}

// parseVideoTags validates and deduplicates submitted tag checkboxes.
func parseVideoTags(vals []string) string {
	var tags []string
	seen := map[string]bool{}
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if validVideoTags[v] && !seen[v] {
			tags = append(tags, v)
			seen[v] = true
		}
	}
	return strings.Join(tags, ",")
}

// sanitizeComment trims, normalizes line endings, strips control characters
// (other than newline/tab), and length-caps a comment body. It does NOT
// HTML-escape — like titles and captions, comment bodies are rendered through
// html/template, whose context-aware output escaping neutralizes any HTML the
// user typed (so "<b>" displays literally rather than rendering). This is how
// "no HTML in comments" is enforced: tags are shown as plain text, never run.
func sanitizeComment(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' || r >= ' ' {
			b.WriteRune(r)
		}
	}
	s = strings.TrimSpace(b.String())
	if len(s) > maxCommentLen {
		s = s[:maxCommentLen]
	}
	return s
}

// checkDuplicateComment guards against accidental double-posting: it reports
// whether this user already has an identical pending or approved comment on the
// same train.
func (app *App) checkDuplicateComment(trainID, userID int64, body string) bool {
	var count int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM comments WHERE train_id=? AND user_id=? AND body=? AND status IN ('pending','approved')`,
		trainID, userID, body,
	).Scan(&count)
	return count > 0
}

func (app *App) checkDuplicateCorridorComment(corridorID, userID int64, body string) bool {
	var count int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM comments WHERE corridor_id=? AND user_id=? AND body=? AND status IN ('pending','approved')`,
		corridorID, userID, body,
	).Scan(&count)
	return count > 0
}

// sanitizeTitle trims and length-caps a title. It does NOT HTML-escape:
// all titles are rendered through html/template, which performs
// context-aware escaping on output. Escaping here too would double-encode
// (e.g. "&" -> "&amp;" stored, then "&amp;amp;" rendered).
func sanitizeTitle(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > maxTitleLen {
		s = s[:maxTitleLen]
	}
	return s
}
