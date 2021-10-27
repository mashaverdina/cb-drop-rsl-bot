package utils

import (
	"time"
)

var MSK, _ = time.LoadLocation("Europe/Moscow")

func MskNow() time.Time {
	return time.Now().In(MSK)
}

func LastDaysInterval(daysN int) (time.Time, time.Time) {
	cy, cm, cd := MskNow().Date()
	to := time.Date(cy, cm, cd, 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, 0, -(daysN - 1))
	return from, to
}

func MonthInterval(mn time.Month) (time.Time, time.Time) {
	cy, cm, _ := MskNow().Date()
	if cm < mn {
		cy = cy - 1
	}

	from := time.Date(cy, mn, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return from, to
}
