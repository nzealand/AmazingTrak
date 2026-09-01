package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// Metra (Chicago-area commuter rail) live positions, sourced from Metra's own
// GTFS-Realtime feed at https://gtfspublic.metrarr.com/gtfs/public/positions
// (metra.com/metra-gtfs-api, confirmed 2026-09-01). Requires a free API key
// from Metra's license-agreement request form, sent as an api_token query
// param — entered by an admin on the Settings page, same as every other
// keyed source in this app (live_sources.api_key). Updated every 30s
// upstream; polling faster buys nothing.
//
// Unlike every other agency in this app, Metra is modeled as 11 separate
// corridors (one per line: BNSF, Union Pacific North/Northwest/West,
// Milwaukee District North/West, Rock Island, Metra Electric, SouthWest
// Service, North Central Service, Heritage Corridor) instead of one shared
// corridor, because train numbers are NOT unique across lines (e.g. "300"
// exists on both Rock Island and UP-N) — this app's trains table enforces
// UNIQUE(corridor_id, train_number), and SEPTA already hit this exact
// collision problem trying to cram multiple lines into one corridor. See
// trains_metra.go / stations_metra.go for the seeded data and sourcing notes.
//
// Confirmed against every one of 8163 rows in Metra's static GTFS trips.txt:
// trip_id is "<ROUTE>_<LINEPREFIX><NUMBER>_V<version>_<variant>", e.g.
// "BNSF_BN1200_V4_A" (train 1200 on BNSF) or "UP-NW_UNW601_V3_A" (train 601
// on UP-NW) — the route code before the first underscore always matches
// trips.txt's own route_id column exactly, and the digit run right before
// "_V<version>" is always the public train number, regardless of the
// line-specific letter prefix (BN/HC/MN/MW/ME/NC/RI/SW/UN/UNW/UW) or
// variant-letter suffix length (A, AA, etc. — schedule/holiday variants of
// the same physical train, same idiom as ACE's trip_id suffixes).
// GTFS-RT's trip_id namespace is the same as the static feed's by spec, so
// this same regex is used both for live matching here and (via a Python
// port of the same pattern) to build the seeded roster in trains_metra.go.
var metraTripIDPattern = regexp.MustCompile(`^([A-Z-]+)_[A-Za-z]*(\d+)_V\d+_[A-Za-z]+$`)

// metraLineCorridorSlug maps a trip_id's route code (matches trips.txt's
// route_id exactly) to this app's corridor slug for that line.
var metraLineCorridorSlug = map[string]string{
	"BNSF":  "metra-bnsf",
	"UP-N":  "metra-up-n",
	"UP-NW": "metra-up-nw",
	"UP-W":  "metra-up-w",
	"MD-N":  "metra-md-n",
	"MD-W":  "metra-md-w",
	"RI":    "metra-ri",
	"ME":    "metra-me",
	"SWS":   "metra-sws",
	"NCS":   "metra-ncs",
	"HC":    "metra-hc",
}

const metraPositionsURL = "https://gtfspublic.metrarr.com/gtfs/public/positions"

// A VehiclePosition entity older than this is dropped rather than shown as
// current — same defensive guard as LIRR (mtalirr.go), since Metra's docs
// don't explicitly promise stale entities are pruned from the feed.
const metraEntityMaxAge = 20 * time.Minute

type metraSource struct{}

func (metraSource) Key() string       { return "metra" }
func (metraSource) NeedsAPIKey() bool { return true }
func (metraSource) Description() string {
	return `Position data comes from Metra's own <a href="https://metra.com/metra-gtfs-api" target="_blank" rel="noopener">GTFS-Realtime feed</a>. Requires a free API key from Metra's license-agreement request form (linked on that page). Covers all 11 Metra lines.`
}

