package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
)

type HelpCommand struct {
	BotCommand
}

func (c *HelpCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	text := messages.HelpHeader +
		"Ты можешь прислать свои идеи по дальнейшему развитию [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
		"А дальше я постараюсь ответить на *часто задаваемые вопросы*\n\n" +
		messages.HelpFAQ +
		"Если у тебя остались вопросы, напиши [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n"
	return chatutils.TextTo(&user, text, keyboards.MainMenuKeyboard), nil
}

func (c *HelpCommand) CanHandle(cmd string) bool {
	return cmd == "help"
}
