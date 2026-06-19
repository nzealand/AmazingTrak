package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

// Admin conductor management: review role requests and assign/change/remove the
// Conductor of a corridor. Gated at permission level 4 (corridors) in main.go.

// handleAdminConductors renders pending requests plus a per-corridor assignment
// table with a searchable (datalist-backed) user picker.
func (app *App) handleAdminConductors(w http.ResponseWriter, r *http.Request) {
	s := sessionFromCtx(r)
	pending, err := conductorRequestsByStatus(app.db, "pending")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	corridors, err := allCorridors(app.db, false)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	users, err := allUsers(app.db)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	type conductorsData struct {
		Pending   []ConductorRequest
		Corridors []Corridor
		Users     []User
	}
	app.renderAdmin(w, r, "conductors.html", adminPage{
		Title:     "Conductors",
		Flash:     getFlash(w, r),
		CSRFToken: s.CSRFToken,
		Data:      conductorsData{Pending: pending, Corridors: corridors, Users: users},
	})
}

// handleAdminConductorApprove approves a request: assigns the user as conductor
// (replacing any existing one) and auto-rejects the corridor's other pending
// requests.
func (app *App) handleAdminConductorApprove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	req, err := conductorRequestByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	s := sessionFromCtx(r)
	if err := setCorridorConductor(app.db, req.CorridorID, req.UserID); err != nil {
		setFlash(w, "Error assigning conductor: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/conductors", http.StatusSeeOther)
		return
	}
	decideConductorRequest(app.db, id, "approved", s.AdminUserID)
	rejectOtherPendingRequests(app.db, req.CorridorID, id, s.AdminUserID)
	app.logAudit(s.AdminUserID, "conductor_approve", "corridor", req.CorridorID, req.Username)
	setFlash(w, req.Username+" is now the Conductor of "+req.CorridorName+".")
	http.Redirect(w, r, app.adminPrefix+"/conductors", http.StatusSeeOther)
}

// handleAdminConductorReject rejects a single pending request.
func (app *App) handleAdminConductorReject(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	req, err := conductorRequestByID(app.db, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	s := sessionFromCtx(r)
	decideConductorRequest(app.db, id, "rejected", s.AdminUserID)
	app.logAudit(s.AdminUserID, "conductor_reject", "corridor", req.CorridorID, req.Username)
	setFlash(w, "Request rejected.")
	http.Redirect(w, r, app.adminPrefix+"/conductors", http.StatusSeeOther)
}

// handleAdminCorridorConductorSet directly assigns a corridor's conductor by
// username (from the datalist picker).
func (app *App) handleAdminCorridorConductorSet(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	corridorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	corridor, err := corridorByID(app.db, corridorID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	backURL := app.adminPrefix + "/conductors"
	username := r.FormValue("username")
	user, err := userByUsername(app.db, username)
	if err == sql.ErrNoRows {
		setFlash(w, "No registered user named "+username+".")
		http.Redirect(w, r, backURL, http.StatusSeeOther)
		return
	}
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	s := sessionFromCtx(r)
	if err := setCorridorConductor(app.db, corridorID, user.ID); err != nil {
		setFlash(w, "Error assigning conductor: "+err.Error())
		http.Redirect(w, r, backURL, http.StatusSeeOther)
		return
	}
	rejectOtherPendingRequests(app.db, corridorID, 0, s.AdminUserID)
	app.logAudit(s.AdminUserID, "conductor_set", "corridor", corridorID, user.Username)
	setFlash(w, user.Username+" is now the Conductor of "+corridor.Name+".")
	http.Redirect(w, r, backURL, http.StatusSeeOther)
}

// handleAdminCorridorConductorRemove clears a corridor's conductor.
func (app *App) handleAdminCorridorConductorRemove(w http.ResponseWriter, r *http.Request) {
	if !app.checkCSRF(r) {
		http.Error(w, "Invalid CSRF token", 403)
		return
	}
	corridorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	s := sessionFromCtx(r)
	if err := clearCorridorConductor(app.db, corridorID); err != nil {
		setFlash(w, "Error removing conductor: "+err.Error())
		http.Redirect(w, r, app.adminPrefix+"/conductors", http.StatusSeeOther)
		return
	}
	app.logAudit(s.AdminUserID, "conductor_remove", "corridor", corridorID, "")
	setFlash(w, "Conductor removed.")
	http.Redirect(w, r, app.adminPrefix+"/conductors", http.StatusSeeOther)
}
