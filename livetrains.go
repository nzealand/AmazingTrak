package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Live train positions, sourced from the Amtraker v3 API (a free, community-run
// mirror of Amtrak's own track-a-train feed). The feature is off unless an admin
// enables it on the Settings page; while off we never call out to the API.
//
// Amtraker refreshes roughly every 90s and Amtrak's own GPS pings lag 1-5
// minutes behind reality, so polling faster than this buys nothing but load.
const (
	amtrakerURL      = "https://api-v3.amtraker.com/v3/trains"
	livePollInterval = 90 * time.Second
	// A snapshot older than this is withheld rather than shown as current — if
	// the upstream API dies we would otherwise leave stale trains frozen on the
	// map, which is worse than showing none.
	liveMaxAge = 10 * time.Minute
)

// amtrakerTrain is the subset of the upstream payload we consume.
type amtrakerTrain struct {
	TrainNum   string  `json:"trainNum"`
	RouteName  string  `json:"routeName"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Heading    string  `json:"heading"`
	Velocity   float64 `json:"velocity"`
	TrainState string  `json:"trainState"`
	Stations   []struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		SchArr string `json:"schArr"`
		Arr    string `json:"arr"`
		Status string `json:"status"`
	} `json:"stations"`
}

// liveTrain is what we hand to the map. Every one of these corresponds to a
// train in our own DB, so TrainSlug always links somewhere real.
type liveTrain struct {
	TrainNum     string  `json:"trainNum"`
	DisplayName  string  `json:"displayName"`
	TrainSlug    string  `json:"trainSlug"`
	CorridorName string  `json:"corridorName"`
	CorridorSlug string  `json:"corridorSlug"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Heading      string  `json:"heading"`
	Speed        int     `json:"speed"`
	DelayMin     int     `json:"delayMin"`
	Status       string  `json:"status"` // ontime | late | verylate
	NextStation  string  `json:"nextStation"`
}

// liveSnapshot is one poll's worth of matched trains.
type liveSnapshot struct {
	Trains    []liveTrain `json:"trains"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// dbTrain is one row of our trains table, used to match a live train to our content.
type dbTrain struct {
	Slug         string
	DisplayName  string
	CorridorName string
	CorridorSlug string
	normRoute    string // normalized corridor name, for disambiguation
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// normRouteName reduces a route or corridor name to a comparable key:
// "Amtrak Cascades", "Cascades" and "cascades" all collapse to "cascades".
// Upstream route names do not always match ours ("Northest Regional" is a typo
// in the feed; we combine "Carl Sandburg / Illinois Zephyr" where they split
// it), so this is a best-effort match backed by the train-number fallback below.
func normRouteName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "amtrak ")
	return nonAlnum.ReplaceAllString(s, "")
}

// routeAliases maps upstream route names we can't derive to the corridor we
// keep them under. Keys and values are already normalized.
var routeAliases = map[string]string{
	"northestregional":    "northeastregional", // upstream typo
	"carlsandburg":        "carlsandburgillinoiszephyr",
	"illinoiszephyr":      "carlsandburgillinoiszephyr",
	"lincolnservice":      "lincolnmissouririverrunner",
	"missouririverrunner": "lincolnmissouririverrunner",
	"lincolnriverrunner":  "lincolnmissouririverrunner",
	"illini":              "illinisaluki",
	"saluki":              "illinisaluki",
}

// loadDBTrains indexes our trains by train number. A number is unique within a
// corridor but not across them, so a number can map to several candidates.
func loadDBTrains(app *App) (map[string][]dbTrain, error) {
	rows, err := app.db.Query(`
		SELECT t.train_number, t.slug, t.display_name, c.name, c.slug
		FROM trains t JOIN corridors c ON c.id = t.corridor_id
		WHERE t.is_active = 1 AND c.is_active = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string][]dbTrain{}
	for rows.Next() {
		var num string
		var t dbTrain
		if err := rows.Scan(&num, &t.Slug, &t.DisplayName, &t.CorridorName, &t.CorridorSlug); err != nil {
			return nil, err
		}
		t.normRoute = normRouteName(t.CorridorName)
		out[num] = append(out[num], t)
	}
	return out, rows.Err()
}

// matchTrain resolves a live train to one of ours, or reports false. When a
// train number is ambiguous across corridors we disambiguate on route name;
// when it is unique we accept it regardless of what the feed calls the route.
func matchTrain(index map[string][]dbTrain, lt amtrakerTrain) (dbTrain, bool) {
	candidates := index[lt.TrainNum]
	if len(candidates) == 0 {
		return dbTrain{}, false
	}
	if len(candidates) == 1 {
		return candidates[0], true
	}

	want := normRouteName(lt.RouteName)
	if alias, ok := routeAliases[want]; ok {
		want = alias
	}
	for _, c := range candidates {
		if c.normRoute == want {
			return c, true
		}
	}
	// Ambiguous number and no route match — better to drop it than to pin a
	// train onto the wrong corridor's page.
	return dbTrain{}, false
}

