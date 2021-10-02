package command

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/storage"
)

const (
	NotificationOn  = "notification_on_"
	NotificationOff = "notification_off_"
)

type NotificationCommand struct {
	notificationStorage *storage.NotificationStorage
}

func NewNotificationCommand(notificationStorage *storage.NotificationStorage) *NotificationCommand {
	return &NotificationCommand{notificationStorage: notificationStorage}
}

func (c *NotificationCommand) Handle(ctx context.Context, user entities.User, commandText string, arguments string) ([]tgbotapi.Chattable, error) {
	arguments = strings.TrimPrefix(commandText, NotificationOn)
	arguments = strings.TrimPrefix(arguments, NotificationOff)

	n, err := c.notificationStorage.GetByAlias(arguments)
	if err != nil || n.Alias == "" {
		return chatutils.TextTo(&user, fmt.Sprintf("Нотификации с названием \"%s\" не найдено", arguments), keyboards.MainMenuKeyboard), nil
	}
	if strings.HasPrefix(commandText, NotificationOn) {
		if err := c.notificationStorage.EnableNotification(entities.DisabledNotifications{
			UserID:         user.UserID,
			NotificationID: n.ID,
		}); err != nil {
			return chatutils.TextTo(&user, "Не получилось включить нотификацию попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Нотификация включена, для отключения введи (нажми) /%s%s", NotificationOff, n.Alias), keyboards.MainMenuKeyboard), nil
	} else if strings.HasPrefix(commandText, NotificationOff) {
		if err := c.notificationStorage.DisableNotification(entities.DisabledNotifications{
			UserID:         user.UserID,
			NotificationID: n.ID,
		}); err != nil {
			return chatutils.TextTo(&user, "Не получилось выключить нотификацию попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextToNoMarkdown(&user, fmt.Sprintf("Нотификация тебя больше не побеспокоит. Для включения: /%s%s", NotificationOn, n.Alias), keyboards.MainMenuKeyboard), nil
	}
	return nil, errors.New("not applicable")
}

func (c *NotificationCommand) CanHandle(cmd string) bool {
	return strings.HasPrefix(cmd, NotificationOff) || strings.HasPrefix(cmd, NotificationOn)
}
