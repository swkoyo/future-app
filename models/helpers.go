package models

import "time"

const GLOBAL_TZ = "PST"
const GLOBAL_TZ_OFFSET = -8 * 60 * 60

func ConvertToFixedTZ(t time.Time) time.Time {
	return t.In(time.FixedZone(GLOBAL_TZ, GLOBAL_TZ_OFFSET))
}

func ParseDateStr(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return ConvertToFixedTZ(parsedDate), nil
}
