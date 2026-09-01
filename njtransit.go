package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NJ Transit rail live positions, sourced from NJT's own RailData API
// (https://raildata.njtransit.com/api/TrainData, confirmed 2026-09-01).
// Unlike every other source in this app, this one needs a username+password
// pair (live_sources.api_key/api_secret — see the credentialPairSource
// interface below) rather than a single static key: NJT's API exchanges
// credentials for a short-lived session token via getToken, and every other
// call takes that token as a form field. Confirmed empirically: NJT caps
// token generation at 5/day, so the token is persisted in
// live_sources.cached_token (see updateLiveSourceCachedToken) and reused
// across polls and process restarts — a fresh token is only requested when
// there's no cached one yet, or an actual call comes back "Invalid token."
//
// NJT is modeled as a single "nj-transit" corridor (not split per line like
// Metra): a live poll across 8 major hub stations plus current vehicle
// positions (2026-09-01, 122 distinct real train numbers) found zero
// cross-line collisions, unlike Metra's confirmed collisions — see
// stations_njtransit.go.
//
// Confirmed empirically against real API responses: GetVehicleData's ID
// field carries a letter prefix for anything that isn't actually an NJT
// train — "A" (an Amtrak train sharing NJT trackage, e.g. Northeast
// Corridor), "S" (a SEPTA train on the shared Atlantic City Line trackage),
// "X" (non-revenue, no passengers) — confirmed against NJT's own
// TrainIdPrefixes documentation and live samples ("A2118" on the Acela,
// "S9750" on Septa). Both Amtrak and SEPTA are already tracked as their own
// sources in this app, so anything not purely numeric is skipped here.
var njtNumericIDPattern = regexp.MustCompile(`^\d+$`)

// njTransitTrainNumbers is a starter roster of real, currently-scheduled
// train numbers — unlike Metra, NJT's public static GTFS has no
// rider-facing train-number field at all (trip_id is just a sequential
// integer), so this couldn't be derived from the static file the way every
// other agency's roster was. Instead, gathered 2026-09-01 via one
// getVehicleData poll (65 active trains) plus getTrainSchedule19Records at
// 8 major hub stations (NY, NP, HB, TR, AC, GL, BH, DO — chosen to cover
// all 12 lines), deduplicated, with any Amtrak/SEPTA-prefixed ID excluded
// (see the filtering note above). Expect to add more via the admin panel as
// other real train numbers are observed running.
var njTransitTrainNumbers = []string{
	"51", "53", "55", "57", "59", "61", "275", "339", "427", "429", "432", "433", "435", "438",
	"439", "440", "442", "643", "645", "657", "660", "674", "682", "807", "858", "880", "881",
	"882", "884", "1009", "1011", "1055", "1061", "1079", "1085", "1087", "1124", "1126", "1165",
	"1167", "1169", "1171", "1173", "1221", "1223", "1225", "1269", "1271", "1274", "1625",
	"1627", "1629", "1631", "1633", "1635", "1637", "3127", "3265", "3266", "3267", "3270",
	"3271", "3275", "3361", "3363", "3373", "3510", "3511", "3513", "3515", "3598", "3725",
	"3736", "3738", "3862", "3864", "3866", "3867", "3868", "3870", "3871", "3872", "3873",
	"3874", "3876", "3949", "3951", "3953", "3957", "3961", "4355", "4372", "4378", "4384",
	"4392", "4398", "4632", "4634", "4638", "4642", "5435", "5441", "5733", "5746", "5939",
	"6252", "6273", "6279", "6355", "6431", "6437", "6643", "6647", "6655", "6662", "6664",
	"6666", "6668", "6670", "6672", "6674", "6676",
}

const njtAPIBase = "https://raildata.njtransit.com/api/TrainData"

// njtLastModifiedFormat parses NJT's own timestamp strings, e.g.
// "01-Sep-2026 05:43:53 PM" (confirmed against a real getVehicleData response).
const njtLastModifiedFormat = "02-Jan-2006 03:04:05 PM"

type njtSource struct{}

func (njtSource) Key() string       { return "njtransit" }
func (njtSource) NeedsAPIKey() bool { return true }

// NeedsCredentialPair marks this source as needing a username+password pair
// (see credentialPairSource) instead of a single API key.
func (njtSource) NeedsCredentialPair() bool { return true }

func (njtSource) Description() string {
	return `Position data comes from NJ Transit's own <a href="https://developer.njtransit.com/" target="_blank" rel="noopener">RailData API</a>. Requires a free username/password from NJT's developer portal registration. Excludes Amtrak and SEPTA trains that share NJT trackage — those are tracked under their own sources above.`
}

