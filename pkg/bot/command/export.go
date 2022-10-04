package command

import (
	"context"
	"vkokarev.com/rslbot/pkg/export"
	"vkokarev.com/rslbot/pkg/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type ExportCommand struct {
	BotCommand
	exporter export.Exporter
}

func NewExportCommand(cbStorage *storage.CbStatStorage) *ExportCommand {
	return &ExportCommand{exporter: export.NewExcelExporter(cbStorage)}
}

func (c *ExportCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	text := "Экспорт сгенерирован!"
	fn, _, err := c.exporter.Export(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	doc := tgbotapi.NewDocumentUpload(user.Chat(), fn)
	result := chatutils.TextTo(&user, text, keyboards.MainMenuKeyboard)
	result = append(result, doc)
	return result, nil
}

func (c *ExportCommand) CanHandle(cmd string) bool {
	return cmd == "export"
}
