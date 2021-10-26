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
	JoinClan  = "joinclan"
	LeaveClan = "leaveclan"
	ClanInfo  = "clan"
)

type ClanCommand struct {
	BotCommand
	userStorage *storage.UserStorage
}

func NewClanCommand(userStorage *storage.UserStorage) *ClanCommand {
	return &ClanCommand{userStorage: userStorage}
}
func (c *ClanCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	switch command {
	case ClanInfo:
		if user.Clan == "" {
			return chatutils.TextTo(&user, fmt.Sprintf("Ты не состоишь ни в каком клане.\nДля присоединения к клану скопируй /%s  ИмяКлана", JoinClan), keyboards.MainMenuKeyboard), nil
		} else {
			return chatutils.TextTo(&user, fmt.Sprintf("Ты состоишь в клане \"%s\"\nДругие члены клана смогут видить твою стату. Для выхода из клана нажми /%s", user.Clan, LeaveClan), keyboards.MainMenuKeyboard), nil
		}
	case JoinClan:
		clanName := strings.Trim(arguments, " ")
		if clanName == "" {
			return chatutils.TextTo(&user, fmt.Sprintf("Укажи имя клана, в который хочешь вступить: /%s ИмяКлана", JoinClan), keyboards.MainMenuKeyboard), nil
		}
		user.Clan = clanName
		if err := c.userStorage.Save(ctx, &user); err != nil {
			return chatutils.TextTo(&user, "Не получилось сменить клан, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextTo(&user, fmt.Sprintf("Ты вступил в клан \"%s\"\nДругие члены клана смогут видить твою стату. Для выхода из клана нажми /%s", clanName, LeaveClan), keyboards.MainMenuKeyboard), nil

	case LeaveClan:
		user.Clan = ""
		if err := c.userStorage.Save(ctx, &user); err != nil {
			return chatutils.TextTo(&user, "Не получилось выйти из клана, попробуй еще раз или напиши нам в саппорт", keyboards.MainMenuKeyboard), nil
		}
		return chatutils.TextTo(&user, fmt.Sprintf("Ты больше не состоишь ни в каком клане.\nДля присоединения к клану скопируй /%s  ИмяКлана", JoinClan), keyboards.MainMenuKeyboard), nil
	}
	return nil, errors.New("not applicable")
}

func (c *ClanCommand) CanHandle(cmd string) bool {
	return cmd == JoinClan || cmd == LeaveClan || cmd == ClanInfo
}
