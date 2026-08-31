package main

import "database/sql"

// stationSeed is one physical station on a corridor's stop list.
type stationSeed struct {
	name      string
	lat, lon  float64
	sortOrder int
}

// septaStops is SEPTA Regional Rail's full station list (one row per
// physical station, not per line/platform), confirmed against SEPTA's
// public static GTFS (google_rail.zip within gtfs_public.zip, 2026-08-31).
// route_stops.txt gives each of the 13 lines' own stop order; a station
// shared by multiple lines (Suburban Station, Jefferson Station, 30th
// Street, Temple University, etc.) is kept once, at the position where it
// was first encountered walking the lines in the order Airport, Chestnut
// Hill East, Chestnut Hill West, Cynwyd, Fox Chase, Lansdale/Doylestown,
// Manayunk/Norristown, Media/Wawa, Paoli/Thorndale, Trenton, Warminster,
// West Trenton, Wilmington/Newark — the same order as the corridor
// description in seed.go/db.go. All 156 stations in stops.txt were covered
// by at least one line's direction_id=0 route_stops rows, so nothing is
// dropped or double-counted.
var septaStops = []stationSeed{
	{"Airport Terminals E & F", 39.87944, -75.23972, 10},
	{"Airport Terminals C & D", 39.87806, -75.24, 20},
	{"Airport Terminal B", 39.87722, -75.24139, 30},
	{"Airport Terminal A", 39.87611, -75.24528, 40},
	{"Eastwick", 39.89278, -75.24389, 50},
	{"Penn Medicine Station", 39.94806, -75.19028, 60},
	{"Gray 30th St Station", 39.95667, -75.18166, 70},
	{"Suburban Station", 39.95389, -75.16778, 80},
	{"Jefferson Station", 39.9525, -75.15806, 90},
	{"Temple University", 39.98139, -75.14944, 100},
	{"Wayne Junction", 40.02222, -75.16, 110},
	{"Wister", 40.03611, -75.16111, 120},
	{"Germantown", 40.0375, -75.17167, 130},
	{"Washington Ln", 40.05083, -75.17139, 140},
	{"Stenton", 40.06055, -75.17861, 150},
	{"Sedgwick", 40.06278, -75.18528, 160},
	{"Mount Airy", 40.06528, -75.19083, 170},
	{"Wyndmoor", 40.07333, -75.19666, 180},
	{"Gravers", 40.0775, -75.20167, 190},
	{"Chestnut Hill East", 40.08111, -75.20722, 200},
	{"Chestnut Hill West", 40.07639, -75.20834, 210},
	{"Highland", 40.07056, -75.21111, 220},
	{"St. Martins", 40.06583, -75.20444, 230},
	{"Richard Allen Ln", 40.0575, -75.19473, 240},
	{"Carpenter", 40.05111, -75.19139, 250},
	{"Upsal", 40.0425, -75.19, 260},
	{"Tulpehocken", 40.03528, -75.18694, 270},
	{"Chelten Av", 40.03, -75.18083, 280},
	{"Queen Ln", 40.02333, -75.17805, 290},
	{"North Philadelphia Septa", 39.99778, -75.15639, 300},
	{"Cynwyd", 40.00667, -75.23167, 310},
	{"Bala", 40.00111, -75.22778, 320},
	{"Wynnefield Av", 39.99, -75.22556, 330},
	{"Olney", 40.03333, -75.12278, 340},
	{"Lawndale", 40.05139, -75.10306, 350},
	{"Cheltenham", 40.05806, -75.09278, 360},
	{"Ryers", 40.06417, -75.08639, 370},
	{"Fox Chase", 40.07639, -75.08334, 380},
	{"North Broad", 39.99222, -75.15389, 390},
	{"Fern Rock Transit Center", 40.04055, -75.13472, 400},
	{"Melrose Park", 40.05944, -75.12917, 410},
	{"Elkins Park", 40.07139, -75.12778, 420},
	{"Jenkintown-Wyncote", 40.09278, -75.1375, 430},
	{"Glenside", 40.10139, -75.15361, 440},
	{"North Hills", 40.11195, -75.16944, 450},
	{"Oreland", 40.11833, -75.18389, 460},
	{"Fort Washington", 40.13583, -75.21222, 470},
	{"Ambler", 40.15361, -75.22472, 480},
	{"Penllyn", 40.17, -75.24416, 490},
	{"Gwynedd Valley", 40.18472, -75.25694, 500},
	{"North Wales", 40.21417, -75.27722, 510},
	{"Pennbrook", 40.23028, -75.28167, 520},
	{"Lansdale", 40.24278, -75.285, 530},
	{"9th St Lansdale", 40.25, -75.27917, 540},
	{"Fortuna", 40.25945, -75.26611, 550},
	{"Colmar", 40.26833, -75.25445, 560},
	{"Link Belt", 40.27389, -75.24667, 570},
	{"Chalfont", 40.28778, -75.20972, 580},
	{"New Britain", 40.2975, -75.17973, 590},
	{"Delaware Valley University", 40.29722, -75.16167, 600},
	{"Doylestown", 40.30639, -75.13028, 610},
	{"Wawa", 39.90068, -75.45856, 620},
	{"Elwyn", 39.9075, -75.41167, 630},
	{"Media", 39.91444, -75.395, 640},
	{"Moylan-Rose Valley", 39.90611, -75.38861, 650},
	{"Wallingford", 39.90361, -75.37194, 660},
	{"Swarthmore", 39.90222, -75.35083, 670},
	{"Morton", 39.90778, -75.32889, 680},
	{"Secane", 39.91583, -75.30972, 690},
	{"Primos", 39.92167, -75.29833, 700},
	{"Clifton-Aldan", 39.92667, -75.29028, 710},
	{"Gladstone", 39.93278, -75.28222, 720},
	{"Lansdowne", 39.9375, -75.27084, 730},
	{"Fernwood-Yeadon", 39.93972, -75.25584, 740},
	{"Angora", 39.94472, -75.23861, 750},
	{"49th St", 39.94361, -75.21667, 760},
	{"Allegheny", 40.00361, -75.16472, 770},
	{"East Falls", 40.01139, -75.19195, 780},
	{"Wissahickon", 40.01667, -75.21028, 790},
	{"Manayunk", 40.02694, -75.225, 800},
	{"Ivy Ridge", 40.03417, -75.23556, 810},
	{"Miquon", 40.05861, -75.26639, 820},
	{"Spring Mill", 40.07417, -75.28611, 830},
	{"Conshohocken", 40.07328, -75.310944, 840},
	{"Norristown Transit Center", 40.11278, -75.34417, 850},
	{"Main St", 40.11722, -75.34861, 860},
	{"Norristown Elm Street", 40.12083, -75.345, 870},
	{"Thorndale", 39.99278, -75.76361, 880},
	{"Downingtown", 40.00219, -75.71078, 890},
	{"Whitford", 40.01472, -75.63805, 900},
	{"Exton", 40.01929, -75.62171, 910},
	{"Malvern", 40.03639, -75.51556, 920},
	{"Paoli", 40.04276, -75.48376, 930},
	{"Daylesford", 40.04306, -75.46056, 940},
	{"Berwyn", 40.04805, -75.44222, 950},
	{"Devon", 40.04722, -75.42278, 960},
	{"Strafford", 40.04945, -75.40305, 970},
	{"Wayne", 40.04583, -75.38667, 980},
	{"St. Davids", 40.04389, -75.3725, 990},
	{"Radnor", 40.04472, -75.35889, 1000},
	{"Villanova", 40.03833, -75.34167, 1010},
	{"Rosemont", 40.02778, -75.32667, 1020},
	{"Bryn Mawr", 40.02195, -75.31639, 1030},
	{"Haverford", 40.01389, -75.29972, 1040},
	{"Ardmore", 40.00828, -75.2904, 1050},
	{"Wynnewood", 40.00278, -75.2725, 1060},
	{"Narberth", 40.00472, -75.26139, 1070},
	{"Merion", 39.99861, -75.25139, 1080},
	{"Overbrook", 39.98944, -75.24944, 1090},
	{"Trenton Transit Center", 40.21851, -74.75393, 1100},
	{"Levittown", 40.14028, -74.81695, 1110},
	{"Bristol", 40.10472, -74.85472, 1120},
	{"Croydon", 40.09361, -74.90667, 1130},
	{"Eddington", 40.08306, -74.93361, 1140},
	{"Cornwells Heights", 40.0709, -74.95432, 1150},
	{"Torresdale", 40.05444, -74.98444, 1160},
	{"Holmesburg Junction", 40.03278, -75.02361, 1170},
	{"Tacony", 40.02333, -75.03889, 1180},
	{"Bridesburg", 40.01056, -75.06973, 1190},
	{"North Philadelphia Amtrak", 39.99678, -75.15511, 1200},
	{"Ardsley", 40.11417, -75.15305, 1210},
	{"Roslyn", 40.12083, -75.13416, 1220},
	{"Crestmont", 40.13334, -75.11861, 1230},
	{"Willow Grove", 40.14389, -75.11417, 1240},
	{"Hatboro", 40.17611, -75.1025, 1250},
	{"Warminster", 40.19528, -75.08916, 1260},
	{"Newark DE", 39.66969, -75.75351, 1270},
	{"Churchman's Crossing", 39.695, -75.6725, 1280},
	{"Wilmington", 39.73726, -75.55109, 1290},
	{"Claymont Transit Center", 39.804167, -75.446111, 1300},
	{"Marcus Hook", 39.82167, -75.41944, 1310},
	{"Highland Av", 39.83361, -75.39333, 1320},
	{"Chester Transit Center", 39.84972, -75.36, 1330},
	{"Eddystone", 39.85722, -75.34222, 1340},
	{"Crum Lynne", 39.87194, -75.33111, 1350},
	{"Ridley Park", 39.88055, -75.32222, 1360},
	{"Prospect Park", 39.88833, -75.30889, 1370},
	{"Norwood", 39.89167, -75.30167, 1380},
	{"Glenolden", 39.89639, -75.29, 1390},
	{"Folcroft", 39.90055, -75.27972, 1400},
	{"Sharon Hill", 39.90445, -75.27084, 1410},
	{"Curtis Park", 39.90805, -75.265, 1420},
	{"Darby", 39.91306, -75.25445, 1430},
	{"Noble", 40.10444, -75.12417, 1440},
	{"Rydal", 40.1075, -75.11056, 1450},
	{"Meadowbrook", 40.11139, -75.0925, 1460},
	{"Bethayres", 40.11666, -75.06834, 1470},
	{"Philmont", 40.12194, -75.04361, 1480},
	{"Forest Hills", 40.12778, -75.02055, 1490},
	{"Somerton", 40.13055, -75.01195, 1500},
	{"Trevose", 40.14028, -74.9825, 1510},
	{"Neshaminy Falls", 40.14695, -74.96167, 1520},
	{"Langhorne", 40.16083, -74.9125, 1530},
	{"Woodbourne", 40.1925, -74.88917, 1540},
	{"Yardley", 40.23528, -74.83056, 1550},
	{"West Trenton", 40.25778, -74.81528, 1560},
}

