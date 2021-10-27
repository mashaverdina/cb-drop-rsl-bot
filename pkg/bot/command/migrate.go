package command

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/notification"
	"vkokarev.com/rslbot/pkg/storage"
)

type MigrateCommand struct {
	BotCommand
	notificationStorage *storage.NotificationStorage
	notificationManager *notification.NotificationManager
}

func NewMigrateCommand(notificationStorage *storage.NotificationStorage, notificationManager *notification.NotificationManager) *MigrateCommand {
	return &MigrateCommand{notificationStorage: notificationStorage, notificationManager: notificationManager}
}
func (c *MigrateCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	if !user.HasSudo {
		return chatutils.TextTo(&user, "Это админская функция, а-та-та", keyboards.MainMenuKeyboard), nil
	}

	users, err := c.notificationStorage.AllUsers()
	if err != nil {
		return chatutils.TextTo(&user, fmt.Sprintf("Ошибка: %v", err), keyboards.MainMenuKeyboard), nil
	}
	disabledUsers := make(map[int64]interface{})
	if len(disabledUsers) == 0 {
		return chatutils.TextTo(&user, "Забыл задать disabledUsers", keyboards.MainMenuKeyboard), nil
	}
	for _, userID := range users {
		if _, disabled := disabledUsers[userID]; !disabled {
			if err := c.notificationManager.AssignDefaultNotifications(userID); err != nil {
				return chatutils.TextTo(&user, fmt.Sprintf("Ошибка: %v", err), keyboards.MainMenuKeyboard), nil
			}
		}
	}

	return chatutils.TextTo(&user, "Готово", keyboards.MainMenuKeyboard), nil
}

func (c *MigrateCommand) CanHandle(cmd string) bool {
	return cmd == "migrate"
}
