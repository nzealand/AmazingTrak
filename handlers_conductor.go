package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Conductor handlers let a corridor's assigned conductor manage its trains
// directly through the public site — they never see the secret admin area.
// Access is gated by conductorGuard (auth.go); CSRF uses the user-session token,
// mirroring the comment flow in handlers_public.go.

// conductorCSRFOK validates the form CSRF token against the user session token.
func conductorCSRFOK(r *http.Request, csrf string) bool {
	return csrf != "" && r.FormValue("csrf_token") == csrf
}

// lastFormValue returns the last submitted value for a form key, or "". A
// checkbox preceded by a hidden field submits both values; the checkbox (last)
// wins.
func lastFormValue(r *http.Request, key string) string {
	vals := r.Form[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[len(vals)-1]
}

// handleConductorRequest records a registered user's request to maintain a
// corridor. The route is wrapped in requireUser, so anonymous visitors are sent
// to /login before reaching here.
func (app *App) handleConductorRequest(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	corridor, err := corridorBySlug(app.db, slug)
	if err == sql.ErrNoRows || (err == nil && !corridor.IsActive) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	user, csrf := app.getUserSession(r)
	if user == nil {
		setFlash(w, "Please log in to request the Conductor role.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !conductorCSRFOK(r, csrf) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}

	corridorURL := "/routes/" + slug

	// Spammers are silently no-op'd.
	if user.IsSpammer {
		setFlash(w, "Your request has been submitted.")
		http.Redirect(w, r, corridorURL, http.StatusSeeOther)
		return
	}
	// When email is enabled, only verified addresses may request the role. When
	// email is off, verification is impossible so this gate is skipped.
	if app.emailEnabled() && !user.EmailConfirmed {
		setFlash(w, "Please verify your email address before requesting the Conductor role.")
		http.Redirect(w, r, corridorURL, http.StatusSeeOther)
		return
	}
	if corridor.ConductorUserID.Valid {
		setFlash(w, "This route already has a Conductor.")
		http.Redirect(w, r, corridorURL, http.StatusSeeOther)
		return
	}
	if pending, _ := pendingConductorRequest(app.db, corridor.ID, user.ID); pending {
		setFlash(w, "You already have a pending request for this route.")
		http.Redirect(w, r, corridorURL, http.StatusSeeOther)
		return
	}

	msg := strings.TrimSpace(r.FormValue("message"))
	if len(msg) > 500 {
		msg = msg[:500]
	}
	if err := createConductorRequest(app.db, corridor.ID, user.ID, msg); err != nil {
		setFlash(w, "Could not submit your request. Please try again.")
		http.Redirect(w, r, corridorURL, http.StatusSeeOther)
		return
	}

	if prefs, err := getSitePrefs(app.db); err == nil && prefs.NotificationEmail != "" {
		go app.sendConductorRequestEmail(prefs.NotificationEmail, corridor, *user, app.baseURL)
	}
	app.maybeNotifyPending()

	setFlash(w, "Thanks! Your request to maintain this route has been submitted for review.")
	http.Redirect(w, r, corridorURL, http.StatusSeeOther)
}

// conductorTrain loads the train for a slug and verifies the logged-in user
// conducts its corridor. On failure it writes the appropriate response (redirect
// or not-found) and returns ok=false.
func (app *App) conductorTrain(w http.ResponseWriter, r *http.Request) (Train, *User, string, bool) {
	slug := r.PathValue("slug")
	train, err := trainBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return Train{}, nil, "", false
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return Train{}, nil, "", false
	}
	user, csrf, ok := app.conductorGuard(r, train.CorridorID)
	if user == nil {
		setFlash(w, "Please log in to manage this train.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return Train{}, nil, "", false
	}
	if !ok {
		setFlash(w, "You don't maintain this route.")
		http.Redirect(w, r, "/trains/"+slug, http.StatusSeeOther)
		return Train{}, nil, "", false
	}
	return train, user, csrf, true
}

// handleConductorTrainEditForm renders the conductor's train edit page.
func (app *App) handleConductorTrainEditForm(w http.ResponseWriter, r *http.Request) {
	train, user, csrf, ok := app.conductorTrain(w, r)
	if !ok {
		return
	}
	corridors, err := corridorsConductedBy(app.db, user.ID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type editData struct {
		Train     Train
		Corridors []Corridor
		CSRFToken string
	}
	app.renderPublic(w, r, "train_edit.html", publicPage{
		Title: "Edit " + train.DisplayName,
		Flash: getFlash(w, r),
		Data:  editData{Train: train, Corridors: corridors, CSRFToken: csrf},
	})
}

// handleConductorTrainEdit saves edits to a train. The conductor may only move a
// train to a corridor they also conduct.
func (app *App) handleConductorTrainEdit(w http.ResponseWriter, r *http.Request) {
	train, user, csrf, ok := app.conductorTrain(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	if !conductorCSRFOK(r, csrf) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}

	displayName := strings.TrimSpace(r.FormValue("display_name"))
	newSlug := strings.TrimSpace(r.FormValue("slug"))
	trainNumber := strings.TrimSpace(r.FormValue("train_number"))
	direction := strings.TrimSpace(r.FormValue("direction"))
	notes := strings.TrimSpace(r.FormValue("notes"))
	corridorID, _ := strconv.ParseInt(r.FormValue("corridor_id"), 10, 64)
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))

	editURL := "/trains/" + train.Slug + "/edit"

	if displayName == "" || newSlug == "" || trainNumber == "" {
		setFlash(w, "Display name, slug, and train number are required.")
		http.Redirect(w, r, editURL, http.StatusSeeOther)
		return
	}

	// A conductor may only reassign the train to a corridor they also conduct.
	if corridorID != train.CorridorID {
		conducts, err := isConductorOf(app.db, user.ID, corridorID)
		if err != nil || !conducts {
			setFlash(w, "You can only move this train to a route you maintain.")
			http.Redirect(w, r, editURL, http.StatusSeeOther)
			return
		}
	}

	_, err := app.db.Exec(
		`UPDATE trains SET display_name=?, slug=?, train_number=?, direction=?, notes=?, corridor_id=?, sort_order=? WHERE id=?`,
		displayName, newSlug, trainNumber, direction, notes, corridorID, sortOrder, train.ID,
	)
	if err != nil {
		setFlash(w, "Error updating train: "+err.Error())
		http.Redirect(w, r, editURL, http.StatusSeeOther)
		return
	}
	setFlash(w, "Train updated.")
	http.Redirect(w, r, "/trains/"+newSlug+"/edit", http.StatusSeeOther)
}

// handleConductorTrainToggle activates/deactivates a train.
func (app *App) handleConductorTrainToggle(w http.ResponseWriter, r *http.Request) {
	train, _, csrf, ok := app.conductorTrain(w, r)
	if !ok {
		return
	}
	if !conductorCSRFOK(r, csrf) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	app.db.Exec(`UPDATE trains SET is_active = CASE WHEN is_active=1 THEN 0 ELSE 1 END WHERE id=?`, train.ID)
	app.invalidateIndexCache()
	setFlash(w, "Train status updated.")
	// Honour an optional return_to (a safe site-relative path) so the button can
	// send the conductor back where they clicked it. The train page 404s on an
	// inactive train, so callers there point return_to at the corridor instead.
	dest := "/trains/" + train.Slug + "/edit"
	if rt := r.FormValue("return_to"); strings.HasPrefix(rt, "/") && !strings.HasPrefix(rt, "//") {
		dest = rt
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

// conductorCorridor loads a corridor by slug and verifies the user conducts it.
func (app *App) conductorCorridor(w http.ResponseWriter, r *http.Request) (Corridor, *User, string, bool) {
	slug := r.PathValue("slug")
	corridor, err := corridorBySlug(app.db, slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return Corridor{}, nil, "", false
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return Corridor{}, nil, "", false
	}
	user, csrf, ok := app.conductorGuard(r, corridor.ID)
	if user == nil {
		setFlash(w, "Please log in to manage this route.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return Corridor{}, nil, "", false
	}
	if !ok {
		setFlash(w, "You don't maintain this route.")
		http.Redirect(w, r, "/routes/"+slug, http.StatusSeeOther)
		return Corridor{}, nil, "", false
	}
	return corridor, user, csrf, true
}

// handleConductorTrainNewForm renders the new-train form for a corridor.
func (app *App) handleConductorTrainNewForm(w http.ResponseWriter, r *http.Request) {
	corridor, _, csrf, ok := app.conductorCorridor(w, r)
	if !ok {
		return
	}
	type newData struct {
		Corridor  Corridor
		CSRFToken string
	}
	app.renderPublic(w, r, "train_new.html", publicPage{
		Title: "New Train — " + corridor.Name,
		Flash: getFlash(w, r),
		Data:  newData{Corridor: corridor, CSRFToken: csrf},
	})
}

// handleConductorTrainCreate creates a train in the conductor's corridor.
func (app *App) handleConductorTrainCreate(w http.ResponseWriter, r *http.Request) {
	corridor, _, csrf, ok := app.conductorCorridor(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	if !conductorCSRFOK(r, csrf) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	newURL := "/routes/" + corridor.Slug + "/trains/new"
	trainNumber := strings.TrimSpace(r.FormValue("train_number"))
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	if trainNumber == "" || displayName == "" {
		setFlash(w, "Train number and display name are required.")
		http.Redirect(w, r, newURL, http.StatusSeeOther)
		return
	}
	slug := slugify(displayName)
	res, err := app.db.Exec(
		`INSERT INTO trains (corridor_id, train_number, display_name, slug) VALUES (?, ?, ?, ?)`,
		corridor.ID, trainNumber, displayName, slug,
	)
	if err != nil {
		setFlash(w, "Error creating train: "+err.Error())
		http.Redirect(w, r, newURL, http.StatusSeeOther)
		return
	}
	id, _ := res.LastInsertId()
	newTrain, err := trainByID(app.db, id)
	if err != nil {
		setFlash(w, "Train created.")
		http.Redirect(w, r, "/routes/"+corridor.Slug, http.StatusSeeOther)
		return
	}
	setFlash(w, "Train created — you can now edit its details.")
	http.Redirect(w, r, "/trains/"+newTrain.Slug+"/edit", http.StatusSeeOther)
}

// scheduleStopsData backs the schedule editor template.
type scheduleStopsData struct {
	Train     Train
	Stops     []ScheduleStop
	CSRFToken string
	Error     string
}

// renderStopsForm renders the schedule editor, optionally with an error and a
// caller-supplied set of rows (used to preserve input on a validation failure).
func (app *App) renderStopsForm(w http.ResponseWriter, r *http.Request, train Train, csrf string, stops []ScheduleStop, errMsg string) {
	app.renderPublic(w, r, "train_edit_stops.html", publicPage{
		Title: "Schedule — " + train.DisplayName,
		Flash: getFlash(w, r),
		Data:  scheduleStopsData{Train: train, Stops: stops, CSRFToken: csrf, Error: errMsg},
	})
}

// handleConductorStopsForm renders the route & schedule editor. It lists every
// station in the train's corridor; the conductor fills in times for the ones the
// train actually stops at.
func (app *App) handleConductorStopsForm(w http.ResponseWriter, r *http.Request) {
	train, _, csrf, ok := app.conductorTrain(w, r)
	if !ok {
		return
	}
	stops, err := corridorStopsWithTrainTimes(app.db, train.ID, train.CorridorID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	app.renderStopsForm(w, r, train, csrf, stops, "")
}

// handleConductorStopsUpdate saves the schedule. Corridor stations define the
// possible stops; entering an arrival or departure time marks a station as an
// actual stop (upserting a train_stops row); leaving both blank removes the stop.
// At least one arrival and one departure time are required overall.
func (app *App) handleConductorStopsUpdate(w http.ResponseWriter, r *http.Request) {
	train, _, csrf, ok := app.conductorTrain(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	if !conductorCSRFOK(r, csrf) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	stops, err := corridorStopsWithTrainTimes(app.db, train.ID, train.CorridorID)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	// Read submitted values (keyed by stop id) into the rows, so we can both
	// validate and, on failure, re-render without losing the conductor's input.
	anyArr, anyDep := false, false
	for i := range stops {
		sid := fmt.Sprintf("%d", stops[i].StopID)
		stops[i].ScheduledArrival = strings.TrimSpace(r.FormValue("arr_" + sid))
		stops[i].ScheduledDeparture = strings.TrimSpace(r.FormValue("dep_" + sid))
		// Each checkbox is preceded by a hidden "0", so a checked box submits
		// ["0","1"] — read the last value to honor it.
		stops[i].RunsWeekday = lastFormValue(r, "wkdy_"+sid) == "1"
		stops[i].RunsWeekend = lastFormValue(r, "wknd_"+sid) == "1"
		stops[i].Stops = stops[i].ScheduledArrival != "" || stops[i].ScheduledDeparture != ""
		if stops[i].ScheduledArrival != "" {
			anyArr = true
		}
		if stops[i].ScheduledDeparture != "" {
			anyDep = true
		}
	}

	if !anyArr || !anyDep {
		app.renderStopsForm(w, r, train, csrf, stops,
			"Please enter at least one arrival time and one departure time. Leave a station blank if the train doesn't stop there.")
		return
	}

	// Apply: upsert stops that have a time, delete those left blank.
	for _, s := range stops {
		if s.Stops {
			var arrVal, depVal interface{}
			if s.ScheduledArrival != "" {
				arrVal = s.ScheduledArrival
			}
			if s.ScheduledDeparture != "" {
				depVal = s.ScheduledDeparture
			}
			app.db.Exec(
				`INSERT INTO train_stops (train_id, stop_id, sort_order, scheduled_arrival, scheduled_departure, runs_weekday, runs_weekend)
				 VALUES (?, ?, ?, ?, ?, ?, ?)
				 ON CONFLICT(train_id, stop_id) DO UPDATE SET
				   scheduled_arrival=excluded.scheduled_arrival,
				   scheduled_departure=excluded.scheduled_departure,
				   runs_weekday=excluded.runs_weekday,
				   runs_weekend=excluded.runs_weekend`,
				train.ID, s.StopID, s.SortOrder, arrVal, depVal, boolToInt(s.RunsWeekday), boolToInt(s.RunsWeekend),
			)
		} else {
			app.db.Exec(`DELETE FROM train_stops WHERE train_id=? AND stop_id=?`, train.ID, s.StopID)
		}
	}

	setFlash(w, "Schedule updated.")
	http.Redirect(w, r, "/trains/"+train.Slug+"/edit/stops", http.StatusSeeOther)
}
