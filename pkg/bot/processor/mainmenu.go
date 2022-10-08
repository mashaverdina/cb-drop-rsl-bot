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
		state.ProcType = entities.StateCb5
		state.Options.WithLevels(5)
		resp := chatutils.TextTo(msg, "Что упало с 5го КБ?", keyboards.ChooseAddDropInlineKeyboard(5))
		return state, resp, nil
	case messages.Cb6:
		state.ProcType = entities.StateCb6
		state.Options.WithLevels(6)
		resp := chatutils.TextTo(msg, "Что упало с 6го КБ?", keyboards.ChooseAddDropInlineKeyboard(6))
		return state, resp, nil
	case messages.Cb4:
		state.ProcType = entities.StateCb4
		state.Options.WithLevels(4)
		resp := chatutils.TextTo(msg, "Что упало с 4го КБ?", keyboards.ChooseAddDropInlineKeyboard(4))
		return state, resp, nil
	case messages.Stats:
		state.ProcType = entities.StateStats
		resp := chatutils.TextTo(msg, "Что тебе показать?", keyboards.StatsKeyboard)
		return state, resp, nil
	case messages.Clans:
		state.ProcType = entities.StateClans
		resp := chatutils.TextTo(msg, "Добро пожаловать в меню клана", keyboards.ClansKeyboard)
		return state, resp, nil
	case messages.FullStats:
		resp := chatutils.TextTo(msg, messages.FullStatsMsg, keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.Help:
		text := messages.HelpHeader +
			"Ты можешь прислать свои идеи по дальнейшему развитию [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
			"А дальше я постараюсь ответить на *часто задаваемые вопросы*\n\n" +
			messages.HelpFAQ +
			"Если у тебя остались вопросы, напиши [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n"
		resp := chatutils.TextTo(msg, text, keyboards.MainMenuKeyboard)
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *MainProcessor) CanHandle(msg *ProcessingMessage) bool {
	switch msg.Text {
	case messages.Cb4, messages.Cb5, messages.Cb6, messages.Stats, messages.FullStats, messages.Clans, messages.Help:
		return true
	default:
		return false
	}
}

func (p *MainProcessor) CancelFor(userID int64) {
}
