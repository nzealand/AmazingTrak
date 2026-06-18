package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//go:embed data.sql
var stopsSQL string

//go:embed capitol_corridor_stops.sql
var capitolCorridorStopsSQL string

type corridorSeed struct {
	name, slug, region, description string
	sort                            int
}

// corridorSeeds must be ordered so that auto-increment IDs match data.sql references (1..44).
// Corridors 1–5: NEC (pure Amtrak). 6–30: State-Supported. 31–44: Long Distance (pure Amtrak).
// 45–46: Seasonal.
var corridorSeeds = []corridorSeed{
	// 1 NEC
	{"Amtrak Acela", "amtrak-acela", "Northeast Corridor",
		"Amtrak's premium high-speed service connecting Boston, New York, Philadelphia, and Washington D.C. The fastest train in the Americas, reaching speeds up to 150 mph along the Northeast Corridor.", 1},
	// 2
	{"Amtrak Northeast Regional", "amtrak-northeast-regional", "Northeast Corridor",
		"Amtrak's busiest long-distance service, running the full length of the Northeast Corridor from Boston and New York south through Philadelphia to Washington D.C., with extensions into Virginia.", 2},
	// 3
	{"Amtrak NEC — Richmond/Newport News/Norfolk", "amtrak-nec-richmond", "Mid-Atlantic",
		"Northeast Regional extensions south from Washington D.C. into Virginia, serving Richmond, Newport News, Williamsburg, and Norfolk.", 3},
	// 4
	{"Amtrak NEC — Roanoke", "amtrak-nec-roanoke", "Mid-Atlantic",
		"Northeast Regional service extending from Washington D.C. through Charlottesville to Roanoke, Virginia.", 4},
	// 5
	{"Amtrak NEC — Springfield", "amtrak-nec-springfield", "New England",
		"Hartford Line and Valley Flyer service connecting New Haven, Hartford, and Springfield.", 5},
	// 6 State-Supported
	{"Adirondack", "adirondack", "Northeast",
		"Daily service connecting New York City and Montreal, Quebec, running along the Hudson River and through the Adirondack region.", 6},
	// 7
	{"Blue Water", "blue-water", "Midwest",
		"Michigan corridor service connecting Port Huron and Chicago via Flint, East Lansing, and Kalamazoo.", 7},
	// 8
	{"Borealis", "borealis", "Midwest",
		"Daily service between Chicago and St. Paul-Minneapolis, connecting Wisconsin communities along the route.", 8},
	// 9
	{"Capitol Corridor", "capitol-corridor", "California",
		"Amtrak's busiest non-NEC corridor, connecting the San Francisco Bay Area to Sacramento through the heart of Northern California with up to seven daily round trips.", 9},
	// 10
	{"Carl Sandburg / Illinois Zephyr", "carl-sandburg-illinois-zephyr", "Midwest",
		"Daily trains between Chicago and Quincy, Illinois, named for the poet Carl Sandburg and the historic Zephyr service.", 10},
	// 11
	{"Carolinian", "carolinian", "Southeast",
		"Daily service connecting New York City and Charlotte, North Carolina, sharing the NEC south to Washington before heading into the Carolinas.", 11},
	// 12
	{"Amtrak Cascades", "amtrak-cascades", "Pacific Northwest",
		"Washington and Oregon state-supported service connecting Vancouver, B.C. through Seattle and Portland to Eugene, Oregon.", 12},
	// 13
	{"Downeaster", "downeaster", "New England",
		"Maine DOT-supported service between Boston's North Station and Brunswick, Maine, with frequent daily departures.", 13},
	// 14
	{"Ethan Allen Express", "ethan-allen-express", "Northeast",
		"Daily service between New York City and Burlington, Vermont, named for the Revolutionary War hero.", 14},
	// 15
	{"Gold Runner", "gold-runner", "California",
		"California's San Joaquin Valley corridor connecting the San Francisco Bay Area (Oakland/Emeryville) to Bakersfield. Formerly known as the San Joaquins.", 15},
	// 16
	{"Heartland Flyer", "heartland-flyer", "South Central",
		"Daily service between Fort Worth, Texas and Oklahoma City, Oklahoma, the only Amtrak route serving Oklahoma.", 16},
	// 17
	{"Hiawatha", "hiawatha", "Midwest",
		"Frequent corridor service between Chicago and Milwaukee, Wisconsin, with multiple daily round trips.", 17},
	// 18
	{"Illini / Saluki", "illini-saluki", "Midwest",
		"Daily service between Chicago and Carbondale, Illinois, serving Champaign-Urbana and other downstate communities.", 18},
	// 19
	{"Keystone", "keystone", "Mid-Atlantic",
		"Pennsylvania DOT-supported frequent service between New York Penn Station and Harrisburg, PA, with stops along the Main Line.", 19},
	// 20
	{"Lincoln / Missouri River Runner", "lincoln-missouri-river-runner", "Midwest",
		"Combined corridor covering Lincoln Service (Chicago to St. Louis) and Missouri River Runner (St. Louis to Kansas City).", 20},
	// 21
	{"Maple Leaf", "maple-leaf", "Northeast",
		"Joint Amtrak/VIA Rail Canada service connecting New York City and Toronto, Ontario, passing through Niagara Falls.", 21},
	// 22
	{"Mardi Gras Service", "mardi-gras-service", "Gulf Coast",
		"Restored Gulf Coast service connecting New Orleans, Louisiana and Mobile, Alabama.", 22},
	// 23
	{"Empire Service — New York to Albany", "empire-service-albany", "Northeast",
		"Sub-segment of Empire Service providing frequent departures between New York Penn Station and Albany-Rensselaer.", 23},
	// 24
	{"Empire Service", "empire-service", "Northeast",
		"New York state-supported service connecting New York City to Albany, Utica, Syracuse, Rochester, Buffalo, and Niagara Falls.", 24},
	// 25
	{"Pacific Surfliner", "pacific-surfliner", "California",
		"California's coastal corridor connecting San Luis Obispo, Santa Barbara, Los Angeles, and San Diego with frequent daily departures.", 25},
	// 26
	{"Pennsylvanian", "pennsylvanian", "Mid-Atlantic",
		"Daily service between New York City and Pittsburgh, Pennsylvania, traversing the Allegheny Mountains via Harrisburg.", 26},
	// 27
	{"Pere Marquette", "pere-marquette", "Midwest",
		"Michigan DOT-supported daily service between Chicago and Grand Rapids, Michigan.", 27},
	// 28
	{"Piedmont", "piedmont", "Southeast",
		"North Carolina-supported multiple daily round trips between Raleigh and Charlotte, one of the highest-frequency state-supported corridors.", 28},
	// 29
	{"Vermonter", "vermonter", "New England",
		"Daily service between St. Albans, Vermont and Washington D.C., running the length of Vermont through New England and the Northeast Corridor.", 29},
	// 30
	{"Wolverine", "wolverine", "Midwest",
		"Michigan DOT-supported service connecting Pontiac and Detroit through Ann Arbor, Kalamazoo, and Niles to Chicago.", 30},
	// 31 Long Distance
	{"Amtrak Auto Train", "amtrak-auto-train", "Southeast",
		"The world's longest passenger train, running daily between Lorton, Virginia (near Washington D.C.) and Sanford, Florida, carrying passengers and their vehicles overnight.", 31},
	// 32
	{"Amtrak California Zephyr", "amtrak-california-zephyr", "Transcontinental",
		"Often called America's most scenic train, the California Zephyr runs daily between Chicago and Emeryville (San Francisco Bay Area) through Denver, the Rockies, and the Sierra Nevada.", 32},
	// 33
	{"Amtrak Cardinal", "amtrak-cardinal", "Midwest",
		"Tri-weekly service connecting New York City and Chicago through Washington D.C., the New River Gorge, Cincinnati, and Indianapolis.", 33},
	// 34
	{"Amtrak City of New Orleans", "amtrak-city-of-new-orleans", "South Central",
		"Daily overnight service between Chicago and New Orleans, immortalized in the song by Steve Goodman. Serves Champaign, Memphis, and Jackson.", 34},
	// 35
	{"Amtrak Coast Starlight", "amtrak-coast-starlight", "Pacific Coast",
		"Daily service between Seattle and Los Angeles along the Pacific Coast, passing through Portland, Sacramento, the Bay Area, and Santa Barbara.", 35},
	// 36
	{"Amtrak Crescent", "amtrak-crescent", "Southeast",
		"Daily overnight service between New York City and New Orleans through Washington D.C., Charlotte, Atlanta, and Birmingham.", 36},
	// 37
	{"Amtrak Empire Builder", "amtrak-empire-builder", "Northern Transcontinental",
		"Daily service between Chicago and both Seattle and Portland, splitting at Spokane. Passes through Milwaukee, Minneapolis, Glacier National Park, and the Columbia River Gorge.", 37},
	// 38
	{"Amtrak Floridian", "amtrak-floridian", "Southeast",
		"Long-distance service connecting Chicago and Miami/Tampa via Washington D.C. and Jacksonville, restored to provide a new connection to Florida.", 38},
	// 39
	{"Amtrak Lake Shore Limited", "amtrak-lake-shore-limited", "Northeast/Midwest",
		"Daily service between Chicago and both New York City and Boston, splitting at Albany. Runs along the Great Lakes through Cleveland, Erie, and Buffalo.", 39},
	// 40
	{"Amtrak Palmetto", "amtrak-palmetto", "Southeast",
		"Daily service between New York City and Savannah, Georgia, sharing the NEC south before heading into the Carolinas.", 40},
	// 41
	{"Amtrak Silver Meteor", "amtrak-silver-meteor", "Southeast",
		"Daily overnight service between New York City and Miami, serving the Eastern Seaboard through the Carolinas and Florida.", 41},
	// 42
	{"Amtrak Southwest Chief", "amtrak-southwest-chief", "Southern Transcontinental",
		"Daily service between Chicago and Los Angeles through Kansas City, Albuquerque, and Flagstaff. Follows the historic Santa Fe Trail and Route 66 corridor.", 42},
	// 43
	{"Amtrak Sunset Limited", "amtrak-sunset-limited", "Southern Transcontinental",
		"Tri-weekly service between New Orleans and Los Angeles — America's oldest named train — crossing Louisiana, Texas, New Mexico, and Arizona.", 43},
	// 44
	{"Amtrak Texas Eagle", "amtrak-texas-eagle", "South Central",
		"Daily service between Chicago and San Antonio, Texas through St. Louis, Little Rock, and Dallas. Connects at San Antonio with the Sunset Limited.", 44},
	// 45–46 Seasonal
	{"Berkshire Flyer", "berkshire-flyer", "Northeast",
		"Seasonal Friday/Sunday service between New York Penn Station and Pittsfield, Massachusetts, serving the Berkshire region.", 45},
	{"Winter Park Express", "winter-park-express", "Rocky Mountains",
		"Seasonal ski-season service between Denver Union Station and Winter Park/Fraser ski resort in Colorado.", 46},
}

