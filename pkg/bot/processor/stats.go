package processor

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"

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

func (p *StatsProcessor) LastStat(ctx context.Context, msg *ProcessingMessage, resource string, header string, levels []int) ([]tgbotapi.Chattable, error) {
	msgs := []string{header}
	for _, level := range levels {
		last, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, level, resource)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, fmt.Sprintf("С %dго -- %s", level, formatting.TimePast(last)))
	}

	resp := chatutils.JoinResp(
		chatutils.EditTo(msg, strings.Join(msgs, "\n"), nil),
		chatutils.TextTo(msg, "Что тебе показать?", &keyboards.StatsKeyboard),
	)
	return resp, nil
}

func (p *StatsProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Back:
		state.ProcType = entities.StateMainMenu
		resp := chatutils.DisableKeyboardAndSendNew(msg, "До встречи", keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.LastVoidShard:
		state.ProcType = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "void_shard", messages.LastVoidShard, []int{4, 5, 6})
		return state, resp, err
	case messages.LastSacredShard:
		state.ProcType = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "sacred_shard", messages.LastSacredShard, []int{5, 6})
		return state, resp, err
	case messages.LastLegTome:
		state.ProcType = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "leg_tome", messages.LastLegTome, []int{5, 6})
		return state, resp, err
	case messages.MonthStats:
		state.ProcType = entities.StateMonth
		state.Options.WithLevels(4, 5, 6)
		return state, chatutils.EditTo(msg, "📅 Выбери период времени", keyboards.ChooseMonthKeyboard()), nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *StatsProcessor) CancelFor(userID int64) {
}
