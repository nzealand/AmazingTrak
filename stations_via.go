package main

// viaCorridorStops is the station list for the "VIA Rail Corridor" corridor —
// the Quebec City-Windsor Corridor (Toronto/Ottawa/Montreal/Quebec City
// area), which carries nearly all of VIA Rail's current train traffic.
// VIA's remaining long-distance/remote services (The Canadian:
// Toronto-Vancouver, the Ocean: Montreal-Halifax, Winnipeg-Churchill,
// Sudbury-White River, Jasper-Prince Rupert, Montreal-Jonquière/Senneterre)
// are deliberately out of scope — they run only 1-3x/week each and would add
// far more (often tiny, flag-stop) stations for proportionally little
// live-tracking value. Sourced from VIA Rail's own published static GTFS
// (https://viarail.ca/sites/all/files/gtfs/viarail.zip, Open Government
// Licence – Canada v2, confirmed empirically 2026-09-08): every stop_id
// reached by a trip on one of the 9 route_ids that make up the Corridor
// (Toronto-London, Montréal-Toronto, Toronto-Sarnia, Ottawa-Québec,
// Toronto-Windsor, Ottawa-Montréal, Ottawa-Toronto, Québec-Fallowfield,
// Québec-Montréal), plus Grimsby/St. Catharines/Niagara Falls Station added
// manually — those three are only reachable in the static feed via the
// Toronto-New York route (VIA 97/64 and 98/63, the joint VIA/Amtrak Maple
// Leaf), but VIA's own live tracker reports the Toronto-Niagara Falls leg
// under plain train numbers 97/98 (confirmed live 2026-09-08 — see
// trains_via.go), so those stations and numbers are included even though
// the through service to New York itself is Amtrak's to track (already
// covered as Amtrak trains 63/64). Ordered alphabetically, same idiom as
// metrolinkStops — the Corridor isn't one single line either.
var viaCorridorStops = []stationSeed{
	{"Aldershot", 43.312886, -79.855009, 10},
	{"Alexandria", 45.31805, -74.639672, 20},
	{"Belleville", 44.17961, -77.37455, 30},
	{"Brampton", 43.68685, -79.76448, 40},
	{"Brantford", 43.146686, -80.265709, 50},
	{"Brockville", 44.592113, -75.692713, 60},
	{"Casselman", 45.312287, -75.087197, 70},
	{"Charny", 46.713546, -71.270612, 80},
	{"Chatham", 42.397264, -82.179615, 90},
	{"Cobourg", 43.96778, -78.17151, 100},
	{"Cornwall", 45.042228, -74.743463, 110},
	{"Coteau", 45.27488, -74.23292, 120},
	{"Dorval", 45.44892, -73.74128, 130},
	{"Drummondville", 45.882431, -72.487905, 140},
	{"Fallowfield", 45.299303, -75.736708, 150},
	{"Gananoque", 44.36881, -76.15373, 160},
	{"Georgetown", 43.65551, -79.91912, 170},
	{"Glencoe", 42.745843, -81.71155, 180},
	{"Grimsby", 43.196023, -79.558179, 190},
	{"Guelph", 43.545251, -80.245602, 200},
	{"Guildwood", 43.754824, -79.198263, 210},
	{"Ingersoll", 43.04063, -80.887906, 220},
	{"Kingston", 44.257016, -76.536623, 230},
	{"Kitchener", 43.455741, -80.493207, 240},
	{"London", 42.98099, -81.246734, 250},
	{"Malton", 43.704865, -79.638141, 260},
	{"Montréal", 45.49992, -73.566179, 270},
	{"Napanee", 44.25365, -76.954262, 280},
	{"Niagara Falls Station", 43.109733, -79.054869, 290},
	{"Oakville", 43.454895, -79.682462, 300},
	{"Oshawa", 43.870828, -78.884626, 310},
	{"Ottawa", 45.41599, -75.65151, 320},
	{"Port Hope", 43.943592, -78.299017, 330},
	{"Québec", 46.818571, -71.214614, 340},
	{"Saint-Hyacinthe", 45.627708, -72.948931, 350},
	{"Saint-Lambert", 45.499079, -73.50746, 360},
	{"Sainte-Foy", 46.754238, -71.297575, 370},
	{"Sarnia", 42.95709, -82.38892, 380},
	{"Smiths Falls", 44.912455, -76.024023, 390},
	{"St. Catharines", 43.147664, -79.256377, 400},
	{"St. Marys", 43.260401, -81.13639, 410},
	{"Stratford", 43.364373, -80.975741, 420},
	{"Strathroy", 42.954682, -81.62292, 430},
	{"Toronto", 43.64481, -79.38032, 440},
	{"Trenton Junction", 44.115252, -77.598806, 450},
	{"Windsor", 42.324933, -83.007134, 460},
	{"Woodstock", 43.126528, -80.752123, 470},
	{"Wyoming", 42.94796, -82.12226, 480},
}
