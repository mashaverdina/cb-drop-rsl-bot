package processor

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/formatting"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/storage"
)

type StatsProcessor struct {
	cbStatStorage *storage.CbStatStorage
}

func NewStatsProcessor(cbStatStorage *storage.CbStatStorage) *StatsProcessor {
	return &StatsProcessor{
		cbStatStorage: cbStatStorage,
	}
}

func (p *StatsProcessor) LastStat(ctx context.Context, msg *ProcessingMessage, resource string, header string) ([]tgbotapi.Chattable, error) {
	lastFrom5, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 5, resource)
	if err != nil {
		return nil, err
	}
	lastFrom6, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 6, resource)
	if err != nil {
		return nil, err
	}

	resp := chatutils.JoinResp(
		chatutils.EditTo(msg, strings.Join([]string{
			header,
			fmt.Sprintf("С 5го -- %s", formatting.TimePast(lastFrom5)),
			fmt.Sprintf("С 6го -- %s", formatting.TimePast(lastFrom6)),
		}, "\n"), nil),
		chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
	)
	return resp, nil
}

func (p *StatsProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Back:
		state.State = entities.StateMainMenu
		resp := chatutils.DisableKeyboardAndSendNew(msg, "До встречи", keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.LastVoidShard:
		state.State = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "void_shard", messages.LastVoidShard+" осколок")
		return state, resp, err
	case messages.LastSacredShard:
		state.State = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "sacred_shard", messages.LastSacredShard+" осколок")
		return state, resp, err
	case messages.LastLegTome:
		state.State = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "leg_tome", messages.LastLegTome)
		return state, resp, err
	case messages.MonthStats:
		state.State = entities.StateMonth
		return state, chatutils.EditTo(msg, "📅 Выбери месяц", &keyboards.ChooseMonthKeyboard), nil
	default:
		resp := chatutils.TextTo(msg, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)", nil)
		return state, resp, nil
	}
}

func mothInterval(month string) (time.Time, time.Time) {
	mn := monthMap[month]
	cy, cm, _ := time.Now().Date()
	if cm < mn {
		cy = cy - 1
	}

	from := time.Date(cy, mn, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return from, to
}
