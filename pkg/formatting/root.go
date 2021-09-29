package formatting

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/messages"
)

const (
	dateFormat = "02.01.2006"
)

const (
	ancientSymbol = "💙"
	voidSymbol    = "💜"
	sacredSymbol  = "💛"
	epicSymbol    = "📘"
	legSymbol     = "📙"
)

var ShortMapping = map[string]string{
	messages.AncientShard: ancientSymbol,
	messages.VoidShard:    voidSymbol,
	messages.SacredShard:  sacredSymbol,
	messages.EpicTome:     epicSymbol,
	messages.LegTome:      legSymbol,
}

type cbStatPair struct {
	Title string
	Value int
}

func CbStatFields(stat entities.UserCbStat) []cbStatPair {
	return []cbStatPair{
		{messages.AncientShard, stat.AncientShard},
		{messages.VoidShard, stat.VoidShard},
		{messages.SacredShard, stat.SacredShard},
		{messages.EpicTome, stat.EpicTome},
		{messages.LegTome, stat.LegTome},
	}
}

func VerticalCbStat(stat entities.UserCbStat) string {
	lines := []string{}
	for _, p := range CbStatFields(stat) {
		line := fmt.Sprintf("%s --- %d", p.Title, p.Value)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func VerticalCbStatWithHeader(stat entities.UserCbStat, headerPattern string, args ...interface{}) string {
	header := fmt.Sprintf(headerPattern, args...)
	return header + "\n" + VerticalCbStat(stat)
}

func HorizontalCbStat(stat entities.UserCbStat, mapping map[string]string) string {
	line := ""
	for _, p := range CbStatFields(stat) {
		if mapped, ok := mapping[p.Title]; ok {
			line = line + Multiple(mapped, p.Value)
		} else {
			line = line + Multiple(p.Title, p.Value)
		}
	}
	return line
}

func CbStatsFormat(stats []entities.UserCbStat, withTime bool, headerPattern string, args ...interface{}) string {
	lines := []string{fmt.Sprintf(headerPattern, args...)}
	for _, stat := range stats {
		line := HorizontalCbStat(stat, ShortMapping)
		if withTime {
			prefix := stat.RelatedTo.Format(dateFormat)
			line = fmt.Sprintf("_%s_ %s", prefix, line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func Multiple(data string, count int) string {
	s := ""
	for i := 0; i < count; i++ {
		s = s + " " + data
	}
	return s
}

func TimePast(t *time.Time) string {
	if t == nil {
		return "никогда"
	}
	delta := time.Now().Sub(*t)
	if delta.Hours() < 24 {
		return t.Format(dateFormat) + " (сегодня)"
	}
	return t.Format(dateFormat) + " (" + strconv.FormatInt(int64(delta.Hours()/24), 10) + " д. назад)"
}
