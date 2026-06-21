package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type publicPage struct {
	Title       string
	Flash       string
	Data        interface{}
	CurrentUser *User
	UserCSRF    string
	FormValues  map[string]string
}

func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	flash := getFlash(w, r)
	loggedIn, _ := app.getUserSession(r)

	// Serve cached HTML when there is no flash message and no user session
	// (the cached HTML renders the anonymous account nav).
	if flash == "" && loggedIn == nil {
		app.indexCacheMu.RLock()
		cached := app.indexCacheHTML
		app.indexCacheMu.RUnlock()
		if cached != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	corridors, err := allCorridors(app.db, true)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	data := publicPage{
		Title: "AmazingTrak — Amtrak Train Tracker",
		Flash: flash,
		Data:  corridors,
	}

	// Flash responses and logged-in views are per-request; don't cache them.
	if flash != "" || loggedIn != nil {
		app.renderPublic(w, r, "index.html", data)
		return
	}

	b, err := app.renderPublicToBuffer(r, "index.html", data)
	if err != nil {
		log.Printf("render index: %v", err)
		http.Error(w, "Template error", 500)
		return
	}
	app.indexCacheMu.Lock()
	app.indexCacheHTML = b
	app.indexCacheMu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b)
}

func (app *App) handleOverview(w http.ResponseWriter, r *http.Request) {
	type topVideo struct {
		TrainName   string
		TrainSlug   string
		Title       string
		URL         string
		Tags        string
		RarityCount int
		Contributor string
	}
	type topUser struct {
		Username    string
		RarityCount int
	}
	type latestVideo struct {
		TrainName string
		TrainSlug string
		Title     string
		URL       string
		Tags      string
		CreatedAt string
	}
	type missingCorridor struct {
		CorridorName string
		CorridorSlug string
		MissingCount int
	}
	type latestComment struct {
		Username  string
		TrainName string
		TrainSlug string
		Body      string
		CreatedAt string
	}
	type overviewData struct {
		TopVideos        []topVideo
		TopUsers         []topUser
		LatestVideos     []latestVideo
		MissingCorridors []missingCorridor
		LatestComments   []latestComment
	}

	raritySQL := `(CASE WHEN instr(COALESCE(tags,''),'long_consist')>0 THEN 1 ELSE 0 END +
	              CASE WHEN instr(COALESCE(tags,''),'doubleheader')>0 THEN 1 ELSE 0 END +
	              CASE WHEN instr(COALESCE(tags,''),'sandwich_set')>0 THEN 1 ELSE 0 END +
	              CASE WHEN instr(COALESCE(tags,''),'reverse_set')>0 THEN 1 ELSE 0 END)`

	var data overviewData

	// Top 5 videos by rarity count
	rows, err := app.db.Query(`SELECT t.display_name, t.slug, COALESCE(m.title,''), m.url, COALESCE(m.tags,''),
		` + raritySQL + ` AS rc,
		COALESCE(u.username, CASE m.added_by WHEN 'admin' THEN 'Admin' ELSE 'Anonymous' END)
		FROM media m
		JOIN trains t ON t.id=m.train_id
		LEFT JOIN users u ON u.id=m.user_id
		WHERE m.media_type='video' AND ` + raritySQL + `>0
		ORDER BY rc DESC, m.created_at DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var v topVideo
			rows.Scan(&v.TrainName, &v.TrainSlug, &v.Title, &v.URL, &v.Tags, &v.RarityCount, &v.Contributor)
			data.TopVideos = append(data.TopVideos, v)
		}
		rows.Close()
	}

	// Top 5 users by total approved submissions
	rows, err = app.db.Query(`SELECT u.username, COUNT(*) AS cnt
		FROM suggestions s
		JOIN users u ON u.id=s.user_id
		WHERE s.status='approved'
		GROUP BY u.id
		ORDER BY cnt DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var tu topUser
			rows.Scan(&tu.Username, &tu.RarityCount)
			data.TopUsers = append(data.TopUsers, tu)
		}
		rows.Close()
	}

	// Latest 5 videos
	rows, err = app.db.Query(`SELECT t.display_name, t.slug, COALESCE(m.title,''), m.url, COALESCE(m.tags,''), m.created_at
		FROM media m
		JOIN trains t ON t.id=m.train_id
		WHERE m.media_type='video' AND m.is_published=1
		ORDER BY m.created_at DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var lv latestVideo
			rows.Scan(&lv.TrainName, &lv.TrainSlug, &lv.Title, &lv.URL, &lv.Tags, &lv.CreatedAt)
			data.LatestVideos = append(data.LatestVideos, lv)
		}
		rows.Close()
	}

	// Top corridors with most active trains that have no approved video at all
	rows, err = app.db.Query(`SELECT c.name, c.slug,
		COUNT(CASE WHEN (SELECT COUNT(*) FROM media WHERE train_id=t.id AND media_type='video' AND is_published=1)=0 THEN 1 END) AS missing
		FROM corridors c
		JOIN trains t ON t.corridor_id=c.id AND t.is_active=1
		WHERE c.is_active=1
		GROUP BY c.id
		HAVING missing>0
		ORDER BY missing DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var mc missingCorridor
			rows.Scan(&mc.CorridorName, &mc.CorridorSlug, &mc.MissingCount)
			data.MissingCorridors = append(data.MissingCorridors, mc)
		}
		rows.Close()
	}

	// Latest 5 approved comments
	rows, err = app.db.Query(`SELECT u.username, t.display_name, t.slug, c.body, c.created_at
		FROM comments c
		JOIN users u ON u.id=c.user_id
		JOIN trains t ON t.id=c.train_id
		WHERE c.status='approved'
		ORDER BY c.created_at DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var lc latestComment
			rows.Scan(&lc.Username, &lc.TrainName, &lc.TrainSlug, &lc.Body, &lc.CreatedAt)
			data.LatestComments = append(data.LatestComments, lc)
		}
		rows.Close()
	}

	app.renderPublic(w, r, "overview.html", publicPage{
		Title: "Overview — " + getSiteName(),
		Data:  data,
	})
}

func (app *App) handleTrainsList(w http.ResponseWriter, r *http.Request) {
	all, err := allTrains(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	var active []Train
	for _, t := range all {
		if t.IsActive {
			active = append(active, t)
		}
	}
	app.renderPublic(w, r, "trains_list.html", publicPage{
		Title: "Trains — AmazingTrak",
		Data:  active,
	})
}

func (app *App) handleMap(w http.ResponseWriter, r *http.Request) {
	corridors, _ := allCorridors(app.db, true)
	app.renderPublic(w, r, "map.html", publicPage{
		Title: "Amtrak Route Map — AmazingTrak",
		Data:  corridors,
	})
}

func (app *App) handleCorridors(w http.ResponseWriter, r *http.Request) {
	corridors, err := allCorridors(app.db, true)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	app.renderPublic(w, r, "corridors.html", publicPage{
		Title: "Routes — AmazingTrak",
		Data:  corridors,
	})
}

func (app *App) handleCorridor(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	corridor, err := corridorBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if !corridor.IsActive {
		http.NotFound(w, r)
		return
	}

	// The conductor of this corridor sees inactive trains too, so they can
	// reactivate them; everyone else sees only active trains.
	viewer, _, canManage := app.conductorGuard(r, corridor.ID)

	trains, err := trainsByCorridorID(app.db, corridor.ID, !canManage)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	media, err := mediaByCorridorID(app.db, corridor.ID, true)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	stops, err := stopsByCorridorID(app.db, corridor.ID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	var bestVideo *Media
	for i, m := range media {
		if m.MediaType != "video" {
			continue
		}
		if bestVideo == nil || (m.IsBest && !bestVideo.IsBest) {
			bestVideo = &media[i]
		}
	}

	// Conductor affordances: CanManage when the viewer conducts this corridor;
	// CanRequest when a logged-in non-spammer (who isn't already the conductor)
	// could request the role; HasRequested when they already have one pending;
	// NeedsVerify when email is on and they must verify before requesting.
	emailOn := app.emailEnabled()
	canRequest := false
	hasRequested := false
	needsVerify := false
	if viewer != nil && !viewer.IsSpammer && !canManage && !corridor.ConductorUserID.Valid {
		if pending, _ := pendingConductorRequest(app.db, corridor.ID, viewer.ID); pending {
			hasRequested = true
		} else if emailOn && !viewer.EmailConfirmed {
			needsVerify = true
		} else {
			canRequest = true
		}
	}

	// Corridor comments: approved ones for everyone, plus the viewer's own pending.
	corridorComments, _ := commentsByCorridorID(app.db, corridor.ID, "approved")
	var ownPendingComments []Comment
	if viewer != nil {
		setTimingCookie(w)
		rows, _ := app.db.Query(
			commentSelectBase+` WHERE c.corridor_id=? AND c.user_id=? AND c.status='pending' ORDER BY c.created_at DESC`,
			corridor.ID, viewer.ID,
		)
		if rows != nil {
			ownPendingComments, _ = scanComments(rows)
			rows.Close()
		}
	}

	type corridorDetail struct {
		Corridor           Corridor
		Trains             []Train
		Media              []Media
		BestVideo          *Media
		Stops              []Stop
		CanManage          bool
		CanRequest         bool
		HasRequested       bool
		NeedsVerify        bool
		Comments           []Comment
		OwnPendingComments []Comment
	}
	app.renderPublic(w, r, "corridor.html", publicPage{
		Title: corridor.Name + " — AmazingTrak",
		Flash: getFlash(w, r),
		Data: corridorDetail{
			Corridor: corridor, Trains: trains, Media: media, BestVideo: bestVideo, Stops: stops,
			CanManage: canManage, CanRequest: canRequest, HasRequested: hasRequested, NeedsVerify: needsVerify,
			Comments: corridorComments, OwnPendingComments: ownPendingComments,
		},
	})
}

func (app *App) handleTrain(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	train, err := trainBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if !train.IsActive {
		http.NotFound(w, r)
		return
	}

	allMedia, err := mediaByTrainID(app.db, train.ID, true)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	stops, err := stopsByTrainID(app.db, train.ID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	comments, err := commentsByTrainID(app.db, train.ID, "approved")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	// The comment form is shown to logged-in users; prime the anti-bot timing
	// cookie so a genuine fill always clears the minimum-duration check.
	// Also load the user's own pending comments so they can see them greyed out.
	var ownPendingComments []Comment
	if u, _ := app.getUserSession(r); u != nil {
		setTimingCookie(w)
		rows, _ := app.db.Query(
			commentSelectBase+` WHERE c.train_id=? AND c.user_id=? AND c.status='pending' ORDER BY c.created_at DESC`,
			train.ID, u.ID,
		)
		if rows != nil {
			ownPendingComments, _ = scanComments(rows)
			rows.Close()
		}
	}

	var heroMedia, mapMedia *Media
	var images, videos, websites []Media
	for i, m := range allMedia {
		if train.HeroMediaID.Valid && m.ID == train.HeroMediaID.Int64 {
			heroMedia = &allMedia[i]
		}
		if train.MapMediaID.Valid && m.ID == train.MapMediaID.Int64 {
			mapMedia = &allMedia[i]
		}
		switch m.MediaType {
		case "image":
			isHero := train.HeroMediaID.Valid && m.ID == train.HeroMediaID.Int64
			isMap := train.MapMediaID.Valid && m.ID == train.MapMediaID.Int64
			isThumb := train.ThumbnailMediaID.Valid && m.ID == train.ThumbnailMediaID.Int64
			if !isHero && !isMap && !isThumb {
				images = append(images, m)
			}
		case "video":
			videos = append(videos, m)
		case "website":
			websites = append(websites, m)
		}
	}

	sort.SliceStable(videos, func(i, j int) bool {
		if videos[i].IsBest != videos[j].IsBest {
			return videos[i].IsBest
		}
		if rarityCount(videos[i].Tags) != rarityCount(videos[j].Tags) {
			return rarityCount(videos[i].Tags) > rarityCount(videos[j].Tags)
		}
		return videos[i].CreatedAt > videos[j].CreatedAt
	})

	var bestVideo *Media
	if len(videos) > 0 {
		bestVideo = &videos[0]
	}

	// Summary of rarity tags present across this train's videos, for the
	// anchor-link line above the best video (only tags that actually occur).
	rarityCounts := map[string]int{}
	for _, v := range videos {
		for _, t := range strings.Split(v.Tags, ",") {
			t = strings.TrimSpace(t)
			if rarityTags[t] {
				rarityCounts[t]++
			}
		}
	}
	type raritySummary struct {
		Tag   string
		Emoji string
		Label string
		Count int
	}
	var rarities []raritySummary
	for _, tag := range rarityOrder {
		if n := rarityCounts[tag]; n > 0 {
			rarities = append(rarities, raritySummary{Tag: tag, Emoji: rarityEmojis[tag], Label: rarityLabels[tag], Count: n})
		}
	}

	// CanManage drives the conductor edit affordances; true when the logged-in
	// user is the conductor of this train's corridor.
	_, _, canManage := app.conductorGuard(r, train.CorridorID)

	type trainDetail struct {
		Train              Train
		HeroMedia          *Media
		MapMedia           *Media
		Images             []Media
		Videos             []Media
		BestVideo          *Media
		Rarities           []raritySummary
		Websites           []Media
		Stops              []TrainStop
		Comments           []Comment
		OwnPendingComments []Comment
		CanManage          bool
	}
	app.renderPublic(w, r, "train.html", publicPage{
		Title: train.DisplayName + " — AmazingTrak",
		Flash: getFlash(w, r),
		Data: trainDetail{
			Train:              train,
			HeroMedia:          heroMedia,
			MapMedia:           mapMedia,
			Images:             images,
			Videos:             videos,
			BestVideo:          bestVideo,
			Rarities:           rarities,
			Websites:           websites,
			Stops:              stops,
			Comments:           comments,
			OwnPendingComments: ownPendingComments,
			CanManage:          canManage,
		},
	})
}

func (app *App) handleSuggestForm(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	train, err := trainBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if !train.IsActive {
		http.NotFound(w, r)
		return
	}

	setTimingCookie(w)

	app.renderPublic(w, r, "suggest.html", publicPage{
		Title: "Suggest Media — " + train.DisplayName,
		Flash: getFlash(w, r),
		Data:  train,
	})
}

func (app *App) handleSuggestSubmit(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	train, err := trainBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil || !train.IsActive {
		http.Error(w, "Not available", 400)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	// Honeypot fields — silently succeed without saving
	if r.FormValue("website") != "" || r.FormValue("a") != "" || r.FormValue("b") != "ok" {
		setFlash(w, "Thanks! Your suggestion has been submitted for review.")
		http.Redirect(w, r, "/trains/"+slug, http.StatusSeeOther)
		return
	}

	// Timing check
	if !checkTiming(r) {
		setFlash(w, "Please take a moment to fill out the form.")
		http.Redirect(w, r, "/trains/"+slug+"/suggest", http.StatusSeeOther)
		return
	}

	ipHash := hashIP(r)

	// Rate limits by tier: conductors are unlimited; approved/auto-approved users
	// get the trusted tier; everyone else (anonymous + unapproved) gets the
	// standard tier. A newly-registered user's very first submission also bypasses
	// the per-minute throttle so they don't have to wait right after signing up.
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
		// First-submission-free: a logged-in user with no prior suggestions skips
		// the per-minute wait once.
		if currentUser != nil {
			var priorSubs int
			app.db.QueryRow(`SELECT COUNT(*) FROM suggestions WHERE user_id=?`, currentUser.ID).Scan(&priorSubs)
			if priorSubs == 0 {
				perMin = 1 << 30 // effectively unlimited for this one submission
			}
		}
		if blocked, reason := app.checkSuggestRateLimit(perMin, perHour, perDay); blocked {
			setFlash(w, reason)
			http.Redirect(w, r, "/trains/"+slug+"/suggest", http.StatusSeeOther)
			return
		}
	}

	rawURL := r.FormValue("url")
	title := sanitizeTitle(r.FormValue("title"))
	caption := strings.TrimSpace(r.FormValue("comment"))
	if len(caption) > 500 {
		caption = caption[:500]
	}
	tags := parseVideoTags(r.Form["tags"])

	renderFormError := func(msg string) {
		setTimingCookie(w)
		fv := map[string]string{
			"url":     rawURL,
			"title":   r.FormValue("title"),
			"comment": r.FormValue("comment"),
		}
		app.renderPublic(w, r, "suggest.html", publicPage{
			Title:      "Suggest Media — " + train.DisplayName,
			Flash:      msg,
			Data:       train,
			FormValues: fv,
		})
	}

	domain, mediaType, normURL, ok := classifyPublicURL(rawURL)
	if !ok {
		renderFormError("That URL is not allowed. Submit links from YouTube, Vimeo, Flickr, Imgur, RailPictures.net, or RRPictureArchives.net.")
		return
	}

	if app.checkDuplicateSuggestion(train.ID, normURL) {
		renderFormError("This link has already been submitted.")
		return
	}

	// Fetch the YouTube title when we have one, both to populate a missing
	// title and to drive auto-approval below.
	var ytTitle string
	if domain == "youtube.com" {
		if t, ok := fetchYouTubeTitle(normURL); ok {
			ytTitle = sanitizeTitle(t)
			if title == "" {
				title = ytTitle
			}
		}
	}

	// Auto-approve a video when the fetched YouTube title contains the train's
	// (3-digit) number as a standalone token — a strong signal it's genuine.
	// Videos flagged with a rarity always need manual review, since the
	// rarity claim itself needs human verification before publishing.
	autoApprove := mediaType == "video" && ytTitle != "" && titleMatchesTrainNumber(ytTitle, train.TrainNumber) && !hasRarity(tags)

	ua := r.UserAgent()
	if len(ua) > 200 {
		ua = ua[:200]
	}

	// Attribute the submission to the logged-in user, if any (anonymous = NULL).
	var userID sql.NullInt64
	if u, _ := app.getUserSession(r); u != nil {
		userID = sql.NullInt64{Int64: u.ID, Valid: true}
	}

	status := "pending"
	if autoApprove {
		status = "approved"
	}
	res, err := app.db.Exec(
		`INSERT INTO suggestions (train_id, url, title, caption, tags, media_type, source_domain, submitter_ip_hash, submitter_user_agent, status, auto_approved, reviewed_at, user_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CASE WHEN ?=1 THEN CURRENT_TIMESTAMP ELSE NULL END, ?)`,
		train.ID, normURL, title, caption, tags, mediaType, domain, ipHash, ua, status, boolToInt(autoApprove), boolToInt(autoApprove), userID,
	)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	sugID, _ := res.LastInsertId()

	if autoApprove {
		mdomain, murl, ok := validateAdminURL(normURL)
		if !ok {
			murl, mdomain = normURL, domain
		}
		app.db.Exec(
			`INSERT INTO media (train_id, media_type, source_type, url, title, caption, tags, source_domain, location_source, is_published, added_by, user_id)
			 VALUES (?, ?, 'url', ?, ?, ?, ?, ?, 'unknown', 1, 'approved_suggestion', ?)`,
			train.ID, mediaType, murl, title, caption, tags, mdomain, userID,
		)
		go app.maybeAutoThumbnail(train.ID, murl, mdomain)
	}

	app.recordRateLimit(ipHash)

	// Notify the admin when this submission moves the pending total across a
	// threshold (only pending — auto-approved ones never need review).
	if !autoApprove {
		app.maybeNotifyPending()
	}

	// Best-effort email notification
	if prefs, err := getSitePrefs(app.db); err == nil && prefs.NotificationEmail != "" {
		sug := Suggestion{ID: sugID, URL: normURL, Title: title, MediaType: mediaType, SourceDomain: domain}
		go app.sendSuggestionEmail(prefs.NotificationEmail, train, sug, app.baseURL)
	}

	setFlash(w, "Thanks! Your suggestion has been submitted for review.")
	http.Redirect(w, r, "/trains/"+slug, http.StatusSeeOther)
}

// handleCommentSubmit accepts a comment on a train from a logged-in registered
// user (route gated by requireUser).
func (app *App) handleCommentSubmit(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	train, err := trainBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil || !train.IsActive {
		http.Error(w, "Not available", 400)
		return
	}
	app.submitComment(w, r, train.ID, 0, "/trains/"+slug+"#comments")
}

// handleCorridorCommentSubmit accepts a comment on a corridor (route gated by
// requireUser).
func (app *App) handleCorridorCommentSubmit(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	corridor, err := corridorBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil || !corridor.IsActive {
		http.Error(w, "Not available", 400)
		return
	}
	app.submitComment(w, r, 0, corridor.ID, "/routes/"+slug+"#comments")
}

// submitComment is the shared comment pipeline for trains and corridors: exactly
// one of trainID/corridorID is non-zero. Comments enter as 'pending' and are
// moderated, with honeypot + timing bot defenses, tiered per-user rate limits
// (conductors unlimited, approved users get the trusted tier), duplicate
// detection, and auto-approval of an approved user's first few comments.
func (app *App) submitComment(w http.ResponseWriter, r *http.Request, trainID, corridorID int64, commentURL string) {
	currentUser, csrf := app.getUserSession(r)
	if currentUser == nil {
		setFlash(w, "Please log in to comment.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	if r.FormValue("csrf_token") != csrf {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}

	// Silently drop spammer comments, and silently succeed on honeypot hits.
	if currentUser.IsSpammer ||
		r.FormValue("website") != "" || r.FormValue("a") != "" || r.FormValue("b") != "ok" {
		setFlash(w, "Thanks! Your comment has been submitted for review.")
		http.Redirect(w, r, commentURL, http.StatusSeeOther)
		return
	}
	if !checkTiming(r) {
		setFlash(w, "Please take a moment to write your comment.")
		http.Redirect(w, r, commentURL, http.StatusSeeOther)
		return
	}

	// Rate limits: conductors are unlimited; approved users get the trusted tier;
	// everyone else gets the standard tier. Messages tell the user how long to wait.
	if !isAnyConductor(app.db, currentUser.ID) {
		prefs, _ := getSitePrefs(app.db)
		perHour, perDay := prefs.CommentRatePerHour, prefs.CommentRatePerDay
		if currentUser.IsApproved() {
			perHour, perDay = prefs.TrustedCommentRatePerHour, prefs.TrustedCommentRatePerDay
		}
		if perHour <= 0 {
			perHour = 10
		}
		if perDay <= 0 {
			perDay = 50
		}
		if blocked, reason := app.checkUserCommentRateLimit(currentUser.ID, perHour, perDay); blocked {
			setFlash(w, reason)
			http.Redirect(w, r, commentURL, http.StatusSeeOther)
			return
		}
	}

	body := sanitizeComment(r.FormValue("body"))
	if body == "" {
		setFlash(w, "Your comment was empty.")
		http.Redirect(w, r, commentURL, http.StatusSeeOther)
		return
	}

	dup := false
	if corridorID != 0 {
		dup = app.checkDuplicateCorridorComment(corridorID, currentUser.ID, body)
	} else {
		dup = app.checkDuplicateComment(trainID, currentUser.ID, body)
	}
	if dup {
		setFlash(w, "You've already posted that comment.")
		http.Redirect(w, r, commentURL, http.StatusSeeOther)
		return
	}

	// Auto-approve the first 3 comments from approved users.
	status := "pending"
	if currentUser.IsApproved() {
		var approvedCount int
		app.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE user_id=? AND status='approved'`, currentUser.ID).Scan(&approvedCount)
		if approvedCount < 3 {
			status = "approved"
		}
	}

	var trainVal, corridorVal interface{}
	if trainID != 0 {
		trainVal = trainID
	}
	if corridorID != 0 {
		corridorVal = corridorID
	}
	_, err := app.db.Exec(
		`INSERT INTO comments (train_id, corridor_id, user_id, body, status, submitter_ip_hash) VALUES (?, ?, ?, ?, ?, ?)`,
		trainVal, corridorVal, currentUser.ID, body, status, hashIP(r),
	)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if status == "pending" {
		app.maybeNotifyPending()
		setFlash(w, "Thanks! Your comment has been submitted and will appear once an admin approves it.")
	} else {
		setFlash(w, "Your comment has been posted.")
	}
	http.Redirect(w, r, commentURL, http.StatusSeeOther)
}