// delayFor returns minutes late (negative = early) at the next station the train
// has not yet departed, plus that station's name.
func delayFor(lt amtrakerTrain) (int, string) {
	for _, st := range lt.Stations {
		if st.Status == "Departed" {
			continue
		}
		if st.SchArr == "" || st.Arr == "" {
			return 0, st.Name
		}
		sch, err1 := time.Parse(time.RFC3339, st.SchArr)
		act, err2 := time.Parse(time.RFC3339, st.Arr)
		if err1 != nil || err2 != nil {
			return 0, st.Name
		}
		return int(act.Sub(sch).Minutes()), st.Name
	}
	return 0, ""
}

// delayStatus buckets lateness for map marker coloring.
func delayStatus(min int) string {
	switch {
	case min > 30:
		return "verylate"
	case min > 10:
		return "late"
	default:
		return "ontime"
	}
}

// fetchLiveTrains pulls the upstream feed and reduces it to the trains we track.
func fetchLiveTrains(app *App) ([]liveTrain, error) {
	index, err := loadDBTrains(app)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", amtrakerURL, nil)
	if err != nil {
		return nil, err
	}
	// Identify ourselves: this is a free community API and an anonymous poller
	// is a bad citizen.
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	// The feed is keyed by train number, each value an array of runs of that
	// number currently in service (a daily train can have several active runs).
	var raw map[string][]amtrakerTrain
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var out []liveTrain
	for _, runs := range raw {
		for _, lt := range runs {
			// Predeparture trains report their origin station's coordinates,
			// which would litter the map with trains that aren't moving yet.
			if lt.TrainState != "Active" {
				continue
			}
			if lt.Lat == 0 && lt.Lon == 0 {
				continue
			}
			match, ok := matchTrain(index, lt)
			if !ok {
				continue
			}
			delay, next := delayFor(lt)
			out = append(out, liveTrain{
				TrainNum:     lt.TrainNum,
				DisplayName:  match.DisplayName,
				TrainSlug:    match.Slug,
				CorridorName: match.CorridorName,
				CorridorSlug: match.CorridorSlug,
				Lat:          lt.Lat,
				Lon:          lt.Lon,
				Heading:      lt.Heading,
				Speed:        int(lt.Velocity + 0.5),
				DelayMin:     delay,
				Status:       delayStatus(delay),
				NextStation:  next,
			})
		}
	}
	return out, nil
}

// liveTrainsCache holds the most recent successful poll.
type liveTrainsCache struct {
	mu        sync.RWMutex
	trains    []liveTrain
	updatedAt time.Time
}

func (c *liveTrainsCache) store(trains []liveTrain) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.trains = trains
	c.updatedAt = time.Now()
}

// load returns the cached snapshot, and false if it is missing or too stale to
// be worth showing.
func (c *liveTrainsCache) load() (liveSnapshot, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.updatedAt.IsZero() || time.Since(c.updatedAt) > liveMaxAge {
		return liveSnapshot{}, false
	}
	return liveSnapshot{Trains: c.trains, UpdatedAt: c.updatedAt}, true
}

// findLiveTrain returns the current snapshot's entry for a train, matched by
// slug, if the feature is on and that train is currently running.
func (app *App) findLiveTrain(slug string) *liveTrain {
	if !app.liveTrainsEnabled() {
		return nil
	}
	snap, ok := app.liveTrains.load()
	if !ok {
		return nil
	}
	for i, t := range snap.Trains {
		if t.TrainSlug == slug {
			return &snap.Trains[i]
		}
	}
	return nil
}

// liveTrainsEnabled reports whether an admin has turned the feature on.
func (app *App) liveTrainsEnabled() bool {
	prefs, err := getSitePrefs(app.db)
	if err != nil {
		return false
	}
	return prefs.LiveTrainsEnabled
}

// pollLiveTrains refreshes the cache on a fixed interval for as long as the
// feature is enabled. Polling is skipped entirely while it is off, so a site
// that never turns this on never talks to the upstream API.
func (app *App) pollLiveTrains() {
	poll := func() {
		if !app.liveTrainsEnabled() {
			return
		}
		trains, err := fetchLiveTrains(app)
		if err != nil {
			// Keep serving the last good snapshot until it ages out.
			return
		}
		app.liveTrains.store(trains)
	}

	poll()
	for range time.Tick(livePollInterval) {
		poll()
	}
}

// handleLiveTrains serves the current snapshot to the map. When the feature is
// disabled the endpoint behaves as though it does not exist.
func (app *App) handleLiveTrains(w http.ResponseWriter, r *http.Request) {
	if !app.liveTrainsEnabled() {
		http.NotFound(w, r)
		return
	}
	snap, ok := app.liveTrains.load()
	if !ok {
		http.Error(w, "Live train data unavailable", 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Positions change every 90s upstream; a short cache absorbs bursts without
	// letting a browser show a visibly wrong position.
	w.Header().Set("Cache-Control", "public, max-age=30")
	json.NewEncoder(w).Encode(snap)
}
