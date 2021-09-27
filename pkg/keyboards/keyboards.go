package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Cb5          = "5 КБ"
	Cb6          = "6 КБ"
	Stats        = "Стата"
	Approve      = "OK"
	Reject       = "Закрыть"
	Clear        = "Ввести заново"
	AncientShard = "Древний осколок"
	VoidShard    = "Темный осколок"
	SacredShard  = "Сакральный осколок"
	EpicTome     = "Эпический том"
	LegTome      = "Легендарный том"
)

var HelloKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Cb5),
		tgbotapi.NewKeyboardButton(Cb6),
		tgbotapi.NewKeyboardButton(Stats),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Approve),
		tgbotapi.NewKeyboardButton(Reject),
	),
)

var NumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(AncientShard, AncientShard),
		tgbotapi.NewInlineKeyboardButtonData(VoidShard, VoidShard),
		tgbotapi.NewInlineKeyboardButtonData(SacredShard, SacredShard),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(EpicTome, EpicTome),
		tgbotapi.NewInlineKeyboardButtonData(LegTome, LegTome),
		tgbotapi.NewInlineKeyboardButtonData(Clear, Clear),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Approve, Approve),
		tgbotapi.NewInlineKeyboardButtonData(Reject, Reject),
	),
)
