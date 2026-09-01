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
// minutes behind reality, so polling faster than 90s buys nothing but load.
// The actual interval is admin-configurable (site_preferences.live_trains_poll_seconds,
// see liveTrainsPollInterval) between this floor and a 10-minute ceiling —
// matching other public trackers like transitdocs.com, which polls every 10m.
const (
	amtrakerURL          = "https://api-v3.amtraker.com/v3/trains"
	livePollIntervalMin  = 90 * time.Second
	livePollIntervalMax  = 10 * time.Minute
	livePollIntervalDflt = 120 * time.Second
	// A snapshot older than this is withheld rather than shown as current — if
	// the upstream API dies we would otherwise leave stale trains frozen on the
	// map, which is worse than showing none. Kept well above
	// livePollIntervalMax so a normal max-interval gap (or one slow poll)
	// never makes a fresh snapshot look stale.
	liveMaxAge = 15 * time.Minute
)

// amtrakerStation is one stop in a train's upstream station list, in route
// order. "Departed" means the train has left it; "Station" means the train
// is physically there right now; anything else ("Enroute") means still
// approaching. Arr/Dep are upstream's own continuously re-estimated
// predictions for that stop (not just the static schedule) — they settle
// down to the real observed time once the train actually arrives/departs.
type amtrakerStation struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	SchArr string `json:"schArr"`
	Arr    string `json:"arr"`
	SchDep string `json:"schDep"`
	Dep    string `json:"dep"`
	Status string `json:"status"`
}

