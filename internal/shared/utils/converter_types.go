package utils

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

func GetDayByTimestamp(unixTime int64) string {
	if unixTime == 0 {
		return ""
	}

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

func ConvertDateStrToDay(dateStr string) string {
	layout := "2 January,2006"

	replacer := strings.NewReplacer(
		"Januari", "January",
		"Februari", "February",
		"Maret", "March",
		"April", "April",
		"Mei", "May",
		"Juni", "June",
		"Juli", "July",
		"Agustus", "August",
		"September", "September",
		"Oktober", "October",
		"November", "November",
		"Desember", "December",
	)
	dateStr = replacer.Replace(dateStr)
	log.Printf("Tanggal: %s", dateStr)

	t, err := time.Parse(layout, strings.TrimSpace(dateStr))
	if err != nil {
		log.Printf("Gagal parse tanggal: %s, err: %v", dateStr, err)
		return "-"
	}

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

func ConvertDateStripToDay(dateStr string) string {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dateStr)

	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	day := t.Day()
	month := months[t.Month()-1]
	year := t.Year()

	return fmt.Sprintf("%02d %s, %d", day, month, year)
}

func ParseCookieDomain(clientHost string) string {
	u, err := url.Parse(clientHost)
	if err != nil {
		return ""
	}

	host := u.Hostname()

	if host == "localhost" || host == "127.0.0.1" {
		return ""
	}

	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return "." + strings.Join(parts[len(parts)-2:], ".")
	}

	return "." + host
}
