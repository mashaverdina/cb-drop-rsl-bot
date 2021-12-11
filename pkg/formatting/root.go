package formatting

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/utils"
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
	if stat.Level == 4 {
		return []cbStatPair{
			{messages.AncientShard, stat.AncientShard},
			{messages.VoidShard, stat.VoidShard},
			{messages.EpicTome, stat.EpicTome},
		}
	}
	return []cbStatPair{
		{messages.AncientShard, stat.AncientShard},
		{messages.VoidShard, stat.VoidShard},
		{messages.SacredShard, stat.SacredShard},
		{messages.EpicTome, stat.EpicTome},
		{messages.LegTome, stat.LegTome},
	}
}

type TopFunc func(string, int) string

func VerticalCbStat(stat entities.UserCbStat, topFuncs []TopFunc) string {
	lines := []string{}
	for _, p := range CbStatFields(stat) {
		line := fmt.Sprintf("%s --- %d", p.Title, p.Value)

		topStatLines := make([]string, 0)
		for _, topFunc := range topFuncs {
			if topLine := topFunc(p.Title, p.Value); topLine != "" {
				topStatLines = append(topStatLines, topLine)
			}
		}

		if len(topStatLines) > 0 {
			line = line + " (" + strings.Join(topStatLines, ", ") + ")"
		}

		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func VerticalCbStatWithHeader(stat entities.UserCbStat, topFuncs []TopFunc, headerPattern string, args ...interface{}) string {
	header := fmt.Sprintf(headerPattern, args...)
	return header + "\n" + VerticalCbStat(stat, topFuncs)
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

func CbStatsFormat(stats []entities.UserCbStat, withTime bool) string {
	lines := []string{}
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
	delta := utils.MskNow().Sub(*t)
	if delta.Hours() < 24 {
		return t.Format(dateFormat) + " (сегодня)"
	}
	return t.Format(dateFormat) + " (" + strconv.FormatInt(int64(delta.Hours()/24), 10) + " д. назад)"
}
