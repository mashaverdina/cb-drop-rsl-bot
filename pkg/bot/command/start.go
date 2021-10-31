package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
)

type StartCommand struct {
	BotCommand
}

func (c *StartCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	text := messages.HelpHeader +
		"Ты можешь прислать свои идеи по дальнейшему развитию [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
		"Ответы на часто задаваемые вопросы можешь посмотреть по команде /faq\n"
	return chatutils.TextTo(&user, text, keyboards.MainMenuKeyboard), nil
}

func (c *StartCommand) CanHandle(cmd string) bool {
	return cmd == "start"
}
