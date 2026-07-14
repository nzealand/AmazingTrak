package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

// Vantage spots: public-submitted map pins marking a good place to see or
// photograph trains, each backed by a YouTube link or an uploaded photo.
// Submissions are moderated like suggestions/comments, and only approved
// spots are ever returned by the public API.

type VantageSpot struct {
	ID                 string
	Latitude           float64
	Longitude          float64
	Title              string
	Caption            string
	MediaType          string
	URL                string
	LocalPath          string
	OriginalFilename   string
	StoredFilename     string
	SourceDomain       string
	Status             string
	SubmitterIPHash    string
	SubmitterUserAgent string
	RejectionReason    string
	CreatedAt          string
	ReviewedAt         string
	ReviewedBy         sql.NullInt64
}

// PublicURL returns the browsable URL for the spot's media: the uploaded
// photo's site-relative path, or the linked video URL.
func (v VantageSpot) PublicURL() string {
	if v.LocalPath != "" {
		return "/uploads/" + strings.TrimPrefix(v.LocalPath, "/")
	}
	return v.URL
}

const vantageSpotColumns = `id, latitude, longitude, title, caption, media_type, url, local_path,
	original_filename, stored_filename, source_domain, status, submitter_ip_hash,
	submitter_user_agent, rejection_reason, created_at, reviewed_at, reviewed_by`

