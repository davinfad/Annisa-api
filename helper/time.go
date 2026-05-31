package helper

import "time"

// wibLoc is Western Indonesia Time (UTC+7). The DB stores created_at as WIB
// wall-clock, but the MySQL driver (loc=Local=UTC on Railway) reads/writes it
// as if it were UTC. These helpers convert at the boundaries so that:
//   - storage stays WIB wall-clock (matching existing data), and
//   - the JSON API always emits true UTC ending in "Z" (what the app parses).
var wibLoc = time.FixedZone("WIB", 7*3600)

// WIBStoredToUTC reinterprets a created_at value read from the DB — whose
// wall-clock digits are actually WIB but were tagged UTC by the driver — as the
// correct UTC instant for JSON serialization. e.g. stored "15:59" -> "08:59Z".
func WIBStoredToUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), wibLoc).UTC()
}

// NowWIBStore returns the current time as a value that the driver will persist
// as WIB wall-clock digits (e.g. "15:59"), keeping new rows consistent with the
// existing WIB-stored data. Its formatted clock is also the correct WIB local
// time for business-hour checks.
func NowWIBStore() time.Time {
	now := time.Now().In(wibLoc)
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
}
