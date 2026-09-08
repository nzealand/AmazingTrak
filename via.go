package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// VIA Rail (Canada) live positions, sourced from tsimobile.viarail.ca's own
// JSON endpoint — the same one VIA's own public train-tracking web app
// (https://tsimobile.viarail.ca/) calls, discovered by inspecting that
// page's network traffic and confirmed empirically 2026-09-08. This is NOT
// an official published API (unlike Amtraker, which at least wraps Amtrak's
// own feed and identifies itself as a developer-facing service) — VIA
// publishes no real-time developer API at all, only static GTFS schedules
// (https://viarail.ca/en/developer-resources). Treat this the same way as
// Amtraker: best-effort, could change or disappear without notice, and this
// app identifies itself via User-Agent as a matter of politeness even though
// no terms of use are published to follow.
//
// Covers all 8 VIA corridors — the Quebec City-Windsor Corridor
// (trains_via.go) plus the 7 remote/long-distance named services
// (trains_via_remote.go: The Canadian, the Ocean, Winnipeg-Churchill,
// Sudbury-White River, Jasper-Prince Rupert, Montreal-Jonquière,
// Montreal-Senneterre). Train numbers are confirmed unique across all 8
// (see trains_via_remote.go's header comment), so loadVIADBTrains below
// builds one flat global number index instead of needing Amtrak-style
// route-name disambiguation for ambiguous numbers.
//
// Confirmed empirically against the live endpoint: each entry is keyed by
// "<train number>" normally, or "<train number> (MM-DD)" when more than one
// physical run of that number is active at once (e.g. a multi-day
// long-distance train that departs before its previous run has finished) —
// viaTrainNumberFromKey strips the optional date suffix. A finished trip has
// arrived=true and null lat/lng and is skipped, same idiom as Amtraker's
// TrainState != "Active" check (see fetchLiveTrains). Unlike Amtraker, this
// feed provides no explicit "currently at this station" flag per stop, so
// CurrentStation/CurrentDeparture are intentionally left unset here — only
// NextStation/NextStation2 (the next couple of stops still ahead, found by
// comparing each stop's own timestamp to now) and DelayMin are computed.
const (
	viaLiveDataURL     = "https://tsimobile.viarail.ca/data/allData.json"
	viaKmhToMph        = 0.621371
	viaHTTPTimeout     = 30 * time.Second
	viaMaxResponseSize = 8 << 20
)

type viaStopTime struct {
	Station   string `json:"station"`
	Code      string `json:"code"`
	Scheduled string `json:"scheduled"`
	Estimated string `json:"estimated"`
	DiffMin   *int   `json:"diffMin"`
}

// viaTrainEntry is one train's current status, as returned by the endpoint.
// Lat/Lng are pointers because they're JSON null once a trip has finished.
type viaTrainEntry struct {
	Lat       *float64      `json:"lat"`
	Lng       *float64      `json:"lng"`
	Speed     float64       `json:"speed"` // km/h
	Direction float64       `json:"direction"`
	Departed  bool          `json:"departed"`
	Arrived   bool          `json:"arrived"`
	Times     []viaStopTime `json:"times"`
}

type viaSource struct{}

func (viaSource) Key() string       { return "via-rail" }
func (viaSource) NeedsAPIKey() bool { return false }
func (viaSource) Description() string {
	return `Position data comes from an <strong>unofficial</strong> source: the same JSON endpoint VIA Rail's own public train tracker (tsimobile.viarail.ca) calls in the browser. VIA publishes no official real-time developer API, only static schedules, so this could change or stop working without notice. Covers all of VIA's network: the Quebec City-Windsor Corridor plus the remote/long-distance services (The Canadian, the Ocean, Winnipeg-Churchill, Sudbury-White River, Jasper-Prince Rupert, Montreal-Jonquière, Montreal-Senneterre) — most of the latter run only a few times a week, so an empty map for them most days is expected, not a bug.`
}

// viaCorridorSlugs lists every corridor this source matches live trains
// into. One live_sources row (this one) covers all of them, same "one
// provider, many corridors" idiom Metra uses for its 11 lines.
var viaCorridorSlugs = []string{
	"via-rail-corridor",
	"via-the-canadian",
	"via-the-ocean",
	"via-winnipeg-churchill",
	"via-sudbury-white-river",
	"via-jasper-prince-rupert",
	"via-montreal-jonquiere",
	"via-montreal-senneterre",
}