// trainSeeds maps corridor index (0-based) → train numbers.
// Empty slice means no trains seeded for that corridor.
var trainSeeds = [][]string{
	// 1 Acela
	{"816", "817", "880", "2102", "2103", "2104", "2108", "2109", "2110", "2113", "2115",
		"2121", "2122", "2123", "2124", "2126", "2130", "2150", "2151", "2152", "2153",
		"2154", "2155", "2159", "2162", "2163", "2166", "2167", "2168", "2169", "2170",
		"2171", "2172", "2173", "2174", "2190", "2192", "2193", "2201", "2203", "2205",
		"2206", "2207", "2214", "2215", "2216", "2218", "2220", "2222", "2223", "2224",
		"2226", "2228", "2233", "2247", "2248", "2249", "2250", "2251", "2252", "2253",
		"2254", "2255", "2256", "2257", "2258", "2259", "2262", "2263", "2265", "2271",
		"2274", "2275", "2290", "2292", "2295"},
	// 2 Northeast Regional
	{"65", "66", "67", "82", "84", "85", "86", "87", "88", "93", "94", "95", "96", "99",
		"100", "101", "102", "103", "104", "106", "108", "109", "111", "112", "113", "114",
		"116", "117", "118", "119", "120", "121", "122", "124", "125", "126", "127", "128",
		"129", "130", "131", "132", "133", "134", "135", "136", "137", "138", "139", "140",
		"141", "142", "143", "144", "145", "146", "147", "148", "149", "150", "151", "152",
		"153", "154", "155", "156", "157", "158", "159", "160", "161", "162", "163", "164",
		"165", "166", "167", "168", "169", "170", "171", "172", "173", "174", "175", "176",
		"177", "178", "179", "181", "182", "183", "184", "185", "186", "189", "190", "192",
		"193", "194", "195", "196", "197", "198", "199", "400", "405", "409", "416", "417",
		"425", "426", "450", "460", "461", "463", "464", "465", "467", "470", "471", "473",
		"474", "475", "478", "479", "486", "488", "490", "494", "495", "497", "499",
		"627", "630", "631", "632", "636", "806", "807", "887", "888",
		"1108", "1161", "1175", "1194", "1195", "2107", "2117"},
	// 3 NEC Richmond/Newport News/Norfolk — sub-service, trains in corridor 2
	{},
	// 4 NEC Roanoke — sub-service, trains in corridor 2
	{},
	// 5 NEC Springfield — sub-service, trains in corridor 2
	{},
	// 6 Adirondack
	{"68", "69"},
	// 7 Blue Water
	{"364", "365", "1364"},
	// 8 Borealis
	{"1333", "1340"},
	// 9 Capitol Corridor
	{"520", "521", "522", "523", "524", "525", "526", "527", "528", "529", "530", "531",
		"532", "534", "535", "536", "537", "538", "539", "540", "541", "542", "543", "544",
		"545", "546", "547", "548", "549", "550", "551",
		"720", "723", "724", "727", "728", "729", "732", "733", "734", "736", "737", "738",
		"741", "742", "743", "744", "745", "746", "747", "748", "749", "750", "751"},
	// 10 Carl Sandburg / Illinois Zephyr
	{"380", "381", "382", "383"},
	// 11 Carolinian
	{"79", "80", "105", "1072", "1075", "1123", "1171", "1172"},
	// 12 Amtrak Cascades
	{"500", "502", "503", "504", "505", "506", "507", "508", "509", "511", "516", "517", "518", "519"},
	// 13 Downeaster
	{"680", "681", "682", "683", "684", "685", "686", "687", "688", "689", "690", "691",
		"692", "693", "694", "695", "696", "697", "698", "699", "1689", "1697"},
	// 14 Ethan Allen Express
	{"290", "291"},
	// 15 Gold Runner
	{"701", "702", "703", "704", "710", "711", "712", "713", "714", "715", "716", "717",
		"718", "719", "1701", "1702", "1703", "1704", "1710", "1711", "1712", "1715",
		"1716", "1717", "1718", "1719"},
	// 16 Heartland Flyer
	{"821", "822"},
	// 17 Hiawatha
	{"329", "330", "331", "332", "334", "335", "336", "337", "338", "339", "341", "342", "343"},
	// 18 Illini / Saluki
	{"390", "391", "392", "393"},
	// 19 Keystone
	{"110", "123", "600", "601", "605", "607", "609", "610", "611", "612", "615", "620",
		"622", "623", "624", "626", "637", "639", "640", "641", "642", "643", "644", "645",
		"646", "647", "648", "649", "650", "651", "652", "653", "654", "655", "656", "657",
		"658", "660", "661", "662", "663", "664", "665", "666", "667", "669", "670", "671",
		"672", "674"},
	// 20 Lincoln / Missouri River Runner
	{"300", "301", "302", "305", "306", "307", "311", "316", "318", "319"},
	// 21 Maple Leaf
	{"63", "64"},
	// 22 Mardi Gras Service
	{"23", "24", "25", "26"},
	// 23 Empire Service Albany — sub-corridor, trains in corridor 24
	{},
	// 24 Empire Service
	{"230", "232", "233", "234", "235", "236", "237", "238", "239", "240", "241",
		"243", "244", "245", "280", "281", "283", "284", "1237"},
	// 25 Pacific Surfliner
	{"562", "564", "566", "567", "572", "573", "577", "579", "580", "581", "582",
		"584", "586", "587", "588", "591", "593", "595",
		"757", "761", "765", "769", "770", "774", "777", "779", "782", "784", "785",
		"786", "790", "791", "794",
		"1562", "1591", "1595",
		"1765", "1769", "1770", "1774", "1777", "1784", "1785", "1790"},
	// 26 Pennsylvanian
	{"42", "43"},
	// 27 Pere Marquette
	{"370", "371"},
	// 28 Piedmont
	{"71", "72", "73", "74", "75", "76", "77", "78"},
	// 29 Vermonter
	{"54", "55", "56", "57", "107"},
	// 30 Wolverine
	{"350", "351", "352", "353", "354", "355", "1354"},
	// 31 Auto Train
	{"52", "53"},
	// 32 California Zephyr
	{"5", "6", "1005", "1006"},
	// 33 Cardinal
	{"50", "51"},
	// 34 City of New Orleans
	{"58", "59", "1058", "1059"},
	// 35 Coast Starlight
	{"11", "14", "1011", "1014"},
	// 36 Crescent
	{"19", "20", "1019", "1020"},
	// 37 Empire Builder
	{"7", "8", "27", "28", "1007", "1008", "1027", "1028"},
	// 38 Floridian
	{"40", "41", "1040", "1041"},
	// 39 Lake Shore Limited
	{"48", "49", "448", "449"},
	// 40 Palmetto
	{"89", "90"},
	// 41 Silver Meteor
	{"97", "98", "1098"},
	// 42 Southwest Chief
	{"3", "4", "1003", "1004"},
	// 43 Sunset Limited
	{"1", "2"},
	// 44 Texas Eagle
	{"21", "22", "1021", "1022"},
	// 45 Berkshire Flyer
	{"1233", "1234", "1246"},
	// 46 Winter Park Express
	{"1105", "1106"},
}

