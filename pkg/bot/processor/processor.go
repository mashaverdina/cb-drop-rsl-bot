package processor

import (
	"time"

	"vkokarev.com/rslbot/pkg/messages"
)

var monthMap = map[string]time.Month{
	messages.Jan: time.January,
	messages.Feb: time.February,
	messages.Mar: time.March,
	messages.Apr: time.April,
	messages.May: time.May,
	messages.Jun: time.June,
	messages.Jul: time.July,
	messages.Aug: time.August,
	messages.Sep: time.September,
	messages.Oct: time.October,
	messages.Nov: time.November,
	messages.Dec: time.December,
}

const (
	ancientSymbol = "💙"
	voidSymbol    = "💜"
	sacredSymbol  = "💛"
	epicSymbol    = "📘"
	legSymbol     = "📙"
)