// loadVIADBTrains merges the per-corridor indexes for every VIA corridor
// into one flat number -> train map. Safe as a flat merge (rather than
// needing Amtrak-style route-name disambiguation for a collision) because
// train numbers are confirmed unique across all 8 corridors — see
// trains_via_remote.go's header comment.
func loadVIADBTrains(app *App) (map[string]dbTrain, error) {
	out := map[string]dbTrain{}
	for _, slug := range viaCorridorSlugs {
		idx, err := loadDBTrainsByCorridorSlug(app, slug)
		if err != nil {
			return nil, err
		}
		for num, t := range idx {
			out[num] = t
		}
	}
	return out, nil
}

// viaTrainNumberFromKey strips the optional " (MM-DD)" run-disambiguation
// suffix the endpoint adds when more than one physical run of a train
// number is active at once.
func viaTrainNumberFromKey(key string) string {
	if i := strings.Index(key, " ("); i >= 0 {
		return key[:i]
	}
	return key
}

// viaStopTimestamp resolves one stop's best-available timestamp: the
// live-updated estimate, falling back to the static schedule when the
// estimate isn't a real timestamp (confirmed empirically: the endpoint uses
// the literal string "&mdash;" for a not-yet-estimated far-future stop).
func viaStopTimestamp(st viaStopTime) (time.Time, bool) {
	if t, err := time.Parse(time.RFC3339, st.Estimated); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, st.Scheduled); err == nil {
		return t, true
	}
	return time.Time{}, false
}

type viaProgress struct {
	delayMin               int
	nextStation, nextETA   string
	nextStation2, nextETA2 string
}

// viaProgressFor walks a train's stop list (route order) and finds the
// first stop whose own timestamp is still in the future — unlike Amtraker,
// this feed has no explicit per-stop "departed/at this station/enroute"
// status field, so recency is inferred from the timestamps themselves.
func viaProgressFor(entry viaTrainEntry) viaProgress {
	var out viaProgress
	now := time.Now()

	idx := -1
	for i, st := range entry.Times {
		t, ok := viaStopTimestamp(st)
		if !ok {
			continue
		}
		if t.After(now) {
			idx = i
			break
		}
	}
	if idx == -1 {
		return out
	}

	out.nextStation = entry.Times[idx].Station
	if t, ok := viaStopTimestamp(entry.Times[idx]); ok {
		out.nextETA = t.UTC().Format(time.RFC3339)
	}
	if entry.Times[idx].DiffMin != nil {
		out.delayMin = *entry.Times[idx].DiffMin
	}
	if idx+1 < len(entry.Times) {
		out.nextStation2 = entry.Times[idx+1].Station
		if t, ok := viaStopTimestamp(entry.Times[idx+1]); ok {
			out.nextETA2 = t.UTC().Format(time.RFC3339)
		}
	}
	return out
}

