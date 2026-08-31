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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Shared plumbing for every /api/*-routes endpoint that draws agency route
// lines on the public map's "Route lines" layer (see templates/map.html).
// Amtrak's own endpoint (handlers_public.go) predates this and keeps its
// separate implementation since it already had its own caching in place.

const routeLineCacheTTL = 12 * time.Hour

// routeLineCache is a 12h in-memory cache for one agency's route-line
// GeoJSON, shared by MBTA/Caltrain/ACE/LIRR so each only needs to supply its
// own fetch function.
type routeLineCache struct {
	mu        sync.Mutex
	raw       []byte
	gz        []byte
	fetchedAt time.Time
}

func (c *routeLineCache) handle(w http.ResponseWriter, r *http.Request, fetch func() ([]byte, error)) {
	c.mu.Lock()
	if c.raw != nil && time.Since(c.fetchedAt) < routeLineCacheTTL {
		raw, gz := c.raw, c.gz
		c.mu.Unlock()
		writeRouteGeoJSON(w, r, raw, gz)
		return
	}
	c.mu.Unlock()

	data, err := fetch()
	if err != nil {
		http.Error(w, "Failed to fetch routes", 502)
		return
	}
	gz := gzipBytes(data)

	c.mu.Lock()
	c.raw, c.gz, c.fetchedAt = data, gz, time.Now()
	c.mu.Unlock()

	writeRouteGeoJSON(w, r, data, gz)
}

// writeRouteGeoJSON mirrors writeAmtrakRoutes (handlers_public.go): serve
// the gzip copy when the client accepts it, otherwise the raw JSON.
func writeRouteGeoJSON(w http.ResponseWriter, r *http.Request, raw, gz []byte) {
	w.Header().Set("Content-Type", "application/geo+json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if gz != nil && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Write(gz)
		return
	}
	w.Write(raw)
}

// emptyRouteGeoJSON is served (as a normal 200, not an error) when a source
// needs an API key that hasn't been configured yet — so the map doesn't show
// a load failure just because one agency isn't set up.
func emptyRouteGeoJSON() []byte {
	return []byte(`{"type":"FeatureCollection","features":[]}`)
}

// shapesToFeatureCollection turns parsed GTFS shapes into the GeoJSON shape
// templates/map.html expects. Every feature shares corridorName so they all
// map (via findCorridorSlug) to the one corridor regardless of which
// branch/line/direction the shape belongs to.
func shapesToFeatureCollection(shapesByID map[string][][2]float64, corridorName string) []byte {
	features := make([]map[string]interface{}, 0, len(shapesByID))
	for _, coords := range shapesByID {
		if len(coords) < 2 {
			continue
		}
		features = append(features, map[string]interface{}{
			"type":       "Feature",
			"properties": map[string]interface{}{"name": corridorName},
			"geometry": map[string]interface{}{
				"type":        "LineString",
				"coordinates": coords,
			},
		})
	}
	result := map[string]interface{}{"type": "FeatureCollection", "features": features}
	data, _ := json.Marshal(result)
	return data
}

// parseGTFSShapesZip extracts shapes.txt from a downloaded static GTFS zip
// and groups points into one polyline per shape_id, ordered by
// shape_pt_sequence, as [lon, lat] pairs (GeoJSON coordinate order).
func parseGTFSShapesZip(data []byte) (map[string][][2]float64, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening GTFS zip: %w", err)
	}
	var shapesFile *zip.File
	for _, f := range zr.File {
		if f.Name == "shapes.txt" {
			shapesFile = f
			break
		}
	}
	if shapesFile == nil {
		return nil, fmt.Errorf("GTFS zip has no shapes.txt")
	}
	rc, err := shapesFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	cr := csv.NewReader(rc)
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	col := make(map[string]int, len(header))
	for i, h := range header {
		col[strings.TrimSpace(h)] = i
	}
	idIdx, idOK := col["shape_id"]
	latIdx, latOK := col["shape_pt_lat"]
	lonIdx, lonOK := col["shape_pt_lon"]
	seqIdx, seqOK := col["shape_pt_sequence"]
	if !idOK || !latOK || !lonOK {
		return nil, fmt.Errorf("shapes.txt missing a required column")
	}

	type point struct {
		seq      int
		lat, lon float64
	}
	pointsByShape := map[string][]point{}

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if idIdx >= len(rec) || latIdx >= len(rec) || lonIdx >= len(rec) {
			continue
		}
		lat, err1 := strconv.ParseFloat(rec[latIdx], 64)
		lon, err2 := strconv.ParseFloat(rec[lonIdx], 64)
		if err1 != nil || err2 != nil {
			continue
		}
		seq := 0
		if seqOK && seqIdx < len(rec) {
			seq, _ = strconv.Atoi(rec[seqIdx])
		}
		id := rec[idIdx]
		pointsByShape[id] = append(pointsByShape[id], point{seq: seq, lat: lat, lon: lon})
	}

	out := make(map[string][][2]float64, len(pointsByShape))
	for id, pts := range pointsByShape {
		sort.Slice(pts, func(i, j int) bool { return pts[i].seq < pts[j].seq })
		coords := make([][2]float64, len(pts))
		for i, p := range pts {
			coords[i] = [2]float64{p.lon, p.lat}
		}
		out[id] = coords
	}
	return out, nil
}

// downloadBytes does a plain GET (following redirects, which the default
// client already does — confirmed against MTA's static feed, which 301s to
// an S3 bucket) and returns the body, capped at 64MB.
func downloadBytes(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AmazingTrak/1.0 (+https://foamer.online)")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 64<<20))
}

// fetch511StaticShapes downloads one 511.org operator's static GTFS zip
// (https://511.org/open-data — the same account/key already used for that
// operator's live GTFS-RT feed, see bay511.go) and returns its shapes.
func fetch511StaticShapes(apiKey, operatorID string) (map[string][][2]float64, error) {
	u := fmt.Sprintf("https://api.511.org/transit/datafeeds?api_key=%s&operator_id=%s",
		url.QueryEscape(apiKey), url.QueryEscape(operatorID))
	data, err := downloadBytes(u)
	if err != nil {
		return nil, err
	}
	return parseGTFSShapesZip(data)
}