func scanVantageSpots(rows *sql.Rows) ([]VantageSpot, error) {
	var out []VantageSpot
	for rows.Next() {
		var v VantageSpot
		if err := rows.Scan(&v.ID, &v.Latitude, &v.Longitude, &v.Title, &v.Caption, &v.MediaType, &v.URL,
			&v.LocalPath, &v.OriginalFilename, &v.StoredFilename, &v.SourceDomain, &v.Status,
			&v.SubmitterIPHash, &v.SubmitterUserAgent, &v.RejectionReason, &v.CreatedAt, &v.ReviewedAt,
			&v.ReviewedBy); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func allVantageSpots(db *sql.DB, status string) ([]VantageSpot, error) {
	q := `SELECT ` + vantageSpotColumns + ` FROM vantage_spots`
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = db.Query(q+` WHERE status=? ORDER BY created_at DESC`, status)
	} else {
		rows, err = db.Query(q + ` ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVantageSpots(rows)
}

func vantageSpotByID(db *sql.DB, id string) (VantageSpot, error) {
	rows, err := db.Query(`SELECT `+vantageSpotColumns+` FROM vantage_spots WHERE id=?`, id)
	if err != nil {
		return VantageSpot{}, err
	}
	defer rows.Close()
	list, err := scanVantageSpots(rows)
	if err != nil {
		return VantageSpot{}, err
	}
	if len(list) == 0 {
		return VantageSpot{}, sql.ErrNoRows
	}
	return list[0], nil
}

// ----- Public: suggest a vantage spot -----

func (app *App) handleVantageSpotForm(w http.ResponseWriter, r *http.Request) {
	setTimingCookie(w)
	app.renderPublic(w, r, "vantage_spot_suggest.html", publicPage{
		Title: "Suggest a Vantage Spot",
		Flash: getFlash(w, r),
	})
}

func (app *App) handleVantageSpotSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		r.ParseForm()
	}

	// Honeypot fields — silently succeed without saving
	if r.FormValue("website") != "" || r.FormValue("a") != "" || r.FormValue("b") != "ok" {
		setFlash(w, "Thanks! Your vantage spot has been submitted for review.")
		http.Redirect(w, r, "/map", http.StatusSeeOther)
		return
	}

	if !checkTiming(r) {
		setFlash(w, "Please take a moment to fill out the form.")
		http.Redirect(w, r, "/vantage-spots/suggest", http.StatusSeeOther)
		return
	}

	ipHash := hashIP(r)
	currentUser, _ := app.getUserSession(r)
	conductor := currentUser != nil && isAnyConductor(app.db, currentUser.ID)
	if !conductor {
		prefs, _ := getSitePrefs(app.db)
		perMin := prefs.RatePerMinute
		if perMin <= 0 {
			perMin = defaultRatePerMin
		}
		perHour, perDay := prefs.RatePerHour, prefs.RatePerDay
		if currentUser != nil && currentUser.IsApproved() {
			perHour, perDay = prefs.TrustedRatePerHour, prefs.TrustedRatePerDay
		}
		if perHour <= 0 {
			perHour = defaultRatePerHour
		}
		if perDay <= 0 {
			perDay = defaultRatePerDay
		}
		if blocked, reason := app.checkSuggestRateLimit(perMin, perHour, perDay); blocked {
			setFlash(w, reason)
			http.Redirect(w, r, "/vantage-spots/suggest", http.StatusSeeOther)
			return
		}
	}

	renderFormError := func(msg string) {
		setTimingCookie(w)
		app.renderPublic(w, r, "vantage_spot_suggest.html", publicPage{
			Title: "Suggest a Vantage Spot",
			Flash: msg,
			FormValues: map[string]string{
				"title":   r.FormValue("title"),
				"caption": r.FormValue("caption"),
				"lat":     r.FormValue("lat"),
				"lon":     r.FormValue("lon"),
				"url":     r.FormValue("url"),
			},
		})
	}

	lat, lon := parseOptionalLatLon(r.FormValue("lat"), r.FormValue("lon"))
	if lat == nil || lon == nil {
		renderFormError("Please drop a pin on the map to mark the vantage spot.")
		return
	}

	title := sanitizeTitle(r.FormValue("title"))
	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > 500 {
		caption = caption[:500]
	}

	ua := r.UserAgent()
	if len(ua) > 200 {
		ua = ua[:200]
	}

	id := newToken()

	switch r.FormValue("media_kind") {
	case "youtube":
		domain, _, normURL, ok := classifyPublicURL(r.FormValue("url"))
		if !ok || domain != "youtube.com" {
			renderFormError("Please link a YouTube video showcasing this vantage spot.")
			return
		}
		_, err := app.db.Exec(
			`INSERT INTO vantage_spots (id, latitude, longitude, title, caption, media_type, url, source_domain, submitter_ip_hash, submitter_user_agent)
			 VALUES (?, ?, ?, ?, ?, 'video', ?, ?, ?, ?)`,
			id, *lat, *lon, title, caption, normURL, domain, ipHash, ua,
		)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
	case "photo":
		_, fh, err := r.FormFile("photo")
		if err != nil {
			renderFormError("Please choose a photo to upload, or link a YouTube video instead.")
			return
		}
		result, err := processUpload(fh, "vantage-"+id[:8], app.uploadsDir, app.db)
		if err != nil {
			renderFormError("Upload error: " + err.Error())
			return
		}
		_, err = app.db.Exec(
			`INSERT INTO vantage_spots (id, latitude, longitude, title, caption, media_type, local_path, stored_filename, original_filename, submitter_ip_hash, submitter_user_agent)
			 VALUES (?, ?, ?, ?, ?, 'image', ?, ?, ?, ?, ?)`,
			id, *lat, *lon, title, caption, result.LocalPath, result.StoredFilename, result.OriginalFilename, ipHash, ua,
		)
		if err != nil {
			deleteMediaFile(app.uploadsDir, result.LocalPath)
			http.Error(w, "Database error", 500)
			return
		}
	default:
		renderFormError("Please link a YouTube video or upload a photo showcasing this vantage spot.")
		return
	}

	app.recordRateLimit(ipHash)
	app.maybeNotifyPending()

	setFlash(w, "Thanks! Your vantage spot has been submitted for review.")
	http.Redirect(w, r, "/map", http.StatusSeeOther)
}

// ----- Public: approved spots for the map layer -----

type vantageSpotAPIItem struct {
	ID        string  `json:"id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Title     string  `json:"title"`
	Caption   string  `json:"caption"`
	MediaType string  `json:"mediaType"`
	URL       string  `json:"url,omitempty"`
	PhotoURL  string  `json:"photoUrl,omitempty"`
}

func (app *App) handleVantageSpotsAPI(w http.ResponseWriter, r *http.Request) {
	spots, err := allVantageSpots(app.db, "approved")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	out := make([]vantageSpotAPIItem, 0, len(spots))
	for _, sp := range spots {
		item := vantageSpotAPIItem{
			ID: sp.ID, Lat: sp.Latitude, Lon: sp.Longitude,
			Title: sp.Title, Caption: sp.Caption, MediaType: sp.MediaType,
		}
		if sp.MediaType == "video" {
			item.URL = sp.URL
		} else {
			item.PhotoURL = sp.PublicURL()
		}
		out = append(out, item)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	json.NewEncoder(w).Encode(map[string]interface{}{"spots": out})
}

// ----- Admin: review vantage spots -----

func (app *App) handleAdminVantageSpots(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	spots, err := allVantageSpots(app.db, status)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	s := sessionFromCtx(r)
	app.renderAdmin(w, r, "vantage_spots.html", adminPage{
		Title:     "Vantage Spots",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data: struct {
			Spots        []VantageSpot
			StatusFilter string
		}{spots, status},
	})
}

func (app *App) handleAdminVantageSpotApprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id := r.PathValue("id")
	s := sessionFromCtx(r)
	_, err := app.db.Exec(`UPDATE vantage_spots SET status='approved', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`, s.AdminUserID, id)
	if err != nil {
		setFlash(w, "Error approving: "+err.Error())
	} else {
		app.logAudit(s.AdminUserID, "approve_vantage_spot", "vantage_spot", 0, id)
		setFlash(w, "Vantage spot approved.")
	}
	http.Redirect(w, r, app.adminPrefix+"/vantage-spots", http.StatusSeeOther)
}

func (app *App) handleAdminVantageSpotReject(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id := r.PathValue("id")
	reason := strings.TrimSpace(r.FormValue("reason"))
	s := sessionFromCtx(r)
	_, err := app.db.Exec(`UPDATE vantage_spots SET status='rejected', rejection_reason=?, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`, reason, s.AdminUserID, id)
	if err != nil {
		setFlash(w, "Error rejecting: "+err.Error())
	} else {
		app.logAudit(s.AdminUserID, "reject_vantage_spot", "vantage_spot", 0, id)
		setFlash(w, "Vantage spot rejected.")
	}
	http.Redirect(w, r, app.adminPrefix+"/vantage-spots", http.StatusSeeOther)
}

func (app *App) handleAdminVantageSpotDelete(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id := r.PathValue("id")
	if spot, err := vantageSpotByID(app.db, id); err == nil {
		deleteMediaFile(app.uploadsDir, spot.LocalPath)
	}
	s := sessionFromCtx(r)
	_, err := app.db.Exec(`DELETE FROM vantage_spots WHERE id=?`, id)
	if err != nil {
		setFlash(w, "Error deleting: "+err.Error())
	} else {
		app.logAudit(s.AdminUserID, "delete_vantage_spot", "vantage_spot", 0, id)
		setFlash(w, "Vantage spot deleted.")
	}
	http.Redirect(w, r, app.adminPrefix+"/vantage-spots", http.StatusSeeOther)
}