func (viaSource) Fetch(app *App) ([]liveTrain, error) {
	index, err := loadVIADBTrains(app)
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active VIA Rail trains in the database to match against")
	}

	req, err := http.NewRequest("GET", viaLiveDataURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	client := &http.Client{Timeout: viaHTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, viaMaxResponseSize))
	if err != nil {
		return nil, err
	}

	var raw map[string]viaTrainEntry
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var out []liveTrain
	for key, entry := range raw {
		if entry.Arrived || entry.Lat == nil || entry.Lng == nil {
			continue
		}
		trainNum := viaTrainNumberFromKey(key)
		match, ok := index[trainNum]
		if !ok {
			continue
		}

		prog := viaProgressFor(entry)
		out = append(out, liveTrain{
			TrainNum:     trainNum,
			DisplayName:  match.DisplayName,
			TrainSlug:    match.Slug,
			CorridorName: match.CorridorName,
			CorridorSlug: match.CorridorSlug,
			Lat:          *entry.Lat,
			Lon:          *entry.Lng,
			Heading:      bearingToCompass(entry.Direction),
			Speed:        int(entry.Speed*viaKmhToMph + 0.5),
			DelayMin:     prog.delayMin,
			Status:       delayStatus(prog.delayMin),
			NextStation:  prog.nextStation,
			NextETA:      prog.nextETA,
			NextStation2: prog.nextStation2,
			NextETA2:     prog.nextETA2,
			HasDelayInfo: true,
		})
	}
	return out, nil
}

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// VIA's static GTFS zip (same source as trains_via.go/trains_via_remote.go's
// rosters, no key needed) has its own shapes.txt. Like Metra (metra.go),
// this needs each shape's own route_id (trips.txt) to know which corridor it
// belongs to — unlike Metra, though, every one of VIA's route_ids is
// tracked under some corridor now (none are out of scope), so every shape in
// the feed ends up included, just tagged differently.
//
// viaRouteIDToTagName's values are deliberately plain-ASCII, English-only
// labels distinct from the corridor's cosmetic display name (e.g.
// "VIA Montreal Jonquiere", not "VIA Rail — Montréal–Jonquière") — this
// value becomes a GeoJSON feature's "name" property, which the map's
// jsSlugify() (templates/map.html) turns into a corridor slug to match
// against; jsSlugify only preserves [a-z0-9], so an accented character or
// em-dash would silently break the word boundary it sits next to (e.g.
// "Montréal" -> "montr-al"). Each value here is chosen so
// jsSlugify(value) reproduces that corridor's actual slug exactly (verified
// by inspection, same idiom as metraLineCorridorSlug's route-code ->
// corridor-slug mapping, just inverted).
const viaStaticGTFSURL = "https://viarail.ca/sites/all/files/gtfs/viarail.zip"

var viaRouteIDToTagName = map[string]string{
	// Quebec City-Windsor Corridor (trains_via.go) -> "via-rail-corridor"
	"119-93":  "VIA Rail Corridor",
	"226-119": "VIA Rail Corridor",
	"119-341": "VIA Rail Corridor",
	"617-628": "VIA Rail Corridor",
	"119-618": "VIA Rail Corridor",
	"617-226": "VIA Rail Corridor",
	"617-119": "VIA Rail Corridor",
	"628-576": "VIA Rail Corridor",
	"628-226": "VIA Rail Corridor",
	// Remote/long-distance services (trains_via_remote.go)
	"8-119":   "VIA The Canadian",         // -> "via-the-canadian"
	"226-620": "VIA The Ocean",            // -> "via-the-ocean"
	"388-435": "VIA Winnipeg Churchill",   // -> "via-winnipeg-churchill"
	"149-435": "VIA Winnipeg Churchill",   // -> "via-winnipeg-churchill"
	"621-116": "VIA Sudbury White River",  // -> "via-sudbury-white-river"
	"21-458":  "VIA Jasper Prince Rupert", // -> "via-jasper-prince-rupert"
	"226-444": "VIA Montreal Jonquiere",   // -> "via-montreal-jonquiere"
	"226-460": "VIA Montreal Senneterre",  // -> "via-montreal-senneterre"
}

var viaRouteLineCache routeLineCache

func fetchVIARouteGeoJSON() ([]byte, error) {
	data, err := downloadBytes(viaStaticGTFSURL)
	if err != nil {
		return nil, err
	}
	shapes, err := parseGTFSShapesZip(data)
	if err != nil {
		return nil, err
	}
	shapeRoutes, err := parseGTFSTripShapeRoutes(data)
	if err != nil {
		return nil, err
	}

	features := make([]map[string]interface{}, 0, len(shapes))
	for shapeID, coords := range shapes {
		if len(coords) < 2 {
			continue
		}
		tagName, ok := viaRouteIDToTagName[shapeRoutes[shapeID]]
		if !ok {
			continue
		}
		features = append(features, map[string]interface{}{
			"type":       "Feature",
			"properties": map[string]interface{}{"name": tagName},
			"geometry": map[string]interface{}{
				"type":        "LineString",
				"coordinates": coords,
			},
		})
	}
	if len(features) == 0 {
		return emptyRouteGeoJSON(), nil
	}
	result := map[string]interface{}{"type": "FeatureCollection", "features": features}
	return json.Marshal(result)
}

func (app *App) handleVIARoutes(w http.ResponseWriter, r *http.Request) {
	viaRouteLineCache.handle(w, r, fetchVIARouteGeoJSON)
}
