package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// SEPTA Regional Rail live positions, sourced from SEPTA's own public
// TrainView API (a plain JSON endpoint, not GTFS-RT) — fully open, no API
// key. Confirmed empirically (2026-08-31): unlike every other non-Amtrak
// source in this app, TrainView reports actual delay minutes per train
// ("late"), so this is the first non-Amtrak source that can legitimately set
// HasDelayInfo: true. Per SEPTA's own developer docs, a handful of
// low-ridership lines lack onboard GPS; entries with a zero/unparseable
// lat-lon are dropped rather than plotted at (0,0).
const septaTrainViewURL = "https://www3.septa.org/api/TrainView/index.php"

// flexString unmarshals a JSON string OR a bare JSON number into a Go string.
// Confirmed empirically (2026-08-31): TrainView is inconsistent about
// quoting its numeric fields (lat/lon/heading) — the same field arrived
// quoted in one poll and unquoted in the very next one, so a plain `string`
// field intermittently fails to unmarshal and drops the whole poll.
type flexString string

func (f *flexString) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = flexString(s)
		return nil
	}
	*f = flexString(b)
	return nil
}

type septaTrain struct {
	Lat      flexString `json:"lat"`
	Lon      flexString `json:"lon"`
	TrainNo  string     `json:"trainno"`
	NextStop string     `json:"nextstop"`
	Heading  flexString `json:"heading"`
	Late     int        `json:"late"`
}

type septaSource struct{}

func (septaSource) Key() string       { return "septa" }
func (septaSource) NeedsAPIKey() bool { return false }
func (septaSource) Description() string {
	return `Position and delay data come from SEPTA's own public <a href="https://www3.septa.org/api/" target="_blank" rel="noopener">TrainView API</a> (no API key needed). A handful of low-ridership lines lack onboard GPS and are skipped when that happens.`
}

func (septaSource) Fetch(app *App) ([]liveTrain, error) {
	index, err := loadDBTrainsByCorridorSlug(app, "septa")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active SEPTA trains in the database to match against")
	}

	req, err := http.NewRequest("GET", septaTrainViewURL, nil)
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
		return nil, fmt.Errorf("SEPTA TrainView HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	var trains []septaTrain
	if err := json.Unmarshal(body, &trains); err != nil {
		return nil, err
	}

	var out []liveTrain
	for _, t := range trains {
		// 999 is TrainView's own sentinel for "delay unknown" (confirmed
		// empirically, 2026-08-31 — these entries also carry empty
		// consist/TRACK fields), not a real ~16-hour delay; skip rather than
		// plot a train that isn't actually being tracked.
		if t.Late == 999 {
			continue
		}
		match, ok := index[t.TrainNo]
		if !ok {
			continue
		}
		lat, err1 := strconv.ParseFloat(string(t.Lat), 64)
		lon, err2 := strconv.ParseFloat(string(t.Lon), 64)
		if err1 != nil || err2 != nil || (lat == 0 && lon == 0) {
			continue
		}
		lt := liveTrain{
			TrainNum:     t.TrainNo,
			DisplayName:  match.DisplayName,
			TrainSlug:    match.Slug,
			CorridorName: match.CorridorName,
			CorridorSlug: match.CorridorSlug,
			Lat:          lat,
			Lon:          lon,
			NextStation:  t.NextStop,
			DelayMin:     t.Late,
			Status:       delayStatus(t.Late),
			HasDelayInfo: true,
		}
		if heading, err := strconv.ParseFloat(string(t.Heading), 64); err == nil {
			lt.Heading = bearingToCompass(heading)
		}
		out = append(out, lt)
	}
	return out, nil
}

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// SEPTA's public GTFS download is a zip-of-zips: gtfs_public.zip contains
// google_bus.zip and google_rail.zip separately (confirmed empirically,
// 2026-08-31 — 22MB combined, mostly bus). Only google_rail.zip's shapes are
// wanted; its routes.txt confirmed all 13 rows are route_type=2 (rail), so no
// further per-route filtering is needed once unwrapped.
const septaStaticGTFSURL = "https://www3.septa.org/developer/gtfs_public.zip"

var septaRouteLineCache routeLineCache

// extractZipEntry returns the raw bytes of one named file inside a zip.
func extractZipEntry(data []byte, name string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}
	for _, f := range zr.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("zip has no %s", name)
}

func fetchSEPTARouteGeoJSON() ([]byte, error) {
	outer, err := downloadBytes(septaStaticGTFSURL)
	if err != nil {
		return nil, err
	}
	railZip, err := extractZipEntry(outer, "google_rail.zip")
	if err != nil {
		return nil, err
	}
	shapes, err := parseGTFSShapesZip(railZip)
	if err != nil {
		return nil, err
	}
	return shapesToFeatureCollection(shapes, "SEPTA Regional Rail"), nil
}

func (app *App) handleSEPTARoutes(w http.ResponseWriter, r *http.Request) {
	septaRouteLineCache.handle(w, r, fetchSEPTARouteGeoJSON)
}