func (app *App) handleStation(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	stops, err := stopsBySlug(app.db, slug)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	if len(stops) == 0 {
		http.NotFound(w, r)
		return
	}

	// Use the first stop's name/code as the canonical station identity
	stationName := stops[0].Name
	stationCode := stops[0].StationCode

	type corridorEntry struct {
		Corridor Corridor
		Trains   []StationTrain
	}
	var entries []corridorEntry
	corridorSeen := map[int64]int{} // corridorID → index in entries

	for _, stop := range stops {
		trains, err := trainsByStopID(app.db, stop.ID)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		corridor, err := corridorByID(app.db, stop.CorridorID)
		if err != nil {
			continue
		}
		if idx, ok := corridorSeen[stop.CorridorID]; ok {
			entries[idx].Trains = append(entries[idx].Trains, trains...)
		} else {
			corridorSeen[stop.CorridorID] = len(entries)
			entries = append(entries, corridorEntry{Corridor: corridor, Trains: trains})
		}
	}

	type stationData struct {
		StationName string
		StationCode string
		Slug        string
		Entries     []corridorEntry
	}
	app.renderPublic(w, r, "station.html", publicPage{
		Title: stationName + " — AmazingTrak",
		Data: stationData{
			StationName: stationName,
			StationCode: stationCode,
			Slug:        slug,
			Entries:     entries,
		},
	})
}

