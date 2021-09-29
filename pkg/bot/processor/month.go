package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
		replyMsgLines := []string{}
		from, to := mothInterval(msg.Text)
		for i := 5; i <= 6; i++ {
			monthStat, err := p.cbStatStorage.UserStat(ctx, msg.User.UserID, []int{i}, from, to)
			if err != nil || len(monthStat) == 0 {
				replyMsgLines = append(replyMsgLines, fmt.Sprintf("Статистики по %d кб пока нет", i))
				continue
			} else {
				replyMsgLines = append(
					replyMsgLines,
					formatting.CbStatsFormat(monthStat, true, "Твой дроп с *%d КБ* за *%s*", i, msg.Text),
				)
			}

			monthStatCombined, err := p.cbStatStorage.UserStatCombined(ctx, msg.User.UserID, []int{i}, from, to)
			if err == nil {
				replyMsgLines = append(replyMsgLines, formatting.VerticalCbStat(monthStatCombined), "")
			}
		}
		replyMsg := strings.Join(replyMsgLines, "\n")

		resp := chatutils.EditTo(msg, replyMsg, &keyboards.StatsKeyboard)
		return state, resp, nil
	default:
		resp := chatutils.TextTo(msg, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)", nil)
		return state, resp, nil
	}
}
