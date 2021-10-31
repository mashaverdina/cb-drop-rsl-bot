package command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/storage"
	"vkokarev.com/rslbot/pkg/utils"
)

var timeRe = regexp.MustCompile("[0-9][0-9]?:[0-9]{2}")

type NotificationCommand struct {
	notificationStorage *storage.NotificationStorage
}

func NewNotificationCommand(notificationStorage *storage.NotificationStorage) *NotificationCommand {
	return &NotificationCommand{notificationStorage: notificationStorage}
}

func (c *NotificationCommand) Handle(ctx context.Context, user entities.User, commandText string, arguments string) ([]tgbotapi.Chattable, error) {
	notificationAlias := strings.TrimPrefix(commandText, utils.NotificationOn)
	notificationAlias = strings.TrimPrefix(notificationAlias, utils.NotificationOff)

	n, err := c.notificationStorage.GetByAlias(notificationAlias)
	if err != nil || n.Alias == "" {
		return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Нотификации с названием \"%s\" не найдено", arguments), keyboards.MainMenuKeyboard), nil
	}

	if strings.HasPrefix(commandText, utils.NotificationOn) {
		if !timeRe.MatchString(arguments) {
			return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Я ожидаю время в формате 13:30 Пример: /%s%s 13:30", utils.NotificationOn, n.Alias), keyboards.MainMenuKeyboard), nil
		}
		parts := strings.Split(arguments, ":")
		hour, _ := strconv.ParseInt(parts[0], 10, 64)
		minutes, _ := strconv.ParseInt(parts[1], 10, 64)
		if hour >= 24 {
			return chatutils.TextToNoMarkdown(&user, "Формат времени: 13:30, час не может быть больше 23", keyboards.MainMenuKeyboard), nil
		}
		if hour >= 60 {
			return chatutils.TextToNoMarkdown(&user, "Формат времени: 13:30, минуты не могут быть больше 59", keyboards.MainMenuKeyboard), nil
		}
		n, err = c.notificationStorage.LoadFire(n.NotificationID, int(hour), int(minutes))
		if err != nil {
			return chatutils.TextToNoMarkdown(&user, "Не получилось загрузить нотификацию, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		if err := c.notificationStorage.DisableNotification(user.UserID, n); err != nil {
			log.Println(fmt.Sprintf("error while disabling notification %v", err))
			return chatutils.TextToNoMarkdown(&user, "Не получилось удалить старую нотификацию, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		if err := c.notificationStorage.EnableNotification(user.UserID, n); err != nil {
			log.Println(fmt.Sprintf("error while disabling notification %v", err))
			return chatutils.TextToNoMarkdown(&user, "Не получилось удалить старую нотификацию, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Нотификация включена, для отключения введи (нажми) /%s%s", utils.NotificationOff, n.Alias), keyboards.MainMenuKeyboard), nil
	} else if strings.HasPrefix(commandText, utils.NotificationOff) {
		if err := c.notificationStorage.DisableNotification(user.UserID, n); err != nil {
			log.Println(fmt.Sprintf("error while disabling notification %v", err))
			return chatutils.TextToNoMarkdown(&user, "Не получилось удалить старую нотификацию, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Нотификация тебя больше не побеспокоит. Для включения: /%s%s 13:30", utils.NotificationOn, n.Alias), keyboards.MainMenuKeyboard), nil
	}
	return nil, errors.New("not applicable")
}

func (c *NotificationCommand) CanHandle(cmd string) bool {
	return strings.HasPrefix(cmd, utils.NotificationOff) || strings.HasPrefix(cmd, utils.NotificationOn)
}
