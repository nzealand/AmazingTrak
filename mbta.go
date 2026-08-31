package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MBTA Commuter Rail live positions, sourced from the official MBTA V3 API
// (https://api-v3.mbta.com), a JSON:API — not GTFS-RT protobuf, so this uses
// only encoding/json, no extra dependency. filter[route_type]=2 restricts
// results to Commuter Rail; MBTA's subway/Green Line/bus/ferry are never
// requested by this source.
//
// Confirmed empirically: a /vehicles response's own attributes.label is the
// physical railcar/locomotive number, NOT the public train number — the
// train number is the sideloaded trip resource's attributes.name (hence
// include=trip below).
//
// An API key is optional (20 req/min unauthenticated, 1000/min with a free
// key from api-v3.mbta.com) — one request per poll is nowhere near either
// limit, so NeedsAPIKey is false; a key can still be entered for headroom.
const mbtaVehiclesURL = "https://api-v3.mbta.com/vehicles?filter%5Broute_type%5D=2&include=trip"

type mbtaSource struct{}

func (mbtaSource) Key() string       { return "mbta" }
func (mbtaSource) NeedsAPIKey() bool { return false }
func (mbtaSource) Description() string {
	return `Position data comes from the official <a href="https://api-v3.mbta.com" target="_blank" rel="noopener">MBTA V3 API</a>, filtered to Commuter Rail only (no subway/Green Line/bus/ferry). No API key needed — the unauthenticated rate limit (20 req/min) is far more than this feature uses.`
}

type mbtaResource struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    map[string]interface{} `json:"attributes"`
	Relationships struct {
		Trip struct {
			Data *struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"trip"`
	} `json:"relationships"`
}

type mbtaResponse struct {
	Data     []mbtaResource `json:"data"`
	Included []mbtaResource `json:"included"`
}

func (mbtaSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "mbta")
	if err != nil {
		return nil, err
	}

	index, err := loadDBTrainsByCorridorSlug(app, "mbta-commuter-rail")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active MBTA Commuter Rail trains in the database to match against")
	}

	req, err := http.NewRequest("GET", mbtaVehiclesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")
	if src.APIKey != "" {
		req.Header.Set("x-api-key", src.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MBTA API HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	var feed mbtaResponse
	if err := json.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	// Trip resources are sideloaded in "included"; attributes.name is the
	// public train number (attributes.label on the vehicle itself is the
	// physical railcar, not the train number).
	tripTrainNum := map[string]string{}
	for _, inc := range feed.Included {
		if inc.Type != "trip" {
			continue
		}
		if name, ok := inc.Attributes["name"].(string); ok && name != "" {
			tripTrainNum[inc.ID] = name
		}
	}

	var out []liveTrain
	for _, v := range feed.Data {
		if v.Relationships.Trip.Data == nil {
			continue
		}
		trainNum := tripTrainNum[v.Relationships.Trip.Data.ID]
		if trainNum == "" {
			continue
		}
		match, ok := index[trainNum]
		if !ok {
			continue
		}

		lat, _ := v.Attributes["latitude"].(float64)
		lon, _ := v.Attributes["longitude"].(float64)
		if lat == 0 && lon == 0 {
			continue
		}
		lt := liveTrain{
			TrainNum:     trainNum,
			DisplayName:  match.DisplayName,
			TrainSlug:    match.Slug,
			CorridorName: match.CorridorName,
			CorridorSlug: match.CorridorSlug,
			Lat:          lat,
			Lon:          lon,
		}
		if bearing, ok := v.Attributes["bearing"].(float64); ok {
			lt.Heading = bearingToCompass(bearing)
		}
		if speedMps, ok := v.Attributes["speed"].(float64); ok {
			lt.Speed = int(mpsToMph(speedMps) + 0.5)
		}
		if updatedAt, ok := v.Attributes["updated_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				lt.LastUpdated = t.UTC().Format(time.RFC3339)
			}
		}
		out = append(out, lt)
	}
	return out, nil
}
