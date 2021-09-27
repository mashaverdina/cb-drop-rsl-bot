package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Cb5          = "😈 Добавить дроп с 5 КБ"
	Cb6          = "👹 Добавить дроп с 6 КБ"
	Stats        = "📈 Посмотреть статистику"
	Approve      = "✅ OK"
	Reject       = "❌ Закрыть"
	Clear        = "🔄 Ввести заново"
	AncientShard = "💙 Древний осколок"
	VoidShard    = "💜 Темный осколок"
	SacredShard  = "💛 Сакральный осколок"
	EpicTome     = "📘 Эпический том"
	LegTome      = "📙 Легендарный том"
)

var HelloKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Cb5),
		tgbotapi.NewKeyboardButton(Cb6),
		tgbotapi.NewKeyboardButton(Stats),
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
