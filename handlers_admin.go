package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ----- Auth -----

func (app *App) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	s, _ := app.getSession(r)
	if s != nil {
		http.Redirect(w, r, app.adminPrefix, http.StatusSeeOther)
		return
	}
	app.renderAdmin(w, r, "login.html", adminPage{
		Title: "Admin Login",
		Flash: getFlash(w, r),
	})
}

func (app *App) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	ipHash := hashIP(r)
	if app.checkLoginThrottle(ipHash) {
		setFlash(w, "Too many failed attempts. Please wait 15 minutes.")
		http.Redirect(w, r, app.adminPrefix+"/login", http.StatusSeeOther)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	user, err := app.authenticateAdmin(username, password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if user == nil {
		app.recordLoginAttempt(ipHash, username, false)
		setFlash(w, "Invalid credentials.")
		http.Redirect(w, r, app.adminPrefix+"/login", http.StatusSeeOther)
		return
	}
	app.recordLoginAttempt(ipHash, username, true)
	app.db.Exec(`UPDATE admin_users SET last_login_at = CURRENT_TIMESTAMP WHERE id=?`, user.ID)
	session, err := app.createSession(app.db, user.ID, r)
	if err != nil {
		http.Error(w, "Session error", 500)
		return
	}
	app.setSessionCookie(w, session.ID)
	http.Redirect(w, r, app.adminPrefix, http.StatusSeeOther)
}

func (app *App) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	if s != nil {
		app.deleteSession(app.db, s.ID)
	}
	app.clearSessionCookie(w)
	http.Redirect(w, r, app.adminPrefix+"/login", http.StatusSeeOther)
}

// ----- Dashboard -----