// brightlineStops is Brightline's full Florida-corridor station list,
// confirmed against its public static GTFS (bl_gtfs.zip, 2026-08-31)
// stops.txt — 6 stations total, no duplicates to resolve. Ordered
// geographically south to north (Miami to West Palm Beach), then to
// Orlando, matching the actual route.
var brightlineStops = []stationSeed{
	{"Miami", 25.779968, -80.195641, 10},
	{"Aventura", 25.95846, -80.14728, 20},
	{"Fort Lauderdale", 26.123619, -80.145546, 30},
	{"Boca Raton", 26.353828, -80.087494, 40},
	{"West Palm Beach", 26.712015, -80.055359, 50},
	{"Orlando", 28.411698, -81.307619, 60},
}

// sqlExecer is satisfied by both *sql.DB (the versioned migration in db.go,
// run against an already-live database) and *sql.Tx (seedDB in seed.go, run
// inside its one seeding transaction), so the row data and insert logic
// only need to live once and are shared by both callers.
type sqlExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// seedCorridorStops inserts one corridor's station list, looked up by slug
// rather than a hardcoded corridor_id. Unlike the Amtrak corridors in
// data.sql — whose ids are fixed by corridorSeeds' insertion order on a
// fresh install — a migration running against an already-seeded production
// database can't assume what id a newly-inserted corridor got, so every
// insert here goes through a `WHERE corridors.slug=?` lookup instead.
func seedCorridorStops(ex sqlExecer, corridorSlug string, stops []stationSeed) error {
	for _, s := range stops {
		if _, err := ex.Exec(
			`INSERT INTO stops (corridor_id, name, latitude, longitude, sort_order)
				SELECT id, ?, ?, ?, ? FROM corridors WHERE slug=?`,
			s.name, s.lat, s.lon, s.sortOrder, corridorSlug,
		); err != nil {
			return err
		}
	}
	return nil
}
