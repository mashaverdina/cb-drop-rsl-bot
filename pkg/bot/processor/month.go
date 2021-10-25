package processor

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
	"vkokarev.com/rslbot/pkg/utils"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/formatting"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/storage"
)

type MonthProcessor struct {
	cbStatStorage *storage.CbStatStorage
}

func NewMonthProcessor(cbStatStorage *storage.CbStatStorage) *MonthProcessor {
	return &MonthProcessor{
		cbStatStorage: cbStatStorage,
	}
}

func (p *MonthProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Back:
		state.State = entities.StateStats
		resp := chatutils.EditTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard)
		return state, resp, nil
	case messages.Jan, messages.Feb, messages.Mar, messages.Apr, messages.May, messages.Jun, messages.Jul, messages.Aug, messages.Sep, messages.Oct, messages.Nov, messages.Dec:
		state.State = entities.StateStats
		from, to := monthInterval(msg.Text)
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, msg.Text)

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		return state, resp, nil
	case messages.Days30:
		state.State = entities.StateStats
		from, to := lastDaysInterval(30)
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, "последние 30 дней")

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		return state, resp, nil
	case messages.Days7:
		state.State = entities.StateStats
		from, to := lastDaysInterval(7)
		replyMsg := p.getPeriodDrop(ctx, msg.User.UserID, from, to, "последние 7 дней")

		resp := chatutils.JoinResp(
			chatutils.RemoveAndSendNew(msg, replyMsg, nil),
			chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
		)
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *MonthProcessor) CancelFor(userID int64) {
}

func (p *MonthProcessor) getPeriodDrop(ctx context.Context, userID int64, from time.Time, to time.Time, text string) string {
	replyMsgLines := []string{}
	for i := 4; i <= 6; i++ {
		monthStat, err := p.cbStatStorage.UserStat(ctx, userID, []int{i}, from, to)
		if err != nil || len(monthStat) == 0 {
			replyMsgLines = append(replyMsgLines, fmt.Sprintf("Статистики по *%d кб* за *%s* пока нет", i, text))
			continue
		} else {
			replyMsgLines = append(
				replyMsgLines,
				formatting.CbStatsFormat(monthStat, true, "Твой дроп с *%d КБ* за *%s*", i, text),
			)
		}

		monthStatCombined, err := p.cbStatStorage.UserStatCombined(ctx, userID, []int{i}, from, to)
		if err == nil {
			replyMsgLines = append(replyMsgLines, formatting.VerticalCbStat(monthStatCombined), "")
		}
	}
	return strings.Join(replyMsgLines, "\n")
}

func monthInterval(month string) (time.Time, time.Time) {
	mn := monthMap[month]
	cy, cm, _ := utils.MskNow().Date()
	if cm < mn {
		cy = cy - 1
	}

	from := time.Date(cy, mn, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return from, to
}

func lastDaysInterval(daysN int) (time.Time, time.Time) {
	cy, cm, cd := utils.MskNow().Date()
	to := time.Date(cy, cm, cd, 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, 0, -(daysN - 1))
	return from, to
}
