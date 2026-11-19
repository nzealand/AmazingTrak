package main

import (
	"database/sql"
	_ "embed"
	"encoding/csv"
	"strconv"
	"strings"
)

// amtrak_schedules.csv holds per-train stop times for Amtrak corridors,
// generated from Amtrak's published static GTFS feed (feed_start 2026-11-19,
// valid through 2027-11-15). Amtrak publishes every time in Eastern, so the
// times here were converted to each station's local time. Columns:
// corridor,train,sort,stop,arr,dep,weekday,weekend. Blank arr = origin (or
// same as dep); blank dep = terminus. Where a train has more than one
// timetable variant, the most-operated one is used.
//
//go:embed amtrak_schedules.csv
var amtrakSchedulesCSV string

// seedAmtrakSchedules inserts the embedded schedules for trains that have no
// schedule rows yet (hand-entered schedules are never touched). Stops are
// matched by station code within the corridor, so IDs needn't match across
// databases. Idempotent: a train with rows is skipped on re-run.
func seedAmtrakSchedules(db *sql.DB) error {
	recs, err := csv.NewReader(strings.NewReader(amtrakSchedulesCSV)).ReadAll()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var curKey string
	var trainID int64
	var haveTrain bool
	for _, r := range recs[1:] {
		if len(r) < 8 {
			continue
		}
		key := r[0] + "|" + r[1]
		if key != curKey {
			curKey = key
			haveTrain = tx.QueryRow(`SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
				WHERE c.slug=? AND t.train_number=?
				AND NOT EXISTS (SELECT 1 FROM train_stops ts WHERE ts.train_id=t.id)`,
				r[0], r[1]).Scan(&trainID) == nil
		}
		if !haveTrain {
			continue
		}
		sortOrder, _ := strconv.Atoi(r[2])
		var arr, dep interface{}
		if r[4] != "" {
			arr = r[4]
		}
		if r[5] != "" {
			dep = r[5]
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO train_stops
			(train_id, stop_id, sort_order, scheduled_arrival, scheduled_departure, runs_weekday, runs_weekend)
			SELECT ?, s.id, ?, ?, ?, ?, ? FROM stops s JOIN corridors c ON c.id=s.corridor_id
			WHERE c.slug=? AND s.station_code=?`,
			trainID, sortOrder, arr, dep, r[6], r[7], r[0], r[3]); err != nil {
			return err
		}
	}
	return tx.Commit()
}
