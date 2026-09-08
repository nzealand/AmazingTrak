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
// Only the Quebec City-Windsor Corridor is tracked (train numbers in
// viaCorridorTrainNumbers, trains_via.go) — VIA's remote long-distance
// services are out of scope by deliberate choice, not a technical
// limitation; see that file and stations_via.go for why.
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
	return `Position data comes from an <strong>unofficial</strong> source: the same JSON endpoint VIA Rail's own public train tracker (tsimobile.viarail.ca) calls in the browser. VIA publishes no official real-time developer API, only static schedules, so this could change or stop working without notice. Covers the Quebec City-Windsor Corridor only (Toronto/Ottawa/Montreal/Quebec area) — VIA's remote long-distance services (The Canadian, the Ocean, etc.) aren't tracked.`
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
	index, err := loadDBTrainsByCorridorSlug(app, "via-rail-corridor")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active VIA Rail Corridor trains in the database to match against")
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
// VIA's static GTFS zip (same source as trains_via.go/stations_via.go's
// rosters, no key needed) has its own shapes.txt. Like Metra (metra.go) and
// unlike single-shape sources, this needs each shape's own route_id
// (trips.txt) to know which ones belong to the Corridor — VIA's
// long-distance shapes (The Canadian, the Ocean, etc.) live in the same
// static feed and must be excluded, same reasoning as trains_via.go's
// roster scoping. Unlike Metra, every included shape belongs to the same
// one corridor here, so no per-shape corridor tag is needed —
// shapesToFeatureCollection's single shared name is enough.
const viaStaticGTFSURL = "https://viarail.ca/sites/all/files/gtfs/viarail.zip"

// viaCorridorRouteIDs are the 9 static-GTFS route_ids that make up the
// Quebec City-Windsor Corridor — see stations_via.go's comment for the full
// list of routes and why the rest of VIA's network is excluded.
var viaCorridorRouteIDs = map[string]bool{
	"119-93":  true, // Toronto - London
	"226-119": true, // Montréal - Toronto
	"119-341": true, // Toronto - Sarnia
	"617-628": true, // Ottawa - Québec
	"119-618": true, // Toronto - Windsor
	"617-226": true, // Ottawa - Montréal
	"617-119": true, // Ottawa - Toronto
	"628-576": true, // Québec - Fallowfield
	"628-226": true, // Québec - Montréal
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

	corridorShapes := make(map[string][][2]float64, len(shapes))
	for shapeID, coords := range shapes {
		if viaCorridorRouteIDs[shapeRoutes[shapeID]] {
			corridorShapes[shapeID] = coords
		}
	}
	if len(corridorShapes) == 0 {
		return emptyRouteGeoJSON(), nil
	}
	return shapesToFeatureCollection(corridorShapes, "VIA Rail Corridor"), nil
}

func (app *App) handleVIARoutes(w http.ResponseWriter, r *http.Request) {
	viaRouteLineCache.handle(w, r, fetchVIARouteGeoJSON)
}