func (metraSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "metra")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return nil, fmt.Errorf("no Metra API key configured")
	}

	index := make(map[string]map[string]dbTrain, len(metraLineCorridorSlug))
	for _, slug := range metraLineCorridorSlug {
		idx, err := loadDBTrainsByCorridorSlug(app, slug)
		if err != nil {
			return nil, err
		}
		index[slug] = idx
	}

	reqURL := metraPositionsURL + "?api_token=" + url.QueryEscape(src.APIKey)
	body, err := fetchGTFSRTBody(reqURL)
	if err != nil {
		return nil, err
	}

	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	now := time.Now()
	var out []liveTrain
	for _, entity := range feed.Entity {
		vp := entity.GetVehicle()
		if vp == nil {
			continue
		}
		pos := vp.GetPosition()
		if pos == nil {
			continue
		}
		if ts := vp.GetTimestamp(); ts > 0 && now.Sub(time.Unix(int64(ts), 0)) > metraEntityMaxAge {
			continue
		}

		m := metraTripIDPattern.FindStringSubmatch(vp.GetTrip().GetTripId())
		if m == nil {
			continue
		}
		routeCode, trainNum := m[1], m[2]
		slug, ok := metraLineCorridorSlug[routeCode]
		if !ok {
			continue
		}
		match, ok := index[slug][trainNum]
		if !ok {
			continue
		}

		lt := liveTrain{
			TrainNum:     trainNum,
			DisplayName:  match.DisplayName,
			TrainSlug:    match.Slug,
			CorridorName: match.CorridorName,
			CorridorSlug: match.CorridorSlug,
			Lat:          float64(pos.GetLatitude()),
			Lon:          float64(pos.GetLongitude()),
			Heading:      bearingToCompass(float64(pos.GetBearing())),
			Speed:        int(mpsToMph(float64(pos.GetSpeed())) + 0.5),
		}
		if ts := vp.GetTimestamp(); ts > 0 {
			lt.LastUpdated = time.Unix(int64(ts), 0).UTC().Format(time.RFC3339)
		}
		out = append(out, lt)
	}
	return out, nil
}

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// Metra's static schedule.zip (no key needed — see trains_metra.go) has its
// own shapes.txt, one shape per line/direction. Unlike every other agency's
// route-line fetch in this app (one corridor, so shapesToFeatureCollection's
// single shared name is enough), each shape here must be tagged with its
// own line's corridor so the map assigns it correctly — trips.txt's
// shape_id -> route_id mapping is what makes that possible, since shapes.txt
// itself has no route_id column. A feature's "name" property must slugify
// (via templates/map.html's jsSlugify) to exactly this app's corridor slug,
// e.g. "Metra BNSF" -> "metra-bnsf" — see metraLineCorridorSlug above; the
// route code alone (without the "Metra " prefix) would collide with any
// other agency happening to reuse the same short code.
const metraStaticGTFSURL = "https://schedules.metrarail.com/gtfs/schedule.zip"

var metraRouteLineCache routeLineCache

// parseGTFSTripShapeRoutes extracts trips.txt from a downloaded static GTFS
// zip and returns a shape_id -> route_id map (first trip seen for a given
// shape_id wins; every trip on one shape belongs to the same route in
// Metra's feed, so this is exact, not just a best guess).
func parseGTFSTripShapeRoutes(data []byte) (map[string]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening GTFS zip: %w", err)
	}
	var tripsFile *zip.File
	for _, f := range zr.File {
		if f.Name == "trips.txt" {
			tripsFile = f
			break
		}
	}
	if tripsFile == nil {
		return nil, fmt.Errorf("GTFS zip has no trips.txt")
	}
	rc, err := tripsFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	cr := csv.NewReader(rc)
	cr.TrimLeadingSpace = true
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	col := make(map[string]int, len(header))
	for i, h := range header {
		col[strings.TrimSpace(h)] = i
	}
	routeIdx, routeOK := col["route_id"]
	shapeIdx, shapeOK := col["shape_id"]
	if !routeOK || !shapeOK {
		return nil, fmt.Errorf("trips.txt missing route_id or shape_id column")
	}

	out := map[string]string{}
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if routeIdx >= len(rec) || shapeIdx >= len(rec) {
			continue
		}
		shapeID := strings.TrimSpace(rec[shapeIdx])
		if shapeID == "" {
			continue
		}
		if _, exists := out[shapeID]; exists {
			continue
		}
		out[shapeID] = strings.TrimSpace(rec[routeIdx])
	}
	return out, nil
}

func fetchMetraRouteGeoJSON() ([]byte, error) {
	data, err := downloadBytes(metraStaticGTFSURL)
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
		routeID, ok := shapeRoutes[shapeID]
		if !ok {
			continue
		}
		features = append(features, map[string]interface{}{
			"type":       "Feature",
			"properties": map[string]interface{}{"name": "Metra " + routeID},
			"geometry": map[string]interface{}{
				"type":        "LineString",
				"coordinates": coords,
			},
		})
	}
	result := map[string]interface{}{"type": "FeatureCollection", "features": features}
	return json.Marshal(result)
}

func (app *App) handleMetraRoutes(w http.ResponseWriter, r *http.Request) {
	metraRouteLineCache.handle(w, r, fetchMetraRouteGeoJSON)
}
