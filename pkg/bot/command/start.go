package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type StartCommand struct {
	BotCommand
}

func (c *StartCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	text := "Привет!🤖\n" +
		"Я бот для отслеживания твоего дропа с КБ в Raid SL.\n" +
		"Если ты забираешь *2 последних* сундука с *5 и/или 6* клан босса, отправляй мне информацию о своем дропе, и я запомню ее для тебя. А еще покажу тебе разную интересную статистику.\n\n" +
		"Пока что *я умею* показывать\n" +
		"  – твой дроп за месяц,\n" +
		"  – даты выпадения последних сакрала, войда и лег тома.\n" +
		"Но совсем *скоро я смогу*\n" +
		"  – экспортировать твой дроп в Excel файл,\n" +
		"  – показывать шансы дропа,\n" +
		"  – рассказывать тебе, насколько ты удачлив(а) по сравнению с сокланами и всеми пользователями.\n\n" +
		"Ты можешь прислать свои идеи по дальнейшему развитию [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
		"Ответы на часто задаваемые вопросы можешь посмотреть по команде /faq\n"
	return chatutils.TextTo(&user, text, keyboards.MainMenuKeyboard), nil
}

func (c *StartCommand) CanHandle(cmd string) bool {
	return cmd == "start"
}
