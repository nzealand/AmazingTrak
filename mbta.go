package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// The map's line layer otherwise only draws from Amtrak's own NTAD dataset
// (see fetchAmtrakRoutes in handlers_public.go), which has no MBTA geometry
// at all. This fetches MBTA's own shapes so Commuter Rail lines render too.

const mbtaRouteIDsURL = "https://api-v3.mbta.com/routes?filter%5Btype%5D=2"
const mbtaShapesURLFmt = "https://api-v3.mbta.com/shapes?filter%%5Broute%%5D=%s"

var mbtaRouteLineCache routeLineCache

type mbtaIDResource struct {
	ID string `json:"id"`
}

type mbtaIDsResponse struct {
	Data []mbtaIDResource `json:"data"`
}

type mbtaShapeResource struct {
	ID         string `json:"id"`
	Attributes struct {
		Polyline string `json:"polyline"`
	} `json:"attributes"`
}

type mbtaShapesResponse struct {
	Data []mbtaShapeResource `json:"data"`
}

// decodePolyline decodes a Google encoded-polyline (precision 5) string into
// [lon, lat] pairs, GeoJSON coordinate order.
func decodePolyline(encoded string) [][2]float64 {
	var coords [][2]float64
	index, lat, lon := 0, 0, 0
	for index < len(encoded) {
		latDelta, ok := decodePolylineValue(encoded, &index)
		if !ok {
			break
		}
		lat += latDelta
		lonDelta, ok := decodePolylineValue(encoded, &index)
		if !ok {
			break
		}
		lon += lonDelta
		coords = append(coords, [2]float64{float64(lon) / 1e5, float64(lat) / 1e5})
	}
	return coords
}

func decodePolylineValue(encoded string, index *int) (int, bool) {
	var result, shift uint
	for {
		if *index >= len(encoded) {
			return 0, false
		}
		b := int(encoded[*index]) - 63
		*index++
		result |= uint(b&0x1f) << shift
		shift += 5
		if b < 0x20 {
			break
		}
	}
	if result&1 != 0 {
		return int(^(result >> 1)), true
	}
	return int(result >> 1), true
}

// fetchMBTARouteGeoJSON fetches every Commuter Rail route's shapes and
// returns them as a GeoJSON FeatureCollection matching the shape consumed by
// the map template (see extractLines/findCorridorSlug in templates/map.html).
// All features share the "MBTA Commuter Rail" name so they map to the one
// mbta-commuter-rail corridor regardless of which line they belong to.
func fetchMBTARouteGeoJSON() ([]byte, error) {
	client := &http.Client{Timeout: 45 * time.Second}

	var ids mbtaIDsResponse
	if err := mbtaGetJSON(client, mbtaRouteIDsURL, &ids); err != nil {
		return nil, fmt.Errorf("fetching MBTA commuter rail routes: %w", err)
	}
	if len(ids.Data) == 0 {
		return nil, fmt.Errorf("MBTA API returned no commuter rail routes")
	}
	routeIDs := make([]string, 0, len(ids.Data))
	for _, r := range ids.Data {
		routeIDs = append(routeIDs, r.ID)
	}

	shapesURL := fmt.Sprintf(mbtaShapesURLFmt, strings.Join(routeIDs, "%2C"))
	var shapes mbtaShapesResponse
	if err := mbtaGetJSON(client, shapesURL, &shapes); err != nil {
		return nil, fmt.Errorf("fetching MBTA commuter rail shapes: %w", err)
	}

	shapesByID := make(map[string][][2]float64, len(shapes.Data))
	for _, s := range shapes.Data {
		coords := decodePolyline(s.Attributes.Polyline)
		if len(coords) < 2 {
			continue
		}
		shapesByID[s.ID] = coords
	}
	return shapesToFeatureCollection(shapesByID, "MBTA Commuter Rail"), nil
}

func mbtaGetJSON(client *http.Client, url string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MBTA API HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

// handleMBTARoutes serves MBTA Commuter Rail line geometry, cached for 12h
// via the shared routeLineCache (routelines.go).
func (app *App) handleMBTARoutes(w http.ResponseWriter, r *http.Request) {
	mbtaRouteLineCache.handle(w, r, fetchMBTARouteGeoJSON)
}
