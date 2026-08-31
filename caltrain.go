package main

import (
	"fmt"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// Caltrain live positions, sourced from the 511.org SF Bay Open Data GTFS-RT
// feed (agency code "CT"). Unlike Amtraker, this is a standard GTFS-RT
// VehiclePositions feed (protobuf, not JSON) and requires a free API key
// (registered by an admin at https://511.org/open-data, entered on the
// Settings page — see live_sources.api_key).
//
// Confirmed empirically against a real payload: Caltrain's GTFS-RT trip_id
// is exactly the public-facing train number (e.g. trip_id "118" == Train
// 118), so matching is a plain lookup against trains.train_number within
// the Caltrain corridor — no cross-corridor disambiguation needed, unlike
// Amtrak where one number can appear in several corridors.
//
// v1 only calls VehiclePositions (position/speed/heading) — not TripUpdates
// (delay/next-station), to keep this source's cost to one 511.org request
// per poll. 511.org's default rate limit is 60 requests/hour *per API key*,
// shared across every source that key is used for — ACE (see acetrain.go)
// runs on the same 511.org account, so its request budget comes out of the
// same 60/hour; keep an eye on combined poll intervals rather than tuning
// each source in isolation. TripUpdates support (delay/ETA) can be added
// later as a second request per poll; see HasDelayInfo on liveTrain.
const caltrainVehiclePositionsURL = "https://api.511.org/transit/vehiclepositions?api_key=%s&agency=CT"

type caltrainSource struct{}

func (caltrainSource) Key() string       { return "caltrain" }
func (caltrainSource) NeedsAPIKey() bool { return true }
func (caltrainSource) Description() string {
	return "Position data comes from the 511.org SF Bay Open Data GTFS-Realtime feed. Delay/next-station info isn't available from this source yet — only position, speed, and heading."
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

	body, err := fetchGTFSRTBody(fmt.Sprintf(caltrainVehiclePositionsURL, src.APIKey))
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

		trainNum := vp.GetTrip().GetTripId()
		if trainNum == "" {
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
