package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/storage"
)

const (
	JoinClan  = "joinclan"
	LeaveClan = "leaveclan"
)

type ClanCommand struct {
	BotCommand
	userStorage *storage.UserStorage
}

func NewClanCommand(userStorage *storage.UserStorage) *ClanCommand {
	return &ClanCommand{userStorage: userStorage}
}
func (c *ClanCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {

	return chatutils.TextTo(&user, "Готово", keyboards.MainMenuKeyboard), nil
}

func (c *ClanCommand) CanHandle(cmd string) bool {
	return cmd == JoinClan || cmd == LeaveClan
}
