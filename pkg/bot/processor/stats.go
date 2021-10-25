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

func (p *StatsProcessor) LastStat(ctx context.Context, msg *ProcessingMessage, resource string, header string) ([]tgbotapi.Chattable, error) {
	lastFrom5, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 5, resource)
	if err != nil {
		return nil, err
	}
	lastFrom6, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 6, resource)
	if err != nil {
		return nil, err
	}
	msgs := []string{header}
	if resource == "void_shard" {
		lastFrom4, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 4, resource)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, fmt.Sprintf("С 4го -- %s", formatting.TimePast(lastFrom4)))
	}

	msgs = append(msgs, fmt.Sprintf("С 5го -- %s", formatting.TimePast(lastFrom5)),
		fmt.Sprintf("С 6го -- %s", formatting.TimePast(lastFrom6)))

	resp := chatutils.JoinResp(
		chatutils.EditTo(msg, strings.Join(msgs, "\n"), nil),
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
		resp, err := p.LastStat(ctx, msg, "void_shard", messages.LastVoidShard)
		return state, resp, err
	case messages.LastSacredShard:
		state.State = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "sacred_shard", messages.LastSacredShard)
		return state, resp, err
	case messages.LastLegTome:
		state.State = entities.StateStats
		resp, err := p.LastStat(ctx, msg, "leg_tome", messages.LastLegTome)
		return state, resp, err
	case messages.MonthStats:
		state.State = entities.StateMonth
		return state, chatutils.EditTo(msg, "📅 Выбери период времени", keyboards.ChooseMonthKeyboard()), nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *StatsProcessor) CancelFor(userID int64) {
}
