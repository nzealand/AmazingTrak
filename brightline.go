package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// Brightline live positions, sourced from Brightline's own public
// GTFS-Realtime feed. Confirmed empirically (2026-08-31): no API key needed,
// and unlike the position-only feeds (Caltrain/ACE/MBTA/LIRR), TripUpdates
// carries real per-stop delay data, so this source can legitimately set
// HasDelayInfo: true.
//
// Confirmed against a downloaded copy of Brightline's static GTFS
// (trips.txt, 391/391 rows matching): trip_id is
// "<train number>_BL_<n>_1" — the train number is reproduced verbatim as the
// first underscore-delimited segment, hence brightlineTripIDPattern. Only
// covers the Florida corridor (Miami-Orlando); Brightline West (Las
// Vegas-Rancho Cucamonga) is still under construction and has no feed.
const (
	brightlineVehiclePositionsURL = "http://feed.gobrightline.com/position_updates.pb"
	brightlineTripUpdatesURL      = "http://feed.gobrightline.com/trip_updates.pb"
)

var brightlineTripIDPattern = regexp.MustCompile(`^(\d+)_BL_`)

type brightlineSource struct{}

func (brightlineSource) Key() string       { return "brightline" }
func (brightlineSource) NeedsAPIKey() bool { return false }
func (brightlineSource) Description() string {
	return "Position and delay data come from Brightline's own public GTFS-Realtime feed (no API key needed). Covers the Florida corridor (Miami-Orlando) only."
}

func (brightlineSource) Fetch(app *App) ([]liveTrain, error) {
	index, err := loadDBTrainsByCorridorSlug(app, "brightline")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active Brightline trains in the database to match against")
	}

	// Trip delay is a bonus on top of position — if this leg fails, still
	// report positions rather than dropping the whole poll.
	delayByTrip, _ := fetchBrightlineDelays()

	body, err := fetchGTFSRTBody(brightlineVehiclePositionsURL)
	if err != nil {
		return nil, err
	}
	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

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
		tripID := vp.GetTrip().GetTripId()
		m := brightlineTripIDPattern.FindStringSubmatch(tripID)
		if m == nil {
			continue
		}
		trainNum := m[1]
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
		if delay, ok := delayByTrip[tripID]; ok {
			lt.DelayMin = delay
			lt.Status = delayStatus(delay)
			lt.HasDelayInfo = true
		}
		out = append(out, lt)
	}
	return out, nil
}

// fetchBrightlineDelays pulls the TripUpdates feed and returns each trip's
// delay in minutes, keyed by trip_id, from its first stop-time update.
func fetchBrightlineDelays() (map[string]int, error) {
	body, err := fetchGTFSRTBody(brightlineTripUpdatesURL)
	if err != nil {
		return nil, err
	}
	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, entity := range feed.Entity {
		tu := entity.GetTripUpdate()
		if tu == nil || len(tu.StopTimeUpdate) == 0 {
			continue
		}
		stu := tu.StopTimeUpdate[0]
		var delaySec int32
		if arr := stu.GetArrival(); arr != nil {
			delaySec = arr.GetDelay()
		} else if dep := stu.GetDeparture(); dep != nil {
			delaySec = dep.GetDelay()
		} else {
			continue
		}
		out[tu.GetTrip().GetTripId()] = int(delaySec / 60)
	}
	return out, nil
}

// ---- Route line geometry (for the public map's "Route lines" layer) ----

const brightlineStaticGTFSURL = "http://feed.gobrightline.com/bl_gtfs.zip"

var brightlineRouteLineCache routeLineCache

func fetchBrightlineRouteGeoJSON() ([]byte, error) {
	data, err := downloadBytes(brightlineStaticGTFSURL)
	if err != nil {
		return nil, err
	}
	shapes, err := parseGTFSShapesZip(data)
	if err != nil {
		return nil, err
	}
	return shapesToFeatureCollection(shapes, "Brightline"), nil
}

func (app *App) handleBrightlineRoutes(w http.ResponseWriter, r *http.Request) {
	brightlineRouteLineCache.handle(w, r, fetchBrightlineRouteGeoJSON)
}