// credentialPairSource is an optional interface (see liveSource) implemented
// only by a source needing a username+password pair instead of one static
// API key. Checked via a type assertion in handlers_admin.go so every other
// (single-key) source doesn't need a no-op method.
type credentialPairSource interface {
	NeedsCredentialPair() bool
}

func (njtSource) Fetch(app *App) ([]liveTrain, error) {
	src, err := getLiveSource(app.db, "njtransit")
	if err != nil {
		return nil, err
	}
	if src.APIKey == "" || src.APISecret == "" {
		return nil, fmt.Errorf("no NJ Transit username/password configured")
	}

	index, err := loadDBTrainsByCorridorSlug(app, "nj-transit")
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, fmt.Errorf("no active NJ Transit trains in the database to match against")
	}

	token := src.CachedToken
	if token == "" {
		token, err = njtGetToken(src.APIKey, src.APISecret)
		if err != nil {
			return nil, fmt.Errorf("getting NJ Transit token: %w", err)
		}
		if err := updateLiveSourceCachedToken(app.db, "njtransit", token); err != nil {
			return nil, err
		}
	}

	vehicles, err := njtGetVehicleData(token)
	if err != nil && isNJTInvalidToken(err) {
		// Cached token expired or was revoked — get a fresh one and retry
		// exactly once, mirroring the reference client's own retry pattern.
		token, err = njtGetToken(src.APIKey, src.APISecret)
		if err != nil {
			return nil, fmt.Errorf("refreshing NJ Transit token: %w", err)
		}
		if err := updateLiveSourceCachedToken(app.db, "njtransit", token); err != nil {
			return nil, err
		}
		vehicles, err = njtGetVehicleData(token)
	}
	if err != nil {
		return nil, err
	}

	var out []liveTrain
	for _, v := range vehicles {
		trainNum := strings.TrimSpace(v.ID)
		if !njtNumericIDPattern.MatchString(trainNum) {
			continue
		}
		match, ok := index[trainNum]
		if !ok {
			continue
		}
		lat, err1 := strconv.ParseFloat(v.LATITUDE, 64)
		lon, err2 := strconv.ParseFloat(v.LONGITUDE, 64)
		if err1 != nil || err2 != nil || (lat == 0 && lon == 0) {
			continue
		}
		secLate, _ := strconv.Atoi(v.SEC_LATE)
		delayMin := secLate / 60
		lt := liveTrain{
			TrainNum:     trainNum,
			DisplayName:  match.DisplayName,
			TrainSlug:    match.Slug,
			CorridorName: match.CorridorName,
			CorridorSlug: match.CorridorSlug,
			Lat:          lat,
			Lon:          lon,
			NextStation:  v.NEXT_STOP,
			DelayMin:     delayMin,
			Status:       delayStatus(delayMin),
			HasDelayInfo: true,
		}
		if t, err := time.Parse(njtLastModifiedFormat, v.LAST_MODIFIED); err == nil {
			lt.LastUpdated = t.UTC().Format(time.RFC3339)
		}
		out = append(out, lt)
	}
	return out, nil
}

// njtInvalidTokenError is returned when NJT's API rejects the current
// token, so Fetch knows to get a fresh one and retry instead of giving up.
type njtInvalidTokenError struct{ msg string }

func (e *njtInvalidTokenError) Error() string { return e.msg }

func isNJTInvalidToken(err error) bool {
	_, ok := err.(*njtInvalidTokenError)
	return ok
}

// njtGetToken exchanges a username/password for a session token. NJT caps
// this at 5 calls/day — callers must cache the result (see Fetch above) and
// only call this again after an actual "Invalid token." response.
func njtGetToken(username, password string) (string, error) {
	var resp struct {
		Authenticated string `json:"Authenticated"`
		UserToken     string `json:"UserToken"`
	}
	if err := njtPost("getToken", map[string]string{"username": username, "password": password}, &resp); err != nil {
		return "", err
	}
	if resp.Authenticated != "True" {
		return "", fmt.Errorf("NJ Transit rejected the configured username/password")
	}
	return resp.UserToken, nil
}

// njtVehicle is one entry of a getVehicleData response — confirmed
// empirically 2026-09-01 against a real payload.
type njtVehicle struct {
	ID            string `json:"ID"`
	TRAIN_LINE    string `json:"TRAIN_LINE"`
	LAST_MODIFIED string `json:"LAST_MODIFIED"`
	SEC_LATE      string `json:"SEC_LATE"`
	NEXT_STOP     string `json:"NEXT_STOP"`
	LONGITUDE     string `json:"LONGITUDE"`
	LATITUDE      string `json:"LATITUDE"`
}