// amtrakerTrain is the subset of the upstream payload we consume.
type amtrakerTrain struct {
	TrainNum   string            `json:"trainNum"`
	RouteName  string            `json:"routeName"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Heading    string            `json:"heading"`
	Velocity   float64           `json:"velocity"`
	TrainState string            `json:"trainState"`
	LastValTS  string            `json:"lastValTS"`
	Stations   []amtrakerStation `json:"stations"`
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
	// NextETA is Amtrak's own live-updated predicted arrival at NextStation
	// (upstream's "arr" field for a not-yet-departed station is a running
	// estimate, not just the static schedule) — RFC3339, empty if unknown.
	NextETA string `json:"nextEta,omitempty"`
	// NextStation2/NextETA2 are the station after NextStation, and its own
	// predicted arrival — so a client that's showing/guessing the train has
	// already reached NextStation can display what's coming up after it
	// instead of continuing to name a station it's already at.
	NextStation2 string `json:"nextStation2,omitempty"`
	NextETA2     string `json:"nextEta2,omitempty"`
	// CurrentStation/CurrentDeparture are set only when upstream confirms
	// (status "Station") the train is physically stopped there right now.
	// CurrentDeparture is upstream's own live-updated predicted departure —
	// same idea as NextETA, just for leaving instead of arriving.
	CurrentStation   string `json:"currentStation,omitempty"`
	CurrentDeparture string `json:"currentDeparture,omitempty"`
	// LastUpdated is when this train's position was last actually refreshed
	// upstream (RFC3339) — distinct from our own poll cadence.
	LastUpdated string `json:"lastUpdated,omitempty"`
	// HasDelayInfo is false for a source (e.g. Caltrain's VehiclePositions-only
	// feed, for now) that can't tell us how late a train is running. The map
	// UI must not render DelayMin/Status as "on time" — that's just the zero
	// value — unless this is true.
	HasDelayInfo bool `json:"hasDelayInfo"`
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

// stationProgress is what stationsInfo derives from a train's upstream
// station list: how late it's running, where it's confirmed stopped right
// now (if anywhere), and the next one or two stations still ahead.
type stationProgress struct {
	delayMin         int
	currentStation   string // "" unless upstream confirms (status "Station") the train is here right now
	currentDeparture string // upstream's live-updated predicted departure from currentStation
	nextStation      string
	nextETA          string
	nextStation2     string
	nextETA2         string
}

// stationsInfo walks a train's upstream station list (route order) past any
// already-departed stops, then classifies what's left: a "Station"-status
// entry means the train is physically there right now (currentStation, not
// "next"), everything after that is upcoming. Delay is measured against
// whichever stop is immediately at/ahead of the train, matching how upstream
// itself frames lateness.
func stationsInfo(lt amtrakerTrain) stationProgress {
	var out stationProgress

	var upcoming []amtrakerStation
	for _, st := range lt.Stations {
		if st.Status != "Departed" {
			upcoming = append(upcoming, st)
		}
	}
	if len(upcoming) == 0 {
		return out
	}

	first := upcoming[0]
	if first.SchArr != "" && first.Arr != "" {
		sch, err1 := time.Parse(time.RFC3339, first.SchArr)
		act, err2 := time.Parse(time.RFC3339, first.Arr)
		if err1 == nil && err2 == nil {
			out.delayMin = int(act.Sub(sch).Minutes())
		}
	}

	idx := 0
	if first.Status == "Station" {
		out.currentStation = first.Name
		out.currentDeparture = first.Dep
		idx = 1
	}
	if idx < len(upcoming) {
		out.nextStation = upcoming[idx].Name
		out.nextETA = upcoming[idx].Arr
	}
	if idx+1 < len(upcoming) {
		out.nextStation2 = upcoming[idx+1].Name
		out.nextETA2 = upcoming[idx+1].Arr
	}
	return out
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
			prog := stationsInfo(lt)
			out = append(out, liveTrain{
				TrainNum:         lt.TrainNum,
				DisplayName:      match.DisplayName,
				TrainSlug:        match.Slug,
				CorridorName:     match.CorridorName,
				CorridorSlug:     match.CorridorSlug,
				Lat:              lt.Lat,
				Lon:              lt.Lon,
				Heading:          lt.Heading,
				Speed:            int(lt.Velocity + 0.5),
				DelayMin:         prog.delayMin,
				Status:           delayStatus(prog.delayMin),
				NextStation:      prog.nextStation,
				NextETA:          prog.nextETA,
				NextStation2:     prog.nextStation2,
				NextETA2:         prog.nextETA2,
				CurrentStation:   prog.currentStation,
				CurrentDeparture: prog.currentDeparture,
				LastUpdated:      lt.LastValTS,
				HasDelayInfo:     true,
			})
		}
	}
	return out, nil
}

// liveTrainsCache holds the most recent successful poll of one source.
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

// liveSource is one pluggable live-tracking data provider. Amtrak (Amtraker)
// and Caltrain (511.org GTFS-RT) both implement it; adding another agency
// later is a new file implementing this interface plus a live_sources config
// row (see db.go migration 4) — no changes needed here or to the map/API
// endpoints below, which are all source-agnostic.
type liveSource interface {
	// Key matches live_sources.source_key and is this source's cache key.
	Key() string
	// NeedsAPIKey tells the admin Settings page whether to show an API-key
	// field for this source.
	NeedsAPIKey() bool
	// Description is a short admin-facing note (may contain simple HTML,
	// e.g. a link) about where this source's data comes from and any known
	// limitations, shown on the Settings page.
	Description() string
	// Fetch pulls one poll's worth of trains, already matched to our own DB
	// trains, from this source's upstream feed.
	Fetch(app *App) ([]liveTrain, error)
}

// amtrakSource wraps the pre-existing Amtraker integration above.
type amtrakSource struct{}

func (amtrakSource) Key() string       { return "amtrak" }
func (amtrakSource) NeedsAPIKey() bool { return false }
func (amtrakSource) Description() string {
	return `Position data comes from the third-party <a href="https://amtraker.com/" target="_blank" rel="noopener">Amtraker</a> API; Amtrak's own GPS can lag several minutes behind, so treat positions as approximate.`
}
func (amtrakSource) Fetch(app *App) ([]liveTrain, error) { return fetchLiveTrains(app) }

// registeredLiveSources lists every live-tracking provider the app knows how
// to poll. Each one needs a matching row in live_sources (seeded in db.go)
// to be configurable from the admin Settings page.
var registeredLiveSources = []liveSource{
	amtrakSource{},
	caltrainSource{},
	aceSource{},
	mbtaSource{},
	lirrSource{},
	septaSource{},
	brightlineSource{},
	metraSource{},
	njtSource{},
}

// liveSourceByKey looks up a registered source by its live_sources.source_key,
// or nil if key doesn't match anything registered.
func liveSourceByKey(key string) liveSource {
	for _, src := range registeredLiveSources {
		if src.Key() == key {
			return src
		}
	}
	return nil
}

// newLiveTrainCaches allocates one cache per registered source.
func newLiveTrainCaches() map[string]*liveTrainsCache {
	m := make(map[string]*liveTrainsCache, len(registeredLiveSources))
	for _, src := range registeredLiveSources {
		m[src.Key()] = &liveTrainsCache{}
	}
	return m
}

// findLiveTrain returns the current snapshot's entry for a train, matched by
// slug, across every enabled source.
func (app *App) findLiveTrain(slug string) *liveTrain {
	for _, src := range registeredLiveSources {
		cache, ok := app.liveTrainCaches[src.Key()]
		if !ok {
			continue
		}
		snap, ok := cache.load()
		if !ok {
			continue
		}
		for i, t := range snap.Trains {
			if t.TrainSlug == slug {
				return &snap.Trains[i]
			}
		}
	}
	return nil
}

