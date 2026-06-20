package main

// Pending-items threshold notifications.
//
// "Pending review items" = pending suggestions + pending comments + unapproved
// registrations + pending conductor requests. The admin is emailed when the
// total crosses 1, 10, or 100 (upward). To avoid spamming, the highest threshold
// already notified is stored in site_preferences.pending_notify_level; a given
// threshold only fires again after the total has dropped back below it.

var pendingThresholds = []int{1, 10, 100}

// pendingItemsTotal returns the number of items awaiting admin review.
func (app *App) pendingItemsTotal() int {
	var suggestions, comments, registrations, conductorReqs int
	app.db.QueryRow(`SELECT COUNT(*) FROM suggestions WHERE status='pending'`).Scan(&suggestions)
	app.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE status='pending'`).Scan(&comments)
	app.db.QueryRow(`SELECT COUNT(*) FROM users WHERE status IN ('pending','confirmed')`).Scan(&registrations)
	app.db.QueryRow(`SELECT COUNT(*) FROM conductor_requests WHERE status='pending'`).Scan(&conductorReqs)
	return suggestions + comments + registrations + conductorReqs
}

// highestThresholdCrossed returns the largest threshold <= total, or 0.
func highestThresholdCrossed(total int) int {
	level := 0
	for _, t := range pendingThresholds {
		if total >= t {
			level = t
		}
	}
	return level
}

// maybeNotifyPending recomputes the pending total and emails the admin if it has
// crossed a threshold upward that hasn't already been notified. The stored level
// follows the total down too, so each threshold re-arms once the count drops
// below it. No-op when email is disabled or no notification address is set.
func (app *App) maybeNotifyPending() {
	prefs, err := getSitePrefs(app.db)
	if err != nil {
		return
	}
	total := app.pendingItemsTotal()
	current := highestThresholdCrossed(total)
	stored := prefs.PendingNotifyLevel

	if current == stored {
		return
	}
	// Persist the new level (covers both upward crossings and downward re-arming).
	app.db.Exec(`UPDATE site_preferences SET pending_notify_level=? WHERE id=1`, current)

	// Only email on an upward crossing, and only when email + a recipient exist.
	if current > stored && app.emailEnabled() && prefs.NotificationEmail != "" {
		go app.sendThresholdEmail(prefs.NotificationEmail, current, total)
	}
}
