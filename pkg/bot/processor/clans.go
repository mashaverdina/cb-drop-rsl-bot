package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/storage"
)

type ClansProcessor struct {
	cbStatStorage *storage.CbStatStorage
}

func NewClansProcessor(cbStatStorage *storage.CbStatStorage) *ClansProcessor {
	return &ClansProcessor{
		cbStatStorage: cbStatStorage,
	}
}

func (p *ClansProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Back:
		state.ProcType = entities.StateMainMenu
		resp := chatutils.DisableKeyboardAndSendNew(msg, "До встречи", keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.ClansEnter:
		state.ProcType = entities.StateMainMenu
		resp := chatutils.DisableKeyboardAndSendNew(msg, messages.ClansEnterMsg, keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.ClansExit:
		state.ProcType = entities.StateMainMenu
		user.Clan = ""
		if err := c.userStorage.Save(ctx, &user); err != nil {
			return chatutils.TextTo(&user, "Не получилось выйти из клана, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextTo(&user, fmt.Sprintf("Ты больше не состоишь ни в каком клане.\nДля присоединения к клану скопируй /%s  ИмяКлана", JoinClan), keyboards.MainMenuKeyboard), nil
		resp := chatutils.DisableKeyboardAndSendNew(msg, "Ты больше не состоишь ни в каком клане", keyboards.MainMenuKeyboard)
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *ClansProcessor) CancelFor(userID int64) {
}
