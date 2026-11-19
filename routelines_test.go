package main

import (
	"encoding/json"
	"testing"
)

// shapesToFeatureCollection must drop duplicate and reversed shapes and
// keep a line's endpoints while thinning collinear points.
func TestShapesToFeatureCollectionDedupes(t *testing.T) {
	line := [][2]float64{{-74.0, 40.0}, {-74.001, 40.001}, {-74.002, 40.002}, {-74.01, 40.01}}
	rev := [][2]float64{line[3], line[2], line[1], line[0]}
	other := [][2]float64{{-75.0, 41.0}, {-75.1, 41.1}}
	data := shapesToFeatureCollection(map[string][][2]float64{"a": line, "b": line, "c": rev, "d": other}, "T")

	var fc struct {
		Features []struct {
			Geometry struct {
				Coordinates [][2]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(data, &fc); err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) != 2 {
		t.Fatalf("want 2 distinct features, got %d", len(fc.Features))
	}
	for _, f := range fc.Features {
		c := f.Geometry.Coordinates
		if len(c) < 2 {
			t.Fatalf("degenerate line %v", c)
		}
	}
}