// njtGetVehicleData returns the position/status of every currently-active
// NJT train. Confirmed empirically: a train appears here only if it has
// moved in the last 5 minutes, so — unlike LIRR/Metra — no separate
// entity-age staleness filter is needed on top of this.
func njtGetVehicleData(token string) ([]njtVehicle, error) {
	var out []njtVehicle
	if err := njtPost("getVehicleData", map[string]string{"token": token}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// njtPost issues one RailData API call: a multipart/form-data POST (the API
// does not accept a JSON body, confirmed empirically) to {njtAPIBase}/{method},
// decoding the JSON response into out. NJT signals an invalid/expired token
// with a non-2xx status and a body of {"errorMessage":"Invalid token."}.
func njtPost(method string, fields map[string]string, out interface{}) error {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		fw, err := w.CreateFormField(k)
		if err != nil {
			return err
		}
		if _, err := fw.Write([]byte(v)); err != nil {
			return err
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", njtAPIBase+"/"+method, body)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "text/plain")
	req.Header.Set("content-type", w.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Message string `json:"errorMessage"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Message == "Invalid token." {
			return &njtInvalidTokenError{msg: errResp.Message}
		}
		return fmt.Errorf("NJ Transit API HTTP %d for %s", resp.StatusCode, method)
	}
	if len(respBody) == 0 {
		return fmt.Errorf("NJ Transit API returned an empty response for %s", method)
	}
	return json.Unmarshal(respBody, out)
}

// ---- Route line geometry (for the public map's "Route lines" layer) ----
//
// NJT's static schedule.zip is public with no login needed (confirmed
// 2026-09-01: https://www.njtransit.com/rail_data.zip redirects to a public
// CDN URL) — only the live API above requires credentials. The zip's
// shapes.txt mixes heavy-rail (route_type 2) and light-rail (route_type 0:
// Hudson-Bergen, Newark, River Line) geometry together with no way to tell
// them apart directly, so route_type is looked up per shape via trips.txt
// (parseGTFSTripShapeRoutes, shared with metra.go) and routes.txt, and only
// route_type-2 shapes are kept — light rail is out of scope for this
// corridor, same as every other agency in this app.
const njtStaticGTFSURL = "https://www.njtransit.com/rail_data.zip"

var njtRouteLineCache routeLineCache

// parseGTFSRouteTypes extracts routes.txt from a downloaded static GTFS zip
// and returns a route_id -> route_type map.
func parseGTFSRouteTypes(data []byte) (map[string]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening GTFS zip: %w", err)
	}
	var routesFile *zip.File
	for _, f := range zr.File {
		if f.Name == "routes.txt" {
			routesFile = f
			break
		}
	}
	if routesFile == nil {
		return nil, fmt.Errorf("GTFS zip has no routes.txt")
	}
	rc, err := routesFile.Open()
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
	idIdx, idOK := col["route_id"]
	typeIdx, typeOK := col["route_type"]
	if !idOK || !typeOK {
		return nil, fmt.Errorf("routes.txt missing route_id or route_type column")
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
		if idIdx >= len(rec) || typeIdx >= len(rec) {
			continue
		}
		out[strings.TrimSpace(rec[idIdx])] = strings.TrimSpace(rec[typeIdx])
	}
	return out, nil
}

func fetchNJTransitRouteGeoJSON() ([]byte, error) {
	data, err := downloadBytes(njtStaticGTFSURL)
	if err != nil {
		return nil, err
	}
	shapes, err := parseGTFSShapesZip(data)
	if err != nil {
		return nil, err
	}
	shapeRoutes, err := parseGTFSTripShapeRoutes(data)
	if err != nil {
		return nil, err
	}
	routeTypes, err := parseGTFSRouteTypes(data)
	if err != nil {
		return nil, err
	}

	heavyRail := map[string][][2]float64{}
	for shapeID, coords := range shapes {
		routeID, ok := shapeRoutes[shapeID]
		if !ok || routeTypes[routeID] != "2" {
			continue
		}
		heavyRail[shapeID] = coords
	}
	return shapesToFeatureCollection(heavyRail, "NJ Transit"), nil
}

func (app *App) handleNJTransitRoutes(w http.ResponseWriter, r *http.Request) {
	njtRouteLineCache.handle(w, r, fetchNJTransitRouteGeoJSON)
}
