package utils

import "time"

func GetDayByTimestamp(unixTime int64) string {
	t := time.Unix(unixTime, 0).In(time.Local)
	switch t.Weekday() {
	case time.Monday:
		return "Senin"
	case time.Tuesday:
		return "Selasa"
	case time.Wednesday:
		return "Rabu"
	case time.Thursday:
		return "Kamis"
	case time.Friday:
		return "Jumat"
	case time.Saturday:
		return "Sabtu"
	case time.Sunday:
		return "Minggu"
	default:
		return "-"
	}
}