var rarityTags = map[string]bool{
	"long_consist": true, "doubleheader": true, "sandwich_set": true, "reverse_set": true,
	"scenic": true, "environment": true, "historic": true, "special_event": true,
}

func rarityCount(tags string) int {
	n := 0
	for _, t := range strings.Split(tags, ",") {
		if rarityTags[strings.TrimSpace(t)] {
			n++
		}
	}
	return n
}

// hasRarity reports whether a comma-separated tag string contains any rarity tag.
func hasRarity(tags string) bool {
	return rarityCount(tags) > 0
}

// ----- Amtrak routes proxy (server-side cache to avoid CORS/size issues) -----

var (
	amtrakRoutesMu      sync.Mutex
	amtrakRoutesJSON    []byte
	amtrakRoutesFetched time.Time
)

func fetchAmtrakRoutes() ([]byte, error) {
	type arcFC struct {
		Features   []json.RawMessage `json:"features"`
		Properties struct {
			ExceededTransferLimit bool `json:"exceededTransferLimit"`
		} `json:"properties"`
	}

	client := &http.Client{Timeout: 45 * time.Second}
	var allFeatures []json.RawMessage
	const pageSize = 1000

	for offset := 0; ; offset += pageSize {
		u := fmt.Sprintf(
			"https://services.arcgis.com/xOi1kZaI0eWDREZv/arcgis/rest/services/NTAD_Amtrak_Routes/FeatureServer/0/query"+
				"?where=1%%3D1&outFields=name%%2CROUTE_URL&f=geojson&outSR=4326&resultRecordCount=%d&resultOffset=%d&geometryPrecision=4&maxAllowableOffset=0.01",
			pageSize, offset)
		resp, err := client.Get(u)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		var fc arcFC
		if err := json.Unmarshal(body, &fc); err != nil {
			return nil, err
		}
		allFeatures = append(allFeatures, fc.Features...)
		if !fc.Properties.ExceededTransferLimit || len(fc.Features) == 0 {
			break
		}
	}

	result := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": allFeatures,
	}
	return json.Marshal(result)
}

func (app *App) handleAmtrakRoutes(w http.ResponseWriter, r *http.Request) {
	amtrakRoutesMu.Lock()
	if amtrakRoutesJSON != nil && time.Since(amtrakRoutesFetched) < 12*time.Hour {
		data := amtrakRoutesJSON
		amtrakRoutesMu.Unlock()
		w.Header().Set("Content-Type", "application/geo+json")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write(data)
		return
	}
	amtrakRoutesMu.Unlock()

	data, err := fetchAmtrakRoutes()
	if err != nil {
		http.Error(w, "Failed to fetch routes", 502)
		return
	}

	amtrakRoutesMu.Lock()
	amtrakRoutesJSON = data
	amtrakRoutesFetched = time.Now()
	amtrakRoutesMu.Unlock()

	w.Header().Set("Content-Type", "application/geo+json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(data)
}
