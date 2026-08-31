package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// LIRR (Long Island Rail Road) live positions, sourced from the MTA's public
// GTFS-Realtime feed. Confirmed empirically: unlike Caltrain/ACE, this feed
// requires no API key at all (MTA opened these feeds up publicly some years
// back) — NeedsAPIKey is false and no key field is ever sent.
//
// Confirmed against a downloaded copy of LIRR's static GTFS (trips.txt):
// trip_id looks like "GO201_26_809" or "GO201_26_452_2931_METS" (a schedule
// generation id, a sub-id, then the train number, occasionally followed by a
// special-service suffix); trip_short_name — the public train number — is
// exactly the first digit run after the two GO<n>_<n>_ segments in every one
// of 3620 rows checked, hence lirrTripIDPattern. vehicle.label, like MBTA, is
// the physical railcar number, not the train number — not used for matching.
//
// Also confirmed empirically: this feed does not drop stale entities — some
// VehiclePosition entries can be hours old even in an otherwise-fresh poll
// (e.g. a train that finished its trip and hasn't been rotated out yet), so
// each entity's own Timestamp is checked against lirrEntityMaxAge rather
// than trusting the feed as a whole to only contain current data.
const lirrVehiclePositionsURL = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/lirr%2Fgtfs-lirr"

const lirrEntityMaxAge = 20 * time.Minute

var lirrTripIDPattern = regexp.MustCompile(`^GO\d+_\d+_(\d+)`)

type lirrSource struct{}

func (lirrSource) Key() string       { return "lirr" }
func (lirrSource) NeedsAPIKey() bool { return false }
func (lirrSource) Description() string {
	return "Position data comes from the MTA's public LIRR GTFS-Realtime feed (no API key needed). Trains only — this does not and will not cover the NYC subway."
}

func (lirrSource) Fetch(app *App) ([]liveTrain, error) {
	index, err := loadDBTrainsByCorridorSlug(app, "lirr")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active LIRR trains in the database to match against")
	}

	body, err := fetchGTFSRTBody(lirrVehiclePositionsURL)
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
		if ts := vp.GetTimestamp(); ts > 0 && now.Sub(time.Unix(int64(ts), 0)) > lirrEntityMaxAge {
			continue
		}

		m := lirrTripIDPattern.FindStringSubmatch(vp.GetTrip().GetTripId())
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
		out = append(out, lt)
	}
	return out, nil
}
