package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
)

type MainProcessor struct{}

func (p *MainProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Cb5:
		state.State = entities.StateCb5
		resp := chatutils.TextTo(msg, "Что упало с 5го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case messages.Cb6:
		state.State = entities.StateCb6
		resp := chatutils.TextTo(msg, "Что упало с 6го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case messages.Stats:
		state.State = entities.StateStats
		resp := chatutils.TextTo(msg, "Что тебе показать?", keyboards.StatsKeyboard)
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *MainProcessor) CanHandle(msg *ProcessingMessage) bool {
	switch msg.Text {
	case messages.Cb5, messages.Cb6, messages.Stats:
		return true
	default:
		return false
	}
}

func (p *MainProcessor) CancelFor(userID int64) {
}