// otpSeeds maps corridor ID (1-based) → on-time percent.
// Source: FRA Quarterly Report FY2026 Q1 (Oct–Dec 2025).
var otpSeeds = map[int]float64{
	1:  76,
	2:  80,
	3:  76,
	4:  82,
	5:  88,
	6:  63,
	7:  69,
	8:  74,
	9:  89,
	10: 73,
	11: 71,
	12: 76,
	13: 86,
	14: 70,
	15: 63,
	16: 57,
	17: 87,
	18: 72,
	19: 90,
	20: 51,
	21: 59,
	22: 88,
	23: 86,
	24: 83,
	25: 85,
	26: 89,
	27: 82,
	28: 76,
	29: 72,
	30: 66,
	31: 83,
	32: 71,
	33: 56,
	34: 78,
	35: 68,
	36: 76,
	37: 61,
	38: 44,
	39: 78,
	40: 83,
	41: 67,
	42: 40,
	43: 75,
	44: 62,
}

func seedDB(db *sql.DB, adminUsername, adminPassword string) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM corridors`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("seeding database...")
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Seed corridors (insertion order determines auto-increment IDs 1..N)
	for i, c := range corridorSeeds {
		_, err := tx.Exec(
			`INSERT INTO corridors (name, slug, region, description, sort_order) VALUES (?, ?, ?, ?, ?)`,
			c.name, c.slug, c.region, c.description, i+1,
		)
		if err != nil {
			return fmt.Errorf("seed corridor %q: %w", c.name, err)
		}
	}

	// 2. Seed trains
	for i, nums := range trainSeeds {
		corridorID := i + 1
		for j, num := range nums {
			slug := "amtrak-" + num
			name := "Amtrak " + num
			_, err := tx.Exec(
				`INSERT INTO trains (corridor_id, train_number, display_name, slug, sort_order) VALUES (?, ?, ?, ?, ?)`,
				corridorID, num, name, slug, j+1,
			)
			if err != nil {
				return fmt.Errorf("seed train %s (corridor %d): %w", num, corridorID, err)
			}
		}
	}

	// 3. Seed stops from embedded data.sql
	if err := execSQLScript(tx, stopsSQL); err != nil {
		return fmt.Errorf("seed stops: %w", err)
	}

	// 3b. Seed Capitol Corridor train_stops with schedules and weekday/weekend flags
	if err := execSQLScript(tx, capitolCorridorStopsSQL); err != nil {
		return fmt.Errorf("seed capitol corridor train stops: %w", err)
	}

	// 4. Update OTP data
	for id, pct := range otpSeeds {
		tx.Exec(`UPDATE corridors SET on_time_percent=? WHERE id=?`, pct, id)
	}

	// 5. Seed admin user
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO admin_users (username, password_hash, must_change_password) VALUES (?, ?, 0)`,
		adminUsername, string(hash),
	); err != nil {
		return err
	}

	// 7. Site preferences
	if _, err := tx.Exec(`INSERT INTO site_preferences (id, default_theme) VALUES (1, 'auto')`); err != nil {
		return err
	}

	return tx.Commit()
}

// execSQLScript splits a SQL script by semicolons and executes each non-empty statement.
func execSQLScript(tx *sql.Tx, script string) error {
	for _, stmt := range strings.Split(script, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		// Skip blocks that are only comments
		hasCode := false
		for _, line := range strings.Split(stmt, "\n") {
			t := strings.TrimSpace(line)
			if t != "" && !strings.HasPrefix(t, "--") {
				hasCode = true
				break
			}
		}
		if !hasCode {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("%w\nSQL: %.200s", err, stmt)
		}
	}
	return nil
}
