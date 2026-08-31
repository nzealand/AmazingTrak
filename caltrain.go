package main

import (
	"fmt"
)

// Caltrain live positions, sourced from 511.org's shared SF Bay regional
// GTFS-RT feed (see bay511.go) — not a Caltrain-specific request, so
// enabling ACE alongside this costs no extra 511.org calls. Requires a free
// API key (registered by an admin at https://511.org/open-data, entered on
// the Settings page — see live_sources.api_key).
//
// Confirmed empirically against a real payload: Caltrain's native GTFS-RT
// trip_id (after bay511Vehicles strips the regional feed's "CT:" prefix) is
// exactly the public-facing train number (e.g. trip_id "118" == Train 118),
// so matching is a plain lookup against trains.train_number within the
// Caltrain corridor — no cross-corridor disambiguation needed, unlike Amtrak
// where one number can appear in several corridors.
//
// v1 only uses VehiclePositions (position/speed/heading) — not TripUpdates
// (delay/next-station). TripUpdates support can be added later; see
// HasDelayInfo on liveTrain.
type caltrainSource struct{}

func (caltrainSource) Key() string       { return "caltrain" }
func (caltrainSource) NeedsAPIKey() bool { return true }
func (caltrainSource) Description() string {
	return "Position data comes from the 511.org SF Bay Open Data GTFS-Realtime feed (shared with ACE — see its row below). Delay/next-station info isn't available from this source yet — only position, speed, and heading."
}

func (caltrainSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "caltrain")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return nil, fmt.Errorf("no 511.org API key configured")
	}

	index, err := loadDBTrainsByCorridorSlug(app, "caltrain")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active Caltrain trains in the database to match against")
	}

	feed, err := fetchBay511Regional(src.APIKey)
	if err != nil {
		return nil, err
	}

	var out []liveTrain
	for _, v := range bay511Vehicles(feed, "CT") {
		trainNum := v.TripID
		if trainNum == "" {
			trainNum = v.VehicleLabel
		}
		if trainNum == "" {
			continue
		}
		match, ok := index[trainNum]
		if !ok {
			continue
		}
		out = append(out, bay511LiveTrain(v, match, trainNum))
	}
	return out, nil
}

// loadDBTrainsByCorridorSlug indexes active trains in one corridor by train
// number. Unlike Amtrak's cross-corridor matchTrain, a single-corridor source
// like Caltrain never has an ambiguous number, so no route-name
// disambiguation is needed.
func loadDBTrainsByCorridorSlug(app *App, slug string) (map[string]dbTrain, error) {
	rows, err := app.db.Query(`
		SELECT t.train_number, t.slug, t.display_name, c.name, c.slug
		FROM trains t JOIN corridors c ON c.id = t.corridor_id
		WHERE t.is_active = 1 AND c.is_active = 1 AND c.slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]dbTrain{}
	for rows.Next() {
		var num string
		var t dbTrain
		if err := rows.Scan(&num, &t.Slug, &t.DisplayName, &t.CorridorName, &t.CorridorSlug); err != nil {
			return nil, err
		}
		out[num] = t
	}
	return out, rows.Err()
}

func mpsToMph(mps float64) float64 { return mps / 0.44704 }

// bearingToCompass reduces a 0-360 degree compass bearing to the 8-point
// direction the map's COMPASS_ARROW lookup expects (matches Amtraker's own
// heading strings: N, NE, E, SE, S, SW, W, NW).
func bearingToCompass(deg float64) string {
	for deg < 0 {
		deg += 360
	}
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	idx := int(deg/45.0+0.5) % 8
	return dirs[idx]
}
