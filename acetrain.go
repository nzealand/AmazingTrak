package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// ACE (Altamont Corridor Express) live positions, sourced from the same
// 511.org GTFS-RT infrastructure as Caltrain — agency code "CE" (confirmed
// against 511.org's /transit/gtfsoperators list; the same free 511.org API
// key used for Caltrain works here too). ACE runs only ~8-10 weekday
// peak-commute trains, so this feed is legitimately empty most of the day —
// zero entities is expected, not a bug.
//
// Unlike Caltrain, ACE's GTFS-RT trip_id is not the bare train number: static
// GTFS trips.txt shows ids like "ACE01", "ACE03H", "ACE01fifa2" (holiday and
// special-event schedule variants of the same physical train), confirmed
// against a downloaded copy of ACE's static feed. The digits right after the
// "ACE" prefix are the train number regardless of any trailing variant
// suffix, so aceTripIDPattern strips exactly that.
const aceVehiclePositionsURL = "https://api.511.org/transit/vehiclepositions?api_key=%s&agency=CE"

var aceTripIDPattern = regexp.MustCompile(`^ACE(\d+)`)

type aceSource struct{}

func (aceSource) Key() string       { return "ace" }
func (aceSource) NeedsAPIKey() bool { return true }
func (aceSource) Description() string {
	return "Position data comes from the 511.org SF Bay Open Data GTFS-Realtime feed (same account as Caltrain — the same API key works for both). ACE runs only a handful of weekday peak-commute trains, so it's normal to see nothing outside those windows."
}

func (aceSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "ace")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" {
		return nil, fmt.Errorf("no 511.org API key configured")
	}

	index, err := loadDBTrainsByCorridorSlug(app, "ace")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active ACE trains in the database to match against")
	}

	body, err := fetchGTFSRTBody(fmt.Sprintf(aceVehiclePositionsURL, src.APIKey))
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

		trainNum := ""
		if m := aceTripIDPattern.FindStringSubmatch(vp.GetTrip().GetTripId()); m != nil {
			trainNum = m[1]
		}
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

// fetchGTFSRTBody does a plain authenticated GET and returns the raw body,
// shared by every GTFS-RT-over-511.org source (Caltrain, ACE).
func fetchGTFSRTBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}