func (app *App) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	// Recompute the pending-notify level so thresholds re-arm as the admin clears
	// the queue (downward crossings reset the stored level).
	app.maybeNotifyPending()

	var pending, trains, corridors, mediaCount, pendingComments, pendingRegistrations, pendingConductorReqs int
	app.db.QueryRow(`SELECT COUNT(*) FROM suggestions WHERE status='pending'`).Scan(&pending)
	app.db.QueryRow(`SELECT COUNT(*) FROM trains`).Scan(&trains)
	app.db.QueryRow(`SELECT COUNT(*) FROM corridors`).Scan(&corridors)
	app.db.QueryRow(`SELECT COUNT(*) FROM media`).Scan(&mediaCount)
	app.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE status='pending'`).Scan(&pendingComments)
	app.db.QueryRow(`SELECT COUNT(*) FROM users WHERE status IN ('pending','confirmed')`).Scan(&pendingRegistrations)
	app.db.QueryRow(`SELECT COUNT(*) FROM conductor_requests WHERE status='pending'`).Scan(&pendingConductorReqs)

	type dashData struct {
		PendingCount             int
		TrainCount               int
		CorridorCount            int
		MediaCount               int
		PendingCommentCount      int
		PendingRegistrationCount int
		PendingConductorCount    int
	}
	app.renderAdmin(w, r, "dashboard.html", adminPage{
		Title:     "Dashboard",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      dashData{pending, trains, corridors, mediaCount, pendingComments, pendingRegistrations, pendingConductorReqs},
	})
}

// ----- Corridors -----

func (app *App) handleAdminCorridors(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	corridors, err := allCorridors(app.db, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	app.renderAdmin(w, r, "corridors.html", adminPage{
		Title:     "Routes",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      corridors,
	})
}

func (app *App) handleAdminCorridorCreate(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		setFlash(w, "Name is required.")
		http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
		return
	}
	slug := slugify(name)
	region := strings.TrimSpace(r.FormValue("region"))
	res, err := app.db.Exec(
		`INSERT INTO corridors (name, slug, region) VALUES (?, ?, ?)`,
		name, slug, region,
	)
	if err != nil {
		setFlash(w, "Error creating route: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
		return
	}
	id, _ := res.LastInsertId()
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "create", "corridor", id, name)
	app.invalidateIndexCache()
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorEdit(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	corridor, err := corridorByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	trains, err := trainsByCorridorID(app.db, id, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	type editData struct {
		Corridor Corridor
		Trains   []Train
	}
	app.renderAdmin(w, r, "corridor_edit.html", adminPage{
		Title:     "Edit Route: " + corridor.Name,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      editData{Corridor: corridor, Trains: trains},
	})
}

func (app *App) handleAdminCorridorUpdate(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	slug := strings.TrimSpace(r.FormValue("slug"))
	region := strings.TrimSpace(r.FormValue("region"))
	description := strings.TrimSpace(r.FormValue("description"))
	sqSummary := strings.TrimSpace(r.FormValue("service_quality_summary"))
	scheduleURL := strings.TrimSpace(r.FormValue("schedule_url"))
	if scheduleURL != "" && !strings.HasPrefix(scheduleURL, "https://") && !strings.HasPrefix(scheduleURL, "http://") {
		scheduleURL = ""
	}

	if name == "" || slug == "" {
		setFlash(w, "Name and slug are required.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d", app.adminPrefix, id), http.StatusSeeOther)
		return
	}

	var onTimePct sql.NullFloat64
	if v := strings.TrimSpace(r.FormValue("on_time_percent")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 && f <= 100 {
			onTimePct = sql.NullFloat64{Float64: f, Valid: true}
		}
	}

	var heroTrainID sql.NullInt64
	if v := strings.TrimSpace(r.FormValue("hero_train_id")); v != "" && v != "0" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			heroTrainID = sql.NullInt64{Int64: n, Valid: true}
		}
	}

	_, err = app.db.Exec(
		`UPDATE corridors SET name=?, slug=?, region=?, description=?, on_time_percent=?,
		 service_quality_summary=?, hero_train_id=?, schedule_url=? WHERE id=?`,
		name, slug, region, description, onTimePct, sqSummary, heroTrainID, scheduleURL, id,
	)
	if err != nil {
		setFlash(w, "Error updating: "+err.Error())
	} else {
		s := sessionFromCtx(r)
		app.logAudit(s.AdminUserID, "update", "corridor", id, name)
		app.invalidateIndexCache()
		setFlash(w, "Route updated.")
	}
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorToggle(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	app.db.Exec(`UPDATE corridors SET is_active = CASE WHEN is_active=1 THEN 0 ELSE 1 END WHERE id=?`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "toggle", "corridor", id, "")
	app.invalidateIndexCache()
	http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
}

// ----- Corridor Media -----

func (app *App) handleAdminCorridorMedia(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	corridor, err := corridorByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	media, err := mediaByCorridorID(app.db, id, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	type mediaPageData struct {
		Corridor Corridor
		Media    []Media
	}
	app.renderAdmin(w, r, "corridor_media.html", adminPage{
		Title:     "Media: " + corridor.Name,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      mediaPageData{Corridor: corridor, Media: media},
	})
}

func (app *App) handleAdminCorridorMediaAdd(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	corridor, err := corridorByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		r.ParseForm()
	}

	sourceType := strings.TrimSpace(r.FormValue("source_type"))
	title := sanitizeTitle(r.FormValue("title"))
	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > 500 {
		caption = caption[:500]
	}
	locationName := strings.TrimSpace(r.FormValue("location_name"))
	lat, lon := parseOptionalLatLon(r.FormValue("lat"), r.FormValue("lon"))

	redirectURL := fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id)

	switch sourceType {
	case "url":
		rawURL := strings.TrimSpace(r.FormValue("url"))
		mediaType := r.FormValue("media_type_select")
		if mediaType == "" {
			mediaType = "website"
		}
		domain, normURL, ok := validateAdminURL(rawURL)
		if !ok {
			setFlash(w, "Invalid URL.")
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		if title == "" && domain == "youtube.com" {
			if t, ok := fetchYouTubeTitle(normURL); ok {
				title = sanitizeTitle(t)
			}
		}
		tags := ""
		if mediaType == "video" {
			tags = parseVideoTags(r.Form["tags"])
		}
		latN := toNullFloat64(lat)
		lonN := toNullFloat64(lon)
		_, err = app.db.Exec(
			`INSERT INTO media (corridor_id, media_type, source_type, url, title, caption, tags, source_domain,
			 latitude, longitude, location_name, location_source, is_published, added_by)
			 VALUES (?, ?, 'url', ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'admin')`,
			id, mediaType, normURL, title, caption, tags, domain, latN, lonN, locationName,
			locSourceStr(lat),
		)
		if err != nil {
			setFlash(w, "Error adding media: "+err.Error())
		} else {
			s := sessionFromCtx(r)
			app.logAudit(s.AdminUserID, "add_media", "corridor", id, normURL)
			setFlash(w, "Media added.")
		}

	case "upload", "paste":
		fileField := "file"
		if sourceType == "paste" {
			fileField = "paste_file"
		}
		_, fhHeader, err := r.FormFile(fileField)
		if err != nil {
			setFlash(w, "No file received.")
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		result, err := processUpload(fhHeader, corridor.Slug, app.uploadsDir, app.db)
		if err != nil {
			setFlash(w, "Upload error: "+err.Error())
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		// Manual lat/lon overrides EXIF if provided
		if lat != nil && lon != nil {
			result.Lat = lat
			result.Lon = lon
			result.LocationSource = "admin"
		}
		latN := toNullFloat64(result.Lat)
		lonN := toNullFloat64(result.Lon)
		res, err := app.db.Exec(
			`INSERT INTO media (corridor_id, media_type, source_type, local_path, stored_filename, original_filename,
			 title, caption, latitude, longitude, location_name, location_source, is_published, added_by)
			 VALUES (?, 'image', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'admin')`,
			id, sourceType, result.LocalPath, result.StoredFilename, result.OriginalFilename,
			title, caption, latN, lonN, locationName, result.LocationSource,
		)
		if err != nil {
			deleteMediaFile(app.uploadsDir, result.LocalPath)
			setFlash(w, "Error saving media: "+err.Error())
		} else {
			mid, _ := res.LastInsertId()
			s := sessionFromCtx(r)
			app.logAudit(s.AdminUserID, "upload_media", "corridor", mid, result.StoredFilename)
			msg := "Image uploaded."
			if result.LocationSource == "exif" {
				msg += fmt.Sprintf(" GPS extracted: %.4f, %.4f", *result.Lat, *result.Lon)
			}
			postAction := r.FormValue("post_action")
			switch postAction {
			case "set_hero":
				app.db.Exec(`UPDATE corridors SET hero_media_id=? WHERE id=?`, mid, id)
				msg += " Set as hero image."
			case "set_thumbnail":
				app.db.Exec(`UPDATE corridors SET thumbnail_media_id=? WHERE id=?`, mid, id)
				msg += " Set as thumbnail."
			}
			setFlash(w, msg)
		}

	default:
		setFlash(w, "Unknown source type.")
	}

	app.invalidateIndexCache()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaDelete(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)

	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	// Clear hero/thumbnail if this media was selected
	app.db.Exec(`UPDATE corridors SET hero_media_id=NULL WHERE hero_media_id=? AND id=?`, mid, id)
	app.db.Exec(`UPDATE corridors SET thumbnail_media_id=NULL WHERE thumbnail_media_id=? AND id=?`, mid, id)
	deleteMediaFile(app.uploadsDir, m.LocalPath)
	app.db.Exec(`DELETE FROM media WHERE id=?`, mid)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "delete_media", "corridor", mid, "")
	app.invalidateIndexCache()
	setFlash(w, "Media deleted.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaHero(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE corridors SET hero_media_id=? WHERE id=?`, mid, id)
	app.invalidateIndexCache()
	setFlash(w, "Hero image set.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaThumbnail(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE corridors SET thumbnail_media_id=? WHERE id=?`, mid, id)
	app.invalidateIndexCache()
	setFlash(w, "Thumbnail set.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaGeo(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	lat, lon := parseOptionalLatLon(r.FormValue("lat"), r.FormValue("lon"))
	locationName := strings.TrimSpace(r.FormValue("location_name"))
	latN := toNullFloat64(lat)
	lonN := toNullFloat64(lon)
	locSrc := "unknown"
	if lat != nil && lon != nil {
		locSrc = "admin"
	}
	app.db.Exec(`UPDATE media SET latitude=?, longitude=?, location_name=?, location_source=? WHERE id=?`,
		latN, lonN, locationName, locSrc, mid)
	setFlash(w, "Location updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaCaption(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > 500 {
		caption = caption[:500]
	}
	app.db.Exec(`UPDATE media SET caption=? WHERE id=?`, caption, mid)
	setFlash(w, "Caption updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminCorridorMediaTitle(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.CorridorID.Valid || m.CorridorID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	title := sanitizeTitle(r.FormValue("title"))
	app.db.Exec(`UPDATE media SET title=? WHERE id=?`, title, mid)
	app.invalidateIndexCache()
	setFlash(w, "Title updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/routes/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

// ----- Trains -----

func (app *App) handleAdminTrains(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	trains, err := allTrains(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	corridors, err := allCorridors(app.db, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type trainsData struct {
		Trains    []Train
		Corridors []Corridor
	}
	app.renderAdmin(w, r, "trains.html", adminPage{
		Title:     "Trains",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      trainsData{Trains: trains, Corridors: corridors},
	})
}

func (app *App) handleAdminTrainCreate(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	corridorID, err := strconv.ParseInt(r.FormValue("corridor_id"), 10, 64)
	if err != nil {
		setFlash(w, "Invalid route.")
		http.Redirect(w, r, app.adminPrefix+"/trains", http.StatusSeeOther)
		return
	}
	trainNumber := strings.TrimSpace(r.FormValue("train_number"))
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	if trainNumber == "" || displayName == "" {
		setFlash(w, "Train number and display name are required.")
		http.Redirect(w, r, app.adminPrefix+"/trains", http.StatusSeeOther)
		return
	}
	slug := slugify(displayName)
	res, err := app.db.Exec(
		`INSERT INTO trains (corridor_id, train_number, display_name, slug) VALUES (?, ?, ?, ?)`,
		corridorID, trainNumber, displayName, slug,
	)
	if err != nil {
		setFlash(w, "Error creating train: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/trains", http.StatusSeeOther)
		return
	}
	id, _ := res.LastInsertId()
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "create", "train", id, displayName)
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainDetail(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	train, err := trainByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	corridors, err := allCorridors(app.db, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	suggestions, err := suggestionsByTrainID(app.db, id, "pending")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	type detailData struct {
		Train       Train
		Corridors   []Corridor
		Suggestions []Suggestion
	}
	app.renderAdmin(w, r, "train_detail.html", adminPage{
		Title:     "Train: " + train.DisplayName,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      detailData{Train: train, Corridors: corridors, Suggestions: suggestions},
	})
}

func (app *App) handleAdminTrainUpdate(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	displayName := strings.TrimSpace(r.FormValue("display_name"))
	slug := strings.TrimSpace(r.FormValue("slug"))
	trainNumber := strings.TrimSpace(r.FormValue("train_number"))
	direction := strings.TrimSpace(r.FormValue("direction"))
	notes := strings.TrimSpace(r.FormValue("notes"))
	corridorID, _ := strconv.ParseInt(r.FormValue("corridor_id"), 10, 64)
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))

	if displayName == "" || slug == "" || trainNumber == "" {
		setFlash(w, "Display name, slug, and train number are required.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d", app.adminPrefix, id), http.StatusSeeOther)
		return
	}

	_, err = app.db.Exec(
		`UPDATE trains SET display_name=?, slug=?, train_number=?, direction=?, notes=?, corridor_id=?, sort_order=? WHERE id=?`,
		displayName, slug, trainNumber, direction, notes, corridorID, sortOrder, id,
	)
	if err != nil {
		setFlash(w, "Error updating: "+err.Error())
	} else {
		s := sessionFromCtx(r)
		app.logAudit(s.AdminUserID, "update", "train", id, displayName)
		setFlash(w, "Train updated.")
	}
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainToggle(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	app.db.Exec(`UPDATE trains SET is_active = CASE WHEN is_active=1 THEN 0 ELSE 1 END WHERE id=?`, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "toggle", "train", id, "")
	http.Redirect(w, r, app.adminPrefix+"/trains", http.StatusSeeOther)
}

// ----- Train Media -----

func (app *App) handleAdminTrainMedia(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	train, err := trainByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	media, err := mediaByTrainID(app.db, id, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	type mediaPageData struct {
		Train Train
		Media []Media
	}
	app.renderAdmin(w, r, "train_media.html", adminPage{
		Title:     "Media: " + train.DisplayName,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      mediaPageData{Train: train, Media: media},
	})
}

func (app *App) handleAdminTrainMediaAdd(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	train, err := trainByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		r.ParseForm()
	}

	sourceType := strings.TrimSpace(r.FormValue("source_type"))
	title := sanitizeTitle(r.FormValue("title"))
	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > 500 {
		caption = caption[:500]
	}
	locationName := strings.TrimSpace(r.FormValue("location_name"))
	lat, lon := parseOptionalLatLon(r.FormValue("lat"), r.FormValue("lon"))

	redirectURL := fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id)

	switch sourceType {
	case "url":
		rawURL := strings.TrimSpace(r.FormValue("url"))
		mediaType := r.FormValue("media_type_select")
		if mediaType == "" {
			mediaType = "website"
		}
		domain, normURL, ok := validateAdminURL(rawURL)
		if !ok {
			setFlash(w, "Invalid URL.")
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		if title == "" && domain == "youtube.com" {
			if t, ok := fetchYouTubeTitle(normURL); ok {
				title = sanitizeTitle(t)
			}
		}
		tags := ""
		if mediaType == "video" {
			tags = parseVideoTags(r.Form["tags"])
		}
		latN := toNullFloat64(lat)
		lonN := toNullFloat64(lon)
		_, err = app.db.Exec(
			`INSERT INTO media (train_id, media_type, source_type, url, title, caption, tags, source_domain,
			 latitude, longitude, location_name, location_source, is_published, added_by)
			 VALUES (?, ?, 'url', ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'admin')`,
			id, mediaType, normURL, title, caption, tags, domain, latN, lonN, locationName,
			locSourceStr(lat),
		)
		if err != nil {
			setFlash(w, "Error adding media: "+err.Error())
		} else {
			s := sessionFromCtx(r)
			app.logAudit(s.AdminUserID, "add_media", "train", id, normURL)
			setFlash(w, "Media added.")
		}

	case "upload", "paste":
		fileField := "file"
		if sourceType == "paste" {
			fileField = "paste_file"
		}
		_, fhHeader, err := r.FormFile(fileField)
		if err != nil {
			setFlash(w, "No file received.")
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		result, err := processUpload(fhHeader, train.Slug, app.uploadsDir, app.db)
		if err != nil {
			setFlash(w, "Upload error: "+err.Error())
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		if lat != nil && lon != nil {
			result.Lat = lat
			result.Lon = lon
			result.LocationSource = "admin"
		}
		latN := toNullFloat64(result.Lat)
		lonN := toNullFloat64(result.Lon)
		res, err := app.db.Exec(
			`INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, original_filename,
			 title, caption, latitude, longitude, location_name, location_source, is_published, added_by)
			 VALUES (?, 'image', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'admin')`,
			id, "paste", result.LocalPath, result.StoredFilename, result.OriginalFilename,
			title, caption, latN, lonN, locationName, result.LocationSource,
		)
		if err != nil {
			deleteMediaFile(app.uploadsDir, result.LocalPath)
			setFlash(w, "Error saving media: "+err.Error())
		} else {
			mid, _ := res.LastInsertId()
			s := sessionFromCtx(r)
			app.logAudit(s.AdminUserID, "upload_media", "train", mid, result.StoredFilename)
			msg := "Image uploaded."
			if result.LocationSource == "exif" {
				msg += fmt.Sprintf(" GPS extracted: %.4f, %.4f", *result.Lat, *result.Lon)
			}
			postAction := r.FormValue("post_action")
			switch postAction {
			case "set_hero":
				app.db.Exec(`UPDATE trains SET hero_media_id=? WHERE id=?`, mid, id)
				msg += " Set as hero image."
			case "set_thumbnail":
				app.db.Exec(`UPDATE trains SET thumbnail_media_id=? WHERE id=?`, mid, id)
				msg += " Set as thumbnail."
			}
			setFlash(w, msg)
		}

	default:
		setFlash(w, "Unknown source type.")
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaDelete(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE trains SET hero_media_id=NULL WHERE hero_media_id=? AND id=?`, mid, id)
	app.db.Exec(`UPDATE trains SET thumbnail_media_id=NULL WHERE thumbnail_media_id=? AND id=?`, mid, id)
	app.db.Exec(`UPDATE trains SET map_media_id=NULL WHERE map_media_id=? AND id=?`, mid, id)
	deleteMediaFile(app.uploadsDir, m.LocalPath)
	app.db.Exec(`DELETE FROM media WHERE id=?`, mid)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "delete_media", "train", mid, "")
	setFlash(w, "Media deleted.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaHero(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE trains SET hero_media_id=? WHERE id=?`, mid, id)
	setFlash(w, "Hero image set.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaThumbnail(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE trains SET thumbnail_media_id=? WHERE id=?`, mid, id)
	setFlash(w, "Thumbnail set.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaMap(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	app.db.Exec(`UPDATE trains SET map_media_id=? WHERE id=?`, mid, id)
	setFlash(w, "Map image set.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaGeo(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	lat, lon := parseOptionalLatLon(r.FormValue("lat"), r.FormValue("lon"))
	locationName := strings.TrimSpace(r.FormValue("location_name"))
	latN := toNullFloat64(lat)
	lonN := toNullFloat64(lon)
	locSrc := "unknown"
	if lat != nil && lon != nil {
		locSrc = "admin"
	}
	app.db.Exec(`UPDATE media SET latitude=?, longitude=?, location_name=?, location_source=? WHERE id=?`,
		latN, lonN, locationName, locSrc, mid)
	setFlash(w, "Location updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaTags(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	tags := parseVideoTags(r.Form["tags"])
	app.db.Exec(`UPDATE media SET tags=? WHERE id=?`, tags, mid)
	// Adding a rarity to a user-submitted link auto-approves that user.
	if hasRarity(tags) {
		var uid sql.NullInt64
		app.db.QueryRow(`SELECT user_id FROM media WHERE id=?`, mid).Scan(&uid)
		if uid.Valid {
			app.markUserApproved(uid.Int64)
		}
	}
	setFlash(w, "Tags updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaBest(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	if m.IsBest {
		app.db.Exec(`UPDATE media SET is_best=0 WHERE id=?`, mid)
		setFlash(w, "Best video cleared.")
	} else {
		app.db.Exec(`UPDATE media SET is_best=0 WHERE train_id=? AND media_type='video'`, id)
		app.db.Exec(`UPDATE media SET is_best=1 WHERE id=?`, mid)
		setFlash(w, "Best video set.")
	}
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaCaption(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > 500 {
		caption = caption[:500]
	}
	app.db.Exec(`UPDATE media SET caption=? WHERE id=?`, caption, mid)
	setFlash(w, "Caption updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

func (app *App) handleAdminTrainMediaTitle(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	mid, _ := strconv.ParseInt(r.PathValue("mid"), 10, 64)
	m, err := mediaByID(app.db, mid)
	if err != nil || !m.TrainID.Valid || m.TrainID.Int64 != id {
		setFlash(w, "Media not found.")
		http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
		return
	}
	title := sanitizeTitle(r.FormValue("title"))
	app.db.Exec(`UPDATE media SET title=? WHERE id=?`, title, mid)
	app.invalidateIndexCache()
	setFlash(w, "Title updated.")
	http.Redirect(w, r, fmt.Sprintf("%s/trains/%d/media", app.adminPrefix, id), http.StatusSeeOther)
}

// ----- Suggestions -----

func (app *App) handleAdminAllSuggestions(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	statusFilter := r.URL.Query().Get("status")
	suggestions, err := allSuggestions(app.db, statusFilter)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type sugData struct {
		Suggestions  []Suggestion
		StatusFilter string
	}
	app.renderAdmin(w, r, "all_suggestions.html", adminPage{
		Title:     "Suggestions",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      sugData{Suggestions: suggestions, StatusFilter: statusFilter},
	})
}

func (app *App) handleAdminSuggestionApprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	sug, err := suggestionByID(app.db, id)
	if err != nil {
		setFlash(w, "Suggestion not found.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	domain, normURL, ok := validateAdminURL(sug.URL)
	if !ok {
		normURL = sug.URL
		domain = sug.SourceDomain
	}
	// Carry the submitting user (if any) onto the published media row.
	var uid sql.NullInt64
	app.db.QueryRow(`SELECT user_id FROM suggestions WHERE id=?`, id).Scan(&uid)
	_, err = app.db.Exec(
		`INSERT INTO media (train_id, media_type, source_type, url, title, caption, tags, source_domain, location_source, is_published, added_by, user_id)
		 VALUES (?, ?, 'url', ?, ?, ?, ?, ?, 'unknown', 1, 'approved_suggestion', ?)`,
		sug.TrainID, sug.MediaType, normURL, sug.Title, sug.Caption, sug.Tags, domain, uid,
	)
	if err != nil {
		setFlash(w, "Error approving: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	s := sessionFromCtx(r)
	app.db.Exec(
		`UPDATE suggestions SET status='approved', is_spam=0, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
		s.AdminUserID, id,
	)
	// Approving a user-submitted video auto-approves that user.
	if uid.Valid && sug.MediaType == "video" {
		app.markUserApproved(uid.Int64)
	}
	// Auto-set thumbnail if the train doesn't have one yet.
	if sug.MediaType == "video" {
		go app.maybeAutoThumbnail(sug.TrainID, sug.URL, sug.SourceDomain)
	}
	app.logAudit(s.AdminUserID, "approve_suggestion", "suggestion", id, sug.URL)
	setFlash(w, "Suggestion approved and added to train media.")
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

// maybeAutoThumbnail fetches a thumbnail from a video URL and sets it on the
// train if the train currently has no thumbnail. Runs in a goroutine; errors
// are silent (the approval already succeeded at this point).
func (app *App) maybeAutoThumbnail(trainID int64, videoURL, sourceDomain string) {
	// Skip if the train already has a thumbnail.
	var existing sql.NullInt64
	if err := app.db.QueryRow(`SELECT thumbnail_media_id FROM trains WHERE id=?`, trainID).Scan(&existing); err != nil || existing.Valid {
		return
	}

	var thumbURL string
	switch sourceDomain {
	case "youtube.com", "youtu.be":
		vid := ""
		if idx := strings.Index(videoURL, "v="); idx >= 0 {
			vid = videoURL[idx+2:]
			if end := strings.Index(vid, "&"); end >= 0 {
				vid = vid[:end]
			}
		}
		if vid == "" {
			return
		}
		// hqdefault (480×360) is universally available; maxresdefault only exists
		// for videos that have been watched enough to generate it.
		thumbURL = "https://img.youtube.com/vi/" + vid + "/hqdefault.jpg"
	default:
		return
	}

	res, err := app.db.Exec(
		`INSERT INTO media (train_id, media_type, source_type, url, title, source_domain, location_source, is_published, added_by)
		 VALUES (?, 'image', 'url', ?, 'Auto thumbnail', 'img.youtube.com', 'unknown', 1, 'auto_thumbnail')`,
		trainID, thumbURL,
	)
	if err != nil {
		return
	}
	mid, _ := res.LastInsertId()
	app.db.Exec(`UPDATE trains SET thumbnail_media_id=? WHERE id=?`, mid, trainID)
}

func (app *App) handleAdminSuggestionReject(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	reason := strings.TrimSpace(r.FormValue("reason"))
	s := sessionFromCtx(r)
	app.db.Exec(
		`UPDATE suggestions SET status='rejected', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=?, rejection_reason=? WHERE id=?`,
		s.AdminUserID, reason, id,
	)
	app.logAudit(s.AdminUserID, "reject_suggestion", "suggestion", id, reason)
	setFlash(w, "Suggestion rejected.")
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

// handleAdminSuggestionMarkSpam marks a submission as spam and applies cascading
// penalties to the submitting regular user (if any).
//
// No associated user (anonymous / admin-submitted):
//   - Just flag this one submission as pending-spam.
//
// User with no previously-approved submissions:
//   - Mark user as spammer.
//   - Set ALL their submissions to rejected-spam.
//
// User with at least one approved submission (trusted contributor who may have
// had one bad submission sneak through):
//   - Mark user as spammer.
//   - Set all their submissions that lack rarity tags to pending-spam.
//   - Leave rarity-tagged submissions untouched (already verified to be genuine).
func (app *App) handleAdminSuggestionMarkSpam(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	s := sessionFromCtx(r)

	var userID sql.NullInt64
	app.db.QueryRow(`SELECT user_id FROM suggestions WHERE id=?`, id).Scan(&userID)

	if !userID.Valid {
		// Anonymous or admin-authored: only flag this one submission.
		app.db.Exec(
			`UPDATE suggestions SET status='pending', is_spam=1, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
			s.AdminUserID, id,
		)
		app.logAudit(s.AdminUserID, "spam_suggestion", "suggestion", id, "no user")
		setFlash(w, "Submission marked as spam.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}

	uid := userID.Int64
	app.markUserSpammer(uid)

	// Count previously-approved submissions (not counting the current one).
	var approvedCount int
	app.db.QueryRow(
		`SELECT COUNT(*) FROM suggestions WHERE user_id=? AND status='approved'`,
		uid,
	).Scan(&approvedCount)

	if approvedCount == 0 {
		// No approved content: mass-reject everything as spam.
		app.db.Exec(
			`UPDATE suggestions SET status='rejected', is_spam=1, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE user_id=?`,
			s.AdminUserID, uid,
		)
		app.logAudit(s.AdminUserID, "spam_suggestion", "suggestion", id, fmt.Sprintf("user %d — all rejected", uid))
		setFlash(w, "Submission marked as spam. User marked as spammer and all their submissions rejected.")
	} else {
		// Has approved content: only demote non-rarity submissions to pending-spam.
		// Fetch all their submission IDs and tags (drain rows before exec).
		type idTags struct {
			id   int64
			tags string
		}
		rows, err := app.db.Query(`SELECT id, COALESCE(tags,'') FROM suggestions WHERE user_id=?`, uid)
		var subs []idTags
		if err == nil {
			for rows.Next() {
				var it idTags
				rows.Scan(&it.id, &it.tags)
				subs = append(subs, it)
			}
			rows.Close()
		}
		for _, sub := range subs {
			if !hasRarity(sub.tags) {
				app.db.Exec(
					`UPDATE suggestions SET status='pending', is_spam=1, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
					s.AdminUserID, sub.id,
				)
			}
		}
		app.logAudit(s.AdminUserID, "spam_suggestion", "suggestion", id, fmt.Sprintf("user %d — non-rarity pending-spam", uid))
		setFlash(w, "Submission marked as spam. User marked as spammer and non-rarity submissions set to pending-spam.")
	}

	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

// handleAdminSuggestionUnapprove reverts an approved suggestion back to pending
// and removes the media row that approval created (matched by train + URL).
func (app *App) handleAdminSuggestionUnapprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	sug, err := suggestionByID(app.db, id)
	if err != nil {
		setFlash(w, "Suggestion not found.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	if sug.Status != "approved" {
		setFlash(w, "Only approved suggestions can be unapproved.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	// Remove the media that approval added. URL may have been normalized on
	// insert, so match on either the stored or normalized form.
	normURL := sug.URL
	if _, nu, ok := validateAdminURL(sug.URL); ok {
		normURL = nu
	}
	app.db.Exec(
		`DELETE FROM media WHERE train_id=? AND added_by='approved_suggestion' AND url IN (?, ?)`,
		sug.TrainID, sug.URL, normURL,
	)
	app.db.Exec(
		`UPDATE suggestions SET status='pending', reviewed_at=NULL, reviewed_by=NULL, auto_approved=0 WHERE id=?`,
		id,
	)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "unapprove_suggestion", "suggestion", id, sug.URL)
	setFlash(w, "Suggestion unapproved and returned to pending.")
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

// handleAdminSuggestionEdit updates the title and caption of a pending suggestion.
func (app *App) handleAdminSuggestionEdit(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	sug, err := suggestionByID(app.db, id)
	if err != nil {
		setFlash(w, "Suggestion not found.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	if sug.Status != "pending" {
		setFlash(w, "Only pending suggestions can be edited.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	title := sanitizeTitle(strings.TrimSpace(r.FormValue("title")))
	caption := strings.TrimSpace(r.FormValue("caption"))

	// The edit form only exposes rarity checkboxes; preserve any existing
	// non-rarity (highlight) tags rather than wiping them out.
	var keep []string
	for _, t := range strings.Split(sug.Tags, ",") {
		t = strings.TrimSpace(t)
		if t != "" && !rarityTags[t] {
			keep = append(keep, t)
		}
	}
	tags := parseVideoTags(append(keep, r.Form["tags"]...))
	app.db.Exec(`UPDATE suggestions SET title=?, caption=?, tags=? WHERE id=?`, title, caption, tags, id)
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "edit_suggestion", "suggestion", id, title)
	setFlash(w, "Suggestion updated.")
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

func (app *App) handleAdminSuggestionApproveAll(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	rows, err := app.db.Query(
		`SELECT id, train_id, url, title, caption, tags, media_type, source_domain, user_id FROM suggestions WHERE status='pending'`,
	)
	if err != nil {
		setFlash(w, "Database error.")
		http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
		return
	}
	type pendingSug struct {
		Suggestion
		userID sql.NullInt64
	}
	var pending []pendingSug
	for rows.Next() {
		var p pendingSug
		if err := rows.Scan(&p.ID, &p.TrainID, &p.URL, &p.Title, &p.Caption, &p.Tags, &p.MediaType, &p.SourceDomain, &p.userID); err != nil {
			continue
		}
		pending = append(pending, p)
	}
	rows.Close()

	// The DB connection pool is limited to a single connection (see openDB);
	// rows must be fully drained and closed above before issuing the Exec
	// calls below, or those Execs would block forever waiting on a
	// connection that's still held by this open query.
	s := sessionFromCtx(r)
	var count int
	for _, sug := range pending {
		domain, normURL, ok := validateAdminURL(sug.URL)
		if !ok {
			normURL = sug.URL
			domain = sug.SourceDomain
		}
		_, err := app.db.Exec(
			`INSERT INTO media (train_id, media_type, source_type, url, title, caption, tags, source_domain, location_source, is_published, added_by, user_id)
			 VALUES (?, ?, 'url', ?, ?, ?, ?, ?, 'unknown', 1, 'approved_suggestion', ?)`,
			sug.TrainID, sug.MediaType, normURL, sug.Title, sug.Caption, sug.Tags, domain, sug.userID,
		)
		if err != nil {
			continue
		}
		app.db.Exec(
			`UPDATE suggestions SET status='approved', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
			s.AdminUserID, sug.ID,
		)
		if sug.userID.Valid && sug.MediaType == "video" {
			app.markUserApproved(sug.userID.Int64)
		}
		if sug.MediaType == "video" {
			go app.maybeAutoThumbnail(sug.TrainID, sug.URL, sug.SourceDomain)
		}
		count++
	}
	app.logAudit(s.AdminUserID, "approve_all_suggestions", "suggestion", 0, fmt.Sprintf("%d approved", count))
	setFlash(w, fmt.Sprintf("Approved %d pending submission(s).", count))
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

func (app *App) handleAdminSuggestionRejectAll(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	s := sessionFromCtx(r)
	res, err := app.db.Exec(
		`UPDATE suggestions SET status='rejected', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE status='pending'`,
		s.AdminUserID,
	)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		app.logAudit(s.AdminUserID, "reject_all_suggestions", "suggestion", 0, fmt.Sprintf("%d rejected", n))
		setFlash(w, fmt.Sprintf("Rejected %d pending submission(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

func (app *App) handleAdminSuggestionDeleteRejected(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	s := sessionFromCtx(r)
	res, err := app.db.Exec(`DELETE FROM suggestions WHERE status='rejected'`)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		app.logAudit(s.AdminUserID, "delete_rejected_suggestions", "suggestion", 0, fmt.Sprintf("%d deleted", n))
		setFlash(w, fmt.Sprintf("Deleted %d rejected submission(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/suggestions", http.StatusSeeOther)
}

// ----- Comments -----

func (app *App) handleAdminComments(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	statusFilter := r.URL.Query().Get("status")
	comments, err := allComments(app.db, statusFilter)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type commentData struct {
		Comments     []Comment
		StatusFilter string
	}
	app.renderAdmin(w, r, "comments.html", adminPage{
		Title:     "Comments",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      commentData{Comments: comments, StatusFilter: statusFilter},
	})
}

func (app *App) handleAdminCommentApprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	s := sessionFromCtx(r)
	res, err := app.db.Exec(
		`UPDATE comments SET status='approved', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
		s.AdminUserID, id,
	)
	if err != nil {
		setFlash(w, "Error approving comment: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		setFlash(w, "Comment not found.")
		http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
		return
	}
	app.logAudit(s.AdminUserID, "approve_comment", "comment", id, "")
	setFlash(w, "Comment approved and published.")
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

func (app *App) handleAdminCommentReject(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	reason := strings.TrimSpace(r.FormValue("reason"))
	s := sessionFromCtx(r)
	app.db.Exec(
		`UPDATE comments SET status='rejected', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=?, rejection_reason=? WHERE id=?`,
		s.AdminUserID, reason, id,
	)
	app.logAudit(s.AdminUserID, "reject_comment", "comment", id, reason)
	setFlash(w, "Comment rejected.")
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

// handleAdminCommentUnapprove reverts an approved comment back to pending,
// hiding it from the public train page again.
func (app *App) handleAdminCommentUnapprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	s := sessionFromCtx(r)
	app.db.Exec(
		`UPDATE comments SET status='pending', reviewed_at=NULL, reviewed_by=NULL WHERE id=? AND status='approved'`,
		id,
	)
	app.logAudit(s.AdminUserID, "unapprove_comment", "comment", id, "")
	setFlash(w, "Comment returned to pending.")
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

func (app *App) handleAdminCommentApproveAll(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	s := sessionFromCtx(r)
	res, err := app.db.Exec(
		`UPDATE comments SET status='approved', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE status='pending'`,
		s.AdminUserID,
	)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		app.logAudit(s.AdminUserID, "approve_all_comments", "comment", 0, fmt.Sprintf("%d approved", n))
		setFlash(w, fmt.Sprintf("Approved %d pending comment(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

func (app *App) handleAdminCommentRejectAll(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	s := sessionFromCtx(r)
	res, err := app.db.Exec(
		`UPDATE comments SET status='rejected', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE status='pending'`,
		s.AdminUserID,
	)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		app.logAudit(s.AdminUserID, "reject_all_comments", "comment", 0, fmt.Sprintf("%d rejected", n))
		setFlash(w, fmt.Sprintf("Rejected %d pending comment(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

func (app *App) handleAdminCommentDeleteRejected(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	s := sessionFromCtx(r)
	res, err := app.db.Exec(`DELETE FROM comments WHERE status='rejected'`)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		app.logAudit(s.AdminUserID, "delete_rejected_comments", "comment", 0, fmt.Sprintf("%d deleted", n))
		setFlash(w, fmt.Sprintf("Deleted %d rejected comment(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/comments", http.StatusSeeOther)
}

func (app *App) handleAdminTrainDeleteInactive(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	// Collect uploaded file paths for deactivated trains before deleting
	rows, _ := app.db.Query(
		`SELECT COALESCE(m.local_path,'') FROM media m JOIN trains t ON t.id=m.train_id WHERE t.is_active=0 AND m.local_path != ''`,
	)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var p string
			if rows.Scan(&p) == nil && p != "" {
				deleteMediaFile(app.uploadsDir, p)
			}
		}
	}
	res, err := app.db.Exec(`DELETE FROM trains WHERE is_active=0`)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		s := sessionFromCtx(r)
		app.logAudit(s.AdminUserID, "delete_inactive_trains", "train", 0, fmt.Sprintf("%d deleted", n))
		setFlash(w, fmt.Sprintf("Deleted %d deactivated train(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/trains", http.StatusSeeOther)
}

func (app *App) handleAdminCorridorDeleteInactive(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	// Collect uploaded file paths for deactivated corridor trains
	rows, _ := app.db.Query(
		`SELECT COALESCE(m.local_path,'') FROM media m JOIN trains t ON t.id=m.train_id JOIN corridors c ON c.id=t.corridor_id WHERE c.is_active=0 AND m.local_path != ''`,
	)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var p string
			if rows.Scan(&p) == nil && p != "" {
				deleteMediaFile(app.uploadsDir, p)
			}
		}
	}
	// Corridor-level uploaded media
	rows2, _ := app.db.Query(
		`SELECT COALESCE(m.local_path,'') FROM media m JOIN corridors c ON c.id=m.corridor_id WHERE c.is_active=0 AND m.local_path != ''`,
	)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var p string
			if rows2.Scan(&p) == nil && p != "" {
				deleteMediaFile(app.uploadsDir, p)
			}
		}
	}
	res, err := app.db.Exec(`DELETE FROM corridors WHERE is_active=0`)
	if err != nil {
		setFlash(w, "Error: "+err.Error())
	} else {
		n, _ := res.RowsAffected()
		s := sessionFromCtx(r)
		app.logAudit(s.AdminUserID, "delete_inactive_corridors", "corridor", 0, fmt.Sprintf("%d deleted", n))
		app.invalidateIndexCache()
		setFlash(w, fmt.Sprintf("Deleted %d deactivated route(s).", n))
	}
	http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
}

func (app *App) handleAdminSyncScheduleURLs(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}

	amtrakRoutesMu.Lock()
	data := amtrakRoutesJSON
	amtrakRoutesMu.Unlock()

	if data == nil {
		var err error
		data, err = fetchAmtrakRoutes()
		if err != nil {
			setFlash(w, "Could not fetch Amtrak route data: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
			return
		}
		amtrakRoutesMu.Lock()
		amtrakRoutesJSON = data
		amtrakRoutesFetched = time.Now()
		amtrakRoutesMu.Unlock()
	}

	var fc struct {
		Features []struct {
			Properties struct {
				Name     string `json:"name"`
				RouteURL string `json:"ROUTE_URL"`
			} `json:"properties"`
		} `json:"features"`
	}
	if err := json.Unmarshal(data, &fc); err != nil {
		setFlash(w, "Could not parse Amtrak route data.")
		http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
		return
	}

	corridors, err := allCorridors(app.db, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	// Match by normalized slug: strip any leading "amtrak-" from both the
	// corridor slug and the slugified route name so "Texas Eagle" matches
	// the "amtrak-texas-eagle" corridor.
	normSlug := func(s string) string { return strings.TrimPrefix(s, "amtrak-") }
	normToID := make(map[string]int64, len(corridors))
	for _, c := range corridors {
		normToID[normSlug(c.Slug)] = c.ID
	}

	updated := 0
	for _, f := range fc.Features {
		if f.Properties.RouteURL == "" {
			continue
		}
		norm := normSlug(slugify(f.Properties.Name))
		if id, ok := normToID[norm]; ok {
			res, _ := app.db.Exec(
				`UPDATE corridors SET schedule_url=? WHERE id=? AND (schedule_url='' OR schedule_url IS NULL)`,
				f.Properties.RouteURL, id,
			)
			if n, _ := res.RowsAffected(); n > 0 {
				updated++
			}
		}
	}

	if updated > 0 {
		app.invalidateIndexCache()
	}
	s := sessionFromCtx(r)
	app.logAudit(s.AdminUserID, "sync_schedule_urls", "corridor", 0, fmt.Sprintf("%d updated", updated))
	setFlash(w, fmt.Sprintf("Synced Amtrak schedule URLs — %d route(s) updated.", updated))
	http.Redirect(w, r, app.adminPrefix+"/routes", http.StatusSeeOther)
}

// locSourceStr returns "admin" if lat is set, else "unknown".
func locSourceStr(lat *float64) string {
	if lat != nil {
		return "admin"
	}
	return "unknown"
}

// ----- Settings -----

func (app *App) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	user, err := getAdminUser(app.db, s.AdminUserID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	prefs, err := getSitePrefs(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type settingsData struct {
		User         AdminUser
		Prefs        SitePreferences
		HasResendKey bool
	}
	app.renderAdmin(w, r, "settings.html", adminPage{
		Title:     "Settings",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      settingsData{User: user, Prefs: prefs, HasResendKey: app.resendKey != ""},
	})
}

func (app *App) handleAdminSettingsPost(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil && err != http.ErrNotMultipart {
		http.Error(w, "Bad request", 400)
		return
	}

	s := sessionFromCtx(r)
	action := r.FormValue("action")

	switch action {
	case "site_branding":
		siteName := strings.TrimSpace(r.FormValue("site_name"))
		if siteName == "" {
			setFlash(w, "Site name cannot be empty.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		if len(siteName) > 60 {
			siteName = siteName[:60]
		}

		var faviconPath string
		app.db.QueryRow(`SELECT COALESCE(favicon_path,'') FROM site_preferences WHERE id=1`).Scan(&faviconPath)

		if file, header, ferr := r.FormFile("favicon"); ferr == nil {
			defer file.Close()
			newPath, perr := processFaviconUpload(header, app.uploadsDir)
			if perr != nil {
				setFlash(w, "Error uploading favicon: "+perr.Error())
				http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
				return
			}
			faviconPath = newPath
			setFavicon(faviconPath, fmt.Sprintf("%d", time.Now().Unix()))
		} else if ferr != http.ErrMissingFile {
			setFlash(w, "Error reading favicon upload: "+ferr.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}

		if _, err := app.db.Exec(`UPDATE site_preferences SET site_name=?, favicon_path=? WHERE id=1`, siteName, faviconPath); err != nil {
			setFlash(w, "Error saving site branding: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		setSiteName(siteName)
		app.invalidateIndexCache()
		app.logAudit(s.AdminUserID, "update_site_branding", "site_preferences", 1, "")
		setFlash(w, "Site branding saved.")

	case "credentials":
		newUsername := strings.TrimSpace(r.FormValue("username"))
		newPassword := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		currentPassword := r.FormValue("current_password")

		// Re-authenticate before any credential change (defence against an
		// unattended/hijacked session silently taking over the account).
		current, cerr := getAdminUser(app.db, s.AdminUserID)
		if cerr != nil {
			http.Error(w, "Server error", 500)
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(current.PasswordHash), []byte(currentPassword)) != nil {
			setFlash(w, "Current password is incorrect. Credentials were not changed.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}

		if newUsername == "" {
			setFlash(w, "Username cannot be empty.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		if newPassword != "" {
			if newPassword != confirmPassword {
				setFlash(w, "Passwords do not match.")
				http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
				return
			}
			if len(newPassword) < 8 {
				setFlash(w, "Password must be at least 8 characters.")
				http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
				return
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}
			_, err = app.db.Exec(`UPDATE admin_users SET username=?, password_hash=? WHERE id=?`,
				newUsername, string(hash), s.AdminUserID)
			if err != nil {
				setFlash(w, "Error updating credentials: "+err.Error())
				http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
				return
			}
		} else {
			_, err := app.db.Exec(`UPDATE admin_users SET username=? WHERE id=?`, newUsername, s.AdminUserID)
			if err != nil {
				setFlash(w, "Error updating username: "+err.Error())
				http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
				return
			}
		}
		app.logAudit(s.AdminUserID, "update_credentials", "admin_user", s.AdminUserID, "")
		setFlash(w, "Credentials updated. Please log in again if you changed your password.")

	case "email":
		notifyEmail := strings.TrimSpace(r.FormValue("notification_email"))
		senderEmail := strings.TrimSpace(r.FormValue("sender_email"))
		emailEnabled := lastFormValue(r, "email_enabled") == "1"
		expiryHours, eerr := strconv.Atoi(r.FormValue("verify_expiry_hours"))
		if eerr != nil || expiryHours < 1 {
			expiryHours = 24
		}
		// Email can only be enabled when an API key and sender address are present.
		if emailEnabled && (app.resendKey == "" || senderEmail == "") {
			setFlash(w, "To enable email, set RESEND_API_KEY and a sender email address.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		_, err := app.db.Exec(
			`UPDATE site_preferences SET notification_email=?, sender_email=?, email_enabled=?, verify_expiry_hours=? WHERE id=1`,
			notifyEmail, senderEmail, boolToInt(emailEnabled), expiryHours,
		)
		if err != nil {
			setFlash(w, "Error saving email settings: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		app.logAudit(s.AdminUserID, "update_email_settings", "site_preferences", 1, "")
		setFlash(w, "Email settings saved.")

	case "rate_limits":
		perMin, err1 := strconv.Atoi(r.FormValue("rate_per_minute"))
		perHour, err2 := strconv.Atoi(r.FormValue("rate_per_hour"))
		perDay, err3 := strconv.Atoi(r.FormValue("rate_per_day"))
		tHour, err4 := strconv.Atoi(r.FormValue("trusted_rate_per_hour"))
		tDay, err5 := strconv.Atoi(r.FormValue("trusted_rate_per_day"))
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || perMin < 1 || perHour < 1 || perDay < 1 || tHour < 1 || tDay < 1 {
			setFlash(w, "Rate limits must be positive integers.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		_, err := app.db.Exec(`UPDATE site_preferences SET rate_per_minute=?, rate_per_hour=?, rate_per_day=?, trusted_rate_per_hour=?, trusted_rate_per_day=? WHERE id=1`, perMin, perHour, perDay, tHour, tDay)
		if err != nil {
			setFlash(w, "Error saving rate limits: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		app.logAudit(s.AdminUserID, "update_rate_limits", "site_preferences", 1, "")
		setFlash(w, "Rate limits saved.")

	case "register_rate_limits":
		perHour, err1 := strconv.Atoi(r.FormValue("register_rate_per_hour"))
		perDay, err2 := strconv.Atoi(r.FormValue("register_rate_per_day"))
		if err1 != nil || err2 != nil || perHour < 1 || perDay < 1 {
			setFlash(w, "Registration rate limits must be positive integers.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		_, err := app.db.Exec(`UPDATE site_preferences SET register_rate_per_hour=?, register_rate_per_day=? WHERE id=1`, perHour, perDay)
		if err != nil {
			setFlash(w, "Error saving registration rate limits: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		app.logAudit(s.AdminUserID, "update_register_rate_limits", "site_preferences", 1, "")
		setFlash(w, "Registration rate limits saved.")

	case "comment_rate_limits":
		perHour, err1 := strconv.Atoi(r.FormValue("comment_rate_per_hour"))
		perDay, err2 := strconv.Atoi(r.FormValue("comment_rate_per_day"))
		tHour, err3 := strconv.Atoi(r.FormValue("trusted_comment_rate_per_hour"))
		tDay, err4 := strconv.Atoi(r.FormValue("trusted_comment_rate_per_day"))
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || perHour < 1 || perDay < 1 || tHour < 1 || tDay < 1 {
			setFlash(w, "Comment rate limits must be positive integers.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		_, err := app.db.Exec(`UPDATE site_preferences SET comment_rate_per_hour=?, comment_rate_per_day=?, trusted_comment_rate_per_hour=?, trusted_comment_rate_per_day=? WHERE id=1`, perHour, perDay, tHour, tDay)
		if err != nil {
			setFlash(w, "Error saving comment rate limits: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		app.logAudit(s.AdminUserID, "update_comment_rate_limits", "site_preferences", 1, "")
		setFlash(w, "Comment rate limits saved.")

	case "admin_compact":
		on := lastFormValue(r, "admin_compact") == "1"
		if _, err := app.db.Exec(`UPDATE site_preferences SET admin_compact=? WHERE id=1`, on); err != nil {
			setFlash(w, "Error saving compact setting.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		setAdminCompact(on)
		app.logAudit(s.AdminUserID, "update_admin_compact", "site_preferences", 1, fmt.Sprintf("%v", on))
		setFlash(w, "Compact mode updated.")

	case "admin_theme":
		theme := r.FormValue("admin_theme")
		validTheme := false
		for _, n := range adminThemeNames {
			if n == theme {
				validTheme = true
				break
			}
		}
		if !validTheme {
			setFlash(w, "Unknown theme.")
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		_, err := app.db.Exec(`UPDATE site_preferences SET admin_theme=? WHERE id=1`, theme)
		if err != nil {
			setFlash(w, "Error saving theme: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		setAdminTheme(theme)
		app.logAudit(s.AdminUserID, "update_admin_theme", "site_preferences", 1, theme)
		setFlash(w, "Admin theme updated.")

	case "user_approval":
		autoConfirm := lastFormValue(r, "auto_approve_on_confirm") == "1"
		autoVideo := lastFormValue(r, "auto_approve_on_video") == "1"
		if _, err := app.db.Exec(`UPDATE site_preferences SET auto_approve_on_confirm=?, auto_approve_on_video=? WHERE id=1`,
			boolToInt(autoConfirm), boolToInt(autoVideo)); err != nil {
			setFlash(w, "Error saving user approval settings: "+err.Error())
			http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
			return
		}
		app.logAudit(s.AdminUserID, "update_user_approval", "site_preferences", 1, "")
		setFlash(w, "User approval settings saved.")

	default:
		setFlash(w, "Unknown action.")
	}

	http.Redirect(w, r, app.adminPrefix+"/settings", http.StatusSeeOther)
}

func (app *App) handleAdminTrainStops(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	train, err := trainByID(app.db, id)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	stops, err := stopsByTrainID(app.db, id)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	s := sessionFromCtx(r)
	type stopsPage struct {
		Train Train
		Stops []TrainStop
	}
	app.renderAdmin(w, r, "train_stops.html", adminPage{
		Title:     "Stops — " + train.DisplayName,
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      stopsPage{Train: train, Stops: stops},
	})
}

func (app *App) handleAdminTrainStopsUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	stops, err := stopsByTrainID(app.db, id)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	for _, ts := range stops {
		sid := fmt.Sprintf("%d", ts.ID)
		arr := strings.TrimSpace(r.FormValue("arr_" + sid))
		dep := strings.TrimSpace(r.FormValue("dep_" + sid))
		wkdy := r.FormValue("wkdy_"+sid) == "1"
		wknd := r.FormValue("wknd_"+sid) == "1"

		var arrVal, depVal interface{}
		if arr != "" {
			arrVal = arr
		}
		if dep != "" {
			depVal = dep
		}
		app.db.Exec(
			`UPDATE train_stops SET scheduled_arrival=?, scheduled_departure=?, runs_weekday=?, runs_weekend=? WHERE id=?`,
			arrVal, depVal, wkdy, wknd, ts.ID,
		)
	}

	setFlash(w, "Stops updated.")
	http.Redirect(w, r, app.adminPrefix+"/trains/"+fmt.Sprintf("%d", id)+"/stops", http.StatusSeeOther)
}
