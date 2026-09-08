package main

// Train-number rosters for VIA Rail's remote/long-distance services, each
// its own corridor (not merged into "VIA Rail Corridor") since they're
// distinct named services with their own branding and schedules, matching
// the same one-corridor-per-named-service idiom already used for Amtrak's
// own long-distance trains elsewhere in corridorSeeds (seed.go). Sourced from
// VIA Rail's own published static GTFS
// (https://viarail.ca/sites/all/files/gtfs/viarail.zip, confirmed empirically
// 2026-09-08): trip_short_name for every trip on each service's route_id(s).
// Checked for collisions against the Corridor roster (trains_via.go) and
// against each other — none overlap, so via.go's live matcher can use one flat
// global number index across all 8 VIA corridors without route-name
// disambiguation (unlike Amtrak's matchTrain).

// viaCanadianTrainNumbers: 2 unique train numbers (route_id 8-119).
var viaCanadianTrainNumbers = []string{
	"1", "2",
}

// viaOceanTrainNumbers: 2 unique train numbers (route_id 226-620).
var viaOceanTrainNumbers = []string{
	"14", "15",
}

// viaWinnipegChurchillTrainNumbers: 4 unique train numbers (route_ids 388-435, 149-435).
var viaWinnipegChurchillTrainNumbers = []string{
	"690", "691", "692", "693",
}

// viaSudburyWhiteRiverTrainNumbers: 2 unique train numbers (route_id 621-116).
var viaSudburyWhiteRiverTrainNumbers = []string{
	"185", "186",
}

// viaJasperPrinceRupertTrainNumbers: 2 unique train numbers (route_id 21-458).
var viaJasperPrinceRupertTrainNumbers = []string{
	"5", "6",
}

// viaMontrealJonquiereTrainNumbers: 4 unique train numbers (route_id 226-444).
var viaMontrealJonquiereTrainNumbers = []string{
	"600", "601", "602", "605",
}

// viaMontrealSenneterreTrainNumbers: 4 unique train numbers (route_id 226-460).
var viaMontrealSenneterreTrainNumbers = []string{
	"603", "604", "606", "607",
}
