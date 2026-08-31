package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// ACE (Altamont Corridor Express) live positions, sourced from 511.org's
// shared SF Bay regional GTFS-RT feed (see bay511.go) — not an ACE-specific
// request, so enabling Caltrain alongside this costs no extra 511.org calls.
// Agency code "CE" confirmed against 511.org's /transit/gtfsoperators list;
// the same free 511.org API key used for Caltrain works here too. ACE runs
// only ~8-10 weekday peak-commute trains, so this feed is legitimately empty
// most of the day — zero entities is expected, not a bug.
//
// Unlike Caltrain, ACE's native GTFS-RT trip_id (after bay511Vehicles strips
// the regional feed's "CE:" prefix) is not the bare train number: static
// GTFS trips.txt shows ids like "ACE01", "ACE03H", "ACE01fifa2" (holiday and
// special-event schedule variants of the same physical train), confirmed
// against a downloaded copy of ACE's static feed. The digits right after the
// "ACE" prefix are the train number regardless of any trailing variant
// suffix, so aceTripIDPattern strips exactly that.
var aceTripIDPattern = regexp.MustCompile(`^ACE(\d+)`)

type aceSource struct{}

func (aceSource) Key() string       { return "ace" }
func (aceSource) NeedsAPIKey() bool { return true }
func (aceSource) Description() string {
	return "Position data comes from the 511.org SF Bay Open Data GTFS-Realtime feed (shared with Caltrain — the same API key works for both, and enabling both costs no extra requests). ACE runs only a handful of weekday peak-commute trains, so it's normal to see nothing outside those windows."
}

func (aceSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "ace")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return nil, fmt.Errorf("no 511.org API key configured")
	}

	index, err := loadDBTrainsByCorridorSlug(app, "ace")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active ACE trains in the database to match against")
	}

	feed, err := fetchBay511Regional(src.APIKey)
	if err != nil {
		return nil, err
	}

	var out []liveTrain
	for _, v := range bay511Vehicles(feed, "CE") {
		trainNum := ""
		if m := aceTripIDPattern.FindStringSubmatch(v.TripID); m != nil {
			trainNum = m[1]
		}
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

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// ACE's operator_id on 511.org's static-GTFS datafeeds endpoint is "CE",
// confirmed against the same /transit/gtfsoperators list used for the live
// regional feed's agency prefix (bay511.go).
var aceRouteLineCache routeLineCache

func fetchACERouteGeoJSON(app *App) ([]byte, error) {
	src, err := getLiveSource(app.db, "ace")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return emptyRouteGeoJSON(), nil
	}
	shapes, err := fetch511StaticShapes(src.APIKey, "CE")
	if err != nil {
		return nil, err
	}
	return shapesToFeatureCollection(shapes, "ACE"), nil
}

func (app *App) handleACERoutes(w http.ResponseWriter, r *http.Request) {
	aceRouteLineCache.handle(w, r, func() ([]byte, error) { return fetchACERouteGeoJSON(app) })
}

// fetchGTFSRTBody does a plain unauthenticated-or-not GET and returns the raw
// body, shared by every direct (non-regional) GTFS-RT source (LIRR) and by
// the regional Bay Area fetch itself (bay511.go).
func fetchGTFSRTBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}
