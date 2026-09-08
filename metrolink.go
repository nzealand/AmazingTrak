package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// LA Metrolink (Southern California commuter rail) live positions, sourced
// from Metrolink's own official GTFS-Realtime feed
// (https://metrolinktrains.com/about/gtfs/gtfs-rt-access/, confirmed
// 2026-09-08). Requires a free API key requested through the form linked on
// that page, sent as a query param — entered by an admin on the Settings
// page, same as every other keyed source in this app (live_sources.api_key).
// Updated every 30s upstream; polling faster buys nothing.
//
// Unlike Metra, all 7 Metrolink lines (Antelope Valley, Inland
// Empire-Orange County, Orange County, Riverside, San Bernardino, Ventura
// County, 91 Line) share one train-number space with zero collisions
// (confirmed against every one of 319 rows in the static GTFS's trips.txt —
// see trains_metrolink.go), so this is modeled as a single "la-metrolink"
// corridor, matching the NJ Transit/SEPTA idiom instead of Metra's
// one-corridor-per-line split.
//
// Unlike Metra, whose GTFS-RT trip_id embeds the train number in a
// regex-matchable pattern, Metrolink's trip_id is an opaque numeric
// static-GTFS id (e.g. "296000081") bearing no relation to the public train
// number — matching instead requires a trip_id -> train number (trips.txt's
// trip_short_name) lookup table built from the static schedule, refreshed
// periodically (see metrolinkTripMap below), the same static-file-derived
// idiom already used for Metra's route-line shape_id -> route_id mapping,
// just applied to live matching instead of geometry.
const (
	metrolinkVehiclePositionsURL = "https://metrolink-gtfsrt.gbsdigital.us/feed/gtfsrt-vehicles"
	metrolinkStaticGTFSURL       = "https://metrolinktrains.com/globalassets/about/gtfs/gtfs.zip"
	// Metrolink's static schedule changes only a handful of times a year, so
	// re-downloading it every live poll (every 30-90s) would be wasteful —
	// this is refreshed independently, far less often.
	metrolinkTripMapTTL = 6 * time.Hour
	// A VehiclePosition entity older than this is dropped rather than shown
	// as current — same defensive guard as Metra/LIRR, since Metrolink's docs
	// don't explicitly promise stale entities are pruned from the feed.
	metrolinkEntityMaxAge = 20 * time.Minute
)

type metrolinkSource struct{}

func (metrolinkSource) Key() string       { return "metrolink" }
func (metrolinkSource) NeedsAPIKey() bool { return true }
func (metrolinkSource) Description() string {
	return `Position data comes from Metrolink's own official <a href="https://metrolinktrains.com/about/gtfs/gtfs-rt-access/" target="_blank" rel="noopener">GTFS-Realtime feed</a>. Requires a free API key requested through the form on that page. Covers all 7 Metrolink lines. Delay/next-station info isn't available from this source yet — only position, speed, and heading.`
}

// metrolinkTripMapCache holds the static trip_id -> train number lookup,
// refreshed at most every metrolinkTripMapTTL.
var metrolinkTripMapCache struct {
	mu        sync.Mutex
	byTripID  map[string]string
	fetchedAt time.Time
	err       error
}

// fetchMetrolinkTripMap returns the cached trip_id -> train number map,
// downloading and parsing Metrolink's static GTFS zip (no key needed) only
// when the cached copy has aged past metrolinkTripMapTTL.
func fetchMetrolinkTripMap() (map[string]string, error) {
	metrolinkTripMapCache.mu.Lock()
	defer metrolinkTripMapCache.mu.Unlock()

	if metrolinkTripMapCache.byTripID != nil && time.Since(metrolinkTripMapCache.fetchedAt) < metrolinkTripMapTTL {
		return metrolinkTripMapCache.byTripID, metrolinkTripMapCache.err
	}

	data, err := downloadBytes(metrolinkStaticGTFSURL)
	if err != nil {
		metrolinkTripMapCache.err = err
		return metrolinkTripMapCache.byTripID, err
	}
	m, err := parseMetrolinkTripShortNames(data)
	metrolinkTripMapCache.fetchedAt = time.Now()
	if err != nil {
		metrolinkTripMapCache.err = err
		return metrolinkTripMapCache.byTripID, err
	}
	metrolinkTripMapCache.byTripID = m
	metrolinkTripMapCache.err = nil
	return m, nil
}

// parseMetrolinkTripShortNames extracts trip_id -> trip_short_name (the
// public train number) from a downloaded copy of Metrolink's static GTFS
// zip's trips.txt.
func parseMetrolinkTripShortNames(data []byte) (map[string]string, error) {
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
	tripIdx, tripOK := col["trip_id"]
	numIdx, numOK := col["trip_short_name"]
	if !tripOK || !numOK {
		return nil, fmt.Errorf("trips.txt missing trip_id or trip_short_name column")
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
		if tripIdx >= len(rec) || numIdx >= len(rec) {
			continue
		}
		tripID := strings.TrimSpace(rec[tripIdx])
		num := strings.TrimSpace(rec[numIdx])
		if tripID == "" || num == "" {
			continue
		}
		out[tripID] = num
	}
	return out, nil
}

func (metrolinkSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "metrolink")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return nil, fmt.Errorf("no Metrolink API key configured")
	}

	index, err := loadDBTrainsByCorridorSlug(app, "la-metrolink")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active LA Metrolink trains in the database to match against")
	}

	tripMap, err := fetchMetrolinkTripMap()
	if err != nil {
		return nil, err
	}

	reqURL := metrolinkVehiclePositionsURL + "?key=" + url.QueryEscape(src.APIKey)
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
		if ts := vp.GetTimestamp(); ts > 0 && now.Sub(time.Unix(int64(ts), 0)) > metrolinkEntityMaxAge {
			continue
		}

		trainNum := tripMap[vp.GetTrip().GetTripId()]
		if trainNum == "" {
			// Fall back to the vehicle label on the off chance it's ever the
			// train number rather than an equipment id — matches other
			// sources' defensive fallback idiom (e.g. Caltrain).
			trainNum = vp.GetVehicle().GetLabel()
		}
		if trainNum == "" {
			continue
		}
		match, ok := index[trainNum]
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

var metrolinkRouteLineCache routeLineCache

func fetchMetrolinkRouteGeoJSON() ([]byte, error) {
	data, err := downloadBytes(metrolinkStaticGTFSURL)
	if err != nil {
		return nil, err
	}
	shapes, err := parseGTFSShapesZip(data)
	if err != nil {
		return nil, err
	}
	if len(shapes) == 0 {
		return emptyRouteGeoJSON(), nil
	}
	return shapesToFeatureCollection(shapes, "LA Metrolink"), nil
}

func (app *App) handleMetrolinkRoutes(w http.ResponseWriter, r *http.Request) {
	metrolinkRouteLineCache.handle(w, r, fetchMetrolinkRouteGeoJSON)
}
