package main

import (
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

func (app *App) checkRateLimit(ipHash string, perMinute, perHour, perDay int) (blocked bool, reason string) {
	var minCount, hourCount, dayCount int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM rate_limit_log WHERE created_at > datetime('now', '-1 minute')`,
	).Scan(&minCount)
	if minCount >= perMinute {
		return true, "Please wait a moment before submitting again. Register and get approved to submit unlimited links."
	}
	app.db.QueryRow(
		`SELECT COUNT(*) FROM rate_limit_log WHERE created_at > datetime('now', '-1 hour')`,
	).Scan(&hourCount)
	if hourCount >= perHour {
		return true, "Too many submissions this hour. Register and get approved to submit unlimited links."
	}
	app.db.QueryRow(
		`SELECT COUNT(*) FROM rate_limit_log WHERE created_at > datetime('now', '-24 hours')`,
	).Scan(&dayCount)
	if dayCount >= perDay {
		return true, "Daily submission limit reached. Register and get approved to submit unlimited links."
	}
	return false, ""
}

func (app *App) recordRateLimit(ipHash string) {
	app.db.Exec(`INSERT INTO rate_limit_log (ip_hash) VALUES (?)`, ipHash)
}

func (app *App) checkActionRateLimit(action, ipHash string, perHour, perDay int) (blocked bool, reason string) {
	var hourCount, dayCount int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM rate_limit_log WHERE action=? AND ip_hash=? AND created_at > datetime('now', '-1 hour')`,
		action, ipHash,
	).Scan(&hourCount)
	if hourCount >= perHour {
		return true, "Too many registrations from this network. Please try again later."
	}
	app.db.QueryRow(
		`SELECT COUNT(*) FROM rate_limit_log WHERE action=? AND ip_hash=? AND created_at > datetime('now', '-24 hours')`,
		action, ipHash,
	).Scan(&dayCount)
	if dayCount >= perDay {
		return true, "Daily registration limit reached. Please try again tomorrow."
	}
	return false, ""
}

func (app *App) recordActionRateLimit(action, ipHash string) {
	app.db.Exec(`INSERT INTO rate_limit_log (ip_hash, action) VALUES (?, ?)`, ipHash, action)
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