// mergedLiveSnapshot combines every source's current cache into one snapshot
// for the map. UpdatedAt is the oldest of the contributing sources' update
// times, so a client checking freshness isn't misled by one very fresh
// source masking a stale one. Returns false if no source has fresh data.
func (app *App) mergedLiveSnapshot() (liveSnapshot, bool) {
	var out liveSnapshot
	found := false
	for _, src := range registeredLiveSources {
		cache, ok := app.liveTrainCaches[src.Key()]
		if !ok {
			continue
		}
		snap, ok := cache.load()
		if !ok {
			continue
		}
		out.Trains = append(out.Trains, snap.Trains...)
		if !found || snap.UpdatedAt.Before(out.UpdatedAt) {
			out.UpdatedAt = snap.UpdatedAt
		}
		found = true
	}
	return out, found
}

// liveTrainsEnabled reports whether at least one live source is currently on.
func (app *App) liveTrainsEnabled() bool {
	var n int
	app.db.QueryRow(`SELECT COUNT(*) FROM live_sources WHERE enabled=1`).Scan(&n)
	return n > 0
}

// liveTrainsPollInterval returns the fastest configured interval among
// enabled sources, for the frontend's own poll cadence (LiveTrainsPollMs).
func (app *App) liveTrainsPollInterval() time.Duration {
	sources, err := getLiveSources(app.db)
	if err != nil {
		return livePollIntervalDflt
	}
	best := time.Duration(0)
	for _, s := range sources {
		if !s.Enabled {
			continue
		}
		d := clampPollInterval(s.PollSeconds)
		if best == 0 || d < best {
			best = d
		}
	}
	if best == 0 {
		return livePollIntervalDflt
	}
	return best
}

// clampPollInterval enforces [livePollIntervalMin, livePollIntervalMax]
// regardless of what's stored (defense in depth against a bad direct DB edit).
func clampPollInterval(seconds int) time.Duration {
	d := time.Duration(seconds) * time.Second
	if d < livePollIntervalMin {
		return livePollIntervalMin
	}
	if d > livePollIntervalMax {
		return livePollIntervalMax
	}
	return d
}

// pollLiveSource refreshes one source's cache for as long as it's enabled,
// re-checking its configured interval after every poll so an admin's change
// on the Settings page takes effect from the next cycle. Polling is skipped
// entirely while a source is off, so a site that never enables it never talks
// to that source's upstream API.
func (app *App) pollLiveSource(src liveSource) {
	cache := app.liveTrainCaches[src.Key()]
	poll := func() {
		s, err := getLiveSource(app.db, src.Key())
		if err != nil || !s.Enabled {
			return
		}
		trains, err := src.Fetch(app)
		recordLiveSourceResult(app.db, src.Key(), err)
		if err != nil {
			// Keep serving the last good snapshot until it ages out.
			return
		}
		cache.store(trains)
	}

	poll()
	for {
		s, err := getLiveSource(app.db, src.Key())
		interval := livePollIntervalDflt
		if err == nil {
			interval = clampPollInterval(s.PollSeconds)
		}
		time.Sleep(interval)
		poll()
	}
}

// handleLiveTrains serves the current merged snapshot to the map. When no
// source is enabled the endpoint behaves as though it does not exist.
func (app *App) handleLiveTrains(w http.ResponseWriter, r *http.Request) {
	if !app.liveTrainsEnabled() {
		http.NotFound(w, r)
		return
	}
	snap, ok := app.mergedLiveSnapshot()
	if !ok {
		http.Error(w, "Live train data unavailable", 503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Positions change every 90s or so upstream; a short cache absorbs bursts
	// without letting a browser show a visibly wrong position.
	w.Header().Set("Cache-Control", "public, max-age=30")
	json.NewEncoder(w).Encode(snap)
}

// handleLiveTrain serves the cached data for one train, matched by slug, so a
// client that only cares about a single selected train doesn't need to
// re-fetch and re-diff the whole snapshot every poll. Returns 404 for no
// source enabled, a stale cache, or a train not currently running — the
// client doesn't need to distinguish those cases.
func (app *App) handleLiveTrain(w http.ResponseWriter, r *http.Request) {
	t := app.findLiveTrain(r.PathValue("slug"))
	if t == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=30")
	json.NewEncoder(w).Encode(t)
}
