package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/messages"
)

var MainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.Cb5),
		tgbotapi.NewKeyboardButton(messages.Cb6),
		tgbotapi.NewKeyboardButton(messages.Stats),
	),
)

var AddDropInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.AncientShard, messages.AncientShard),
		tgbotapi.NewInlineKeyboardButtonData(messages.VoidShard, messages.VoidShard),
		tgbotapi.NewInlineKeyboardButtonData(messages.SacredShard, messages.SacredShard),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.EpicTome, messages.EpicTome),
		tgbotapi.NewInlineKeyboardButtonData(messages.LegTome, messages.LegTome),
		tgbotapi.NewInlineKeyboardButtonData(messages.Clear, messages.Clear),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Approve, messages.Approve),
		tgbotapi.NewInlineKeyboardButtonData(messages.Reject, messages.Reject),
	),
)

var StatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.MonthStats, messages.MonthStats),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.LastVoidShard, messages.LastVoidShard),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.LastSacredShard, messages.LastSacredShard),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.LastLegTome, messages.LastLegTome),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Back, messages.Back),
	),
)

var ChooseMonthKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Jan, messages.Jan),
		tgbotapi.NewInlineKeyboardButtonData(messages.Feb, messages.Feb),
		tgbotapi.NewInlineKeyboardButtonData(messages.Mar, messages.Mar),
		tgbotapi.NewInlineKeyboardButtonData(messages.Apr, messages.Apr),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.May, messages.May),
		tgbotapi.NewInlineKeyboardButtonData(messages.Jun, messages.Jun),
		tgbotapi.NewInlineKeyboardButtonData(messages.Jul, messages.Jul),
		tgbotapi.NewInlineKeyboardButtonData(messages.Aug, messages.Aug),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Sep, messages.Sep),
		tgbotapi.NewInlineKeyboardButtonData(messages.Oct, messages.Oct),
		tgbotapi.NewInlineKeyboardButtonData(messages.Nov, messages.Nov),
		tgbotapi.NewInlineKeyboardButtonData(messages.Dec, messages.Dec),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Back, messages.Back),
	),
)
