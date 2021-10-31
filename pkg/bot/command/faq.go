package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
)

type FAQCommand struct {
	BotCommand
}

func (c *FAQCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	text := "Привет!🤖\n" +
		"Вот ответы на *часто задаваемые вопросы*\n" +
		"Если здесь нет ответа на твой вопрос, напиши [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
		messages.HelpFAQ
	return chatutils.TextTo(&user, text, keyboards.MainMenuKeyboard), nil
}

func (c *FAQCommand) CanHandle(cmd string) bool {
	return cmd == "faq"
}
