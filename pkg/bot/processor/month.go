package processor

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"vkokarev.com/rslbot/pkg/globalstat"
	"vkokarev.com/rslbot/pkg/utils"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/formatting"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/storage"
)

type PeriodType int

const (
	monthPeriod PeriodType = iota
	days7Period
	days30Period
)

type MonthProcessor struct {
	cbStatStorage     *storage.CbStatStorage
	globalStatManager *globalstat.GlobalStatManager
}

func NewMonthProcessor(cbStatStorage *storage.CbStatStorage, globalStatManager *globalstat.GlobalStatManager) *MonthProcessor {
	return &MonthProcessor{
		cbStatStorage:     cbStatStorage,
		globalStatManager: globalStatManager,
	}
}

func (p *MonthProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Back:
		state.ProcType = entities.StateStats
		state.Options.DropLevels()
		state.Options.DropShowFullStat()
		resp := chatutils.EditTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard)
		return state, resp, nil
	case messages.Jan, messages.Feb, messages.Mar, messages.Apr, messages.May, messages.Jun, messages.Jul, messages.Aug, messages.Sep, messages.Oct, messages.Nov, messages.Dec:
		state.ProcType = entities.StateStats
		from, to := utils.MonthInterval(monthMap[msg.Text])
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, msg.Text, monthPeriod, state.Options.Levels, state.Options.ShowFullStat)

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		state.Options.DropLevels()
		state.Options.DropShowFullStat()
		return state, resp, nil
	case messages.Days30:
		state.ProcType = entities.StateStats
		from, to := utils.LastDaysInterval(30)
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, "последние 30 дней", days30Period, state.Options.Levels, state.Options.ShowFullStat)

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		state.Options.DropLevels()
		state.Options.DropShowFullStat()
		return state, resp, nil
	case messages.Days7:
		state.ProcType = entities.StateStats
		from, to := utils.LastDaysInterval(7)
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, "последние 7 дней", days7Period, state.Options.Levels, state.Options.ShowFullStat)

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		state.Options.DropLevels()
		state.Options.DropShowFullStat()
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *MonthProcessor) CancelFor(userID int64) {
}

func (p *MonthProcessor) getPeriodDrop(ctx context.Context, userID int64, from time.Time, to time.Time, text string, periodType PeriodType, levels []int, showFullStat bool) string {
	replyMsgLines := []string{}
	for _, i := range levels {
		monthStat, err := p.cbStatStorage.UserStat(ctx, userID, []int{i}, from, to)
		if err != nil || len(monthStat) == 0 {
			replyMsgLines = append(replyMsgLines, fmt.Sprintf("Статистики по *%d кб* за *%s* пока нет", i, text))
			continue
		} else {
			replyMsgLines = append(
				replyMsgLines,
				fmt.Sprintf("Твой дроп с *%d КБ* за *%s*", i, text),
			)
			if showFullStat {
				replyMsgLines = append(
					replyMsgLines,
					formatting.CbStatsFormat(monthStat, true),
				)
			}
		}

		monthStatCombined, err := p.cbStatStorage.UserStatCombined(ctx, userID, []int{i}, from, to)

		if err == nil {
			replyMsgLines = append(replyMsgLines, formatting.VerticalCbStat(monthStatCombined, []formatting.TopFunc{
				func(level int) formatting.TopFunc {
					return func(itemType string, itemCount int) string {
						var top float64 = 0.
						var err error = nil
						switch periodType {
						case monthPeriod:
							top, err = p.globalStatManager.TopForMonth(from.Month(), level, itemType, itemCount)
						case days7Period:
							top, err = p.globalStatManager.TopFor7Days(level, itemType, itemCount)
						case days30Period:
							top, err = p.globalStatManager.TopFor30Days(level, itemType, itemCount)
						}
						if err != nil {
							return ""
						}
						if itemCount == 0 {
							return fmt.Sprintf("как и у других %.1f%%", math.Min(100., 2*top*100))
						}
						return fmt.Sprintf("больше, чем у %.1f%%", top*100)
					}
				}(i),
			}), "")
		}
	}
	return strings.Join(replyMsgLines, "\n")
}
