package utils

import (
	"time"
)

var MSK, _ = time.LoadLocation("Europe/Moscow")

func MskNow() time.Time {
	return time.Now().In(MSK)
}
