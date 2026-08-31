package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// bay511 is a shared poller for 511.org's SF Bay regional GTFS-RT feed
// (agency=RG), which returns every real-time-enabled Bay Area agency's
// vehicle positions in a single request. Caltrain and ACE both read from
// this one cached fetch instead of each hitting 511.org separately, so
// enabling both costs exactly the same one request per poll as enabling
// just one — comfortably inside 511.org's 60 requests/hour per-key limit
// (a request every bay511MinInterval = 90s is 40/hour) no matter how many
// Bay Area sources are turned on, and any future Bay Area source added here
// rides the same fetch for free.
//
// Confirmed empirically against a real payload: this combined feed
// namespaces each entity's route_id (and trip_id) as "<agencyCode>:<native
// id>" — e.g. Caltrain trip "124" appears here as "CT:124", not plain "124"
// as it does on Caltrain's own single-agency endpoint. bay511Vehicles strips
// that prefix so callers see exactly the native id their own matching logic
// (already built and verified against each single-agency feed) expects.
const bay511RegionalURL = "https://api.511.org/transit/vehiclepositions?api_key=%s&agency=RG"

const bay511MinInterval = 90 * time.Second

var bay511Cache struct {
	mu        sync.Mutex
	feed      *gtfs.FeedMessage
	fetchedAt time.Time
	err       error
}

// fetchBay511Regional returns the shared regional feed, making a real HTTP
// request only when the cached copy is older than bay511MinInterval — any
// number of callers within that window share the one fetch. The two sources
// that use this (Caltrain, ACE) are on the same 511.org account, so whichever
// one polls first each cycle supplies the key.
func fetchBay511Regional(apiKey string) (*gtfs.FeedMessage, error) {
	bay511Cache.mu.Lock()
	defer bay511Cache.mu.Unlock()

	if bay511Cache.feed != nil && time.Since(bay511Cache.fetchedAt) < bay511MinInterval {
		return bay511Cache.feed, bay511Cache.err
	}

	body, err := fetchGTFSRTBody(fmt.Sprintf(bay511RegionalURL, apiKey))
	bay511Cache.fetchedAt = time.Now()
	if err != nil {
		bay511Cache.err = err
		return nil, err
	}
	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		bay511Cache.err = err
		return nil, err
	}
	bay511Cache.feed = &feed
	bay511Cache.err = nil
	return &feed, nil
}

// bay511Vehicle is one vehicle position from the regional feed, reduced to
// what every per-agency Fetch needs, with its agency prefix already
// stripped from the trip id.
type bay511Vehicle struct {
	TripID       string
	VehicleLabel string
	Lat, Lon     float64
	BearingDeg   float64
	SpeedMps     float64
	Timestamp    int64
}

// bay511Vehicles returns only the vehicles belonging to one agency (matched
// by its route_id prefix, e.g. "CT" for Caltrain, "CE" for ACE).
func bay511Vehicles(feed *gtfs.FeedMessage, agencyPrefix string) []bay511Vehicle {
	prefix := agencyPrefix + ":"
	var out []bay511Vehicle
	for _, e := range feed.Entity {
		vp := e.GetVehicle()
		if vp == nil {
			continue
		}
		trip := vp.GetTrip()
		if !strings.HasPrefix(trip.GetRouteId(), prefix) {
			continue
		}
		pos := vp.GetPosition()
		if pos == nil {
			continue
		}
		out = append(out, bay511Vehicle{
			TripID:       strings.TrimPrefix(trip.GetTripId(), prefix),
			VehicleLabel: vp.GetVehicle().GetLabel(),
			Lat:          float64(pos.GetLatitude()),
			Lon:          float64(pos.GetLongitude()),
			BearingDeg:   float64(pos.GetBearing()),
			SpeedMps:     float64(pos.GetSpeed()),
			Timestamp:    int64(vp.GetTimestamp()),
		})
	}
	return out
}

// bay511LiveTrain builds a liveTrain from a regional-feed vehicle already
// matched to one of our own trains — shared by Caltrain and ACE so the
// position/speed/heading/timestamp conversion logic lives in one place.
func bay511LiveTrain(v bay511Vehicle, match dbTrain, trainNum string) liveTrain {
	lt := liveTrain{
		TrainNum:     trainNum,
		DisplayName:  match.DisplayName,
		TrainSlug:    match.Slug,
		CorridorName: match.CorridorName,
		CorridorSlug: match.CorridorSlug,
		Lat:          v.Lat,
		Lon:          v.Lon,
		Heading:      bearingToCompass(v.BearingDeg),
		Speed:        int(mpsToMph(v.SpeedMps) + 0.5),
	}
	if v.Timestamp > 0 {
		lt.LastUpdated = time.Unix(v.Timestamp, 0).UTC().Format(time.RFC3339)
	}
	return lt
}
