package main

// viaCorridorTrainNumbers is the train-number roster for the single
// "VIA Rail Corridor" corridor (Quebec City-Windsor Corridor only — see
// stations_via.go for why the rest of VIA's network is out of scope).
// Sourced from VIA Rail's own published static GTFS
// (https://viarail.ca/sites/all/files/gtfs/viarail.zip, confirmed
// empirically 2026-09-08): trip_short_name (the public train number) for
// every trip on the 9 Corridor route_ids, plus 97/98 added manually (see
// stations_via.go's comment on the Toronto-Niagara Falls leg of the Maple
// Leaf). Checked for collisions against every one of VIA's other route_ids
// (The Canadian, the Ocean, Winnipeg-Churchill, Sudbury-White River,
// Jasper-Prince Rupert, Montréal-Jonquière/Senneterre) — none of those
// remote services' numbers overlap this list, so scoping the roster to the
// Corridor alone can't misattribute a remote train's live position here.
var viaCorridorTrainNumbers = []string{
	"20", "22", "24", "26", "28", "29", "31", "33", "35", "37", "38", "39", "40", "41", "42",
	"43", "44", "45", "46", "47", "48", "50", "52", "53", "54", "55", "59", "60", "61", "62",
	"63", "64", "65", "66", "67", "68", "69", "70", "71", "72", "73", "75", "76", "78", "79",
	"82", "83", "84", "87", "97", "98", "622", "624", "633", "637", "641", "643", "644", "645",
	"646", "647", "668", "669",
}
