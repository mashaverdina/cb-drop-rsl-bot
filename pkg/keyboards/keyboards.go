package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Cb5             = "😈 Добавить дроп с 5 КБ"
	Cb6             = "👹 Добавить дроп с 6 КБ"
	Stats           = "📈 Посмотреть статистику"
	Approve         = "✅ OK"
	Reject          = "❌ Закрыть"
	Clear           = "🔄 Ввести заново"
	AncientShard    = "💙 Древний осколок"
	VoidShard       = "💜 Темный осколок"
	SacredShard     = "💛 Сакральный осколок"
	EpicTome        = "📘 Эпический том"
	LegTome         = "📙 Легендарный том"
	LastVoidShard   = "💜 Последний темный осколок"
	LastSacredShard = "💛 Последний сакральный осколок"
	LastLegTome     = "📙 Последний легендарный том"
	MonthStats      = "📅 Статистика за месяц"
	Back            = "🔙 Назад"
	Jan             = "Январь"
	Feb             = "Февраль"
	Mar             = "Март"
	Apr             = "Апрель"
	May             = "Май"
	Jun             = "Июнь"
	Jul             = "Июль"
	Aug             = "Август"
	Sep             = "Сентябрь"
	Oct             = "Октябрь"
	Nov             = "Ноябрь"
	Dec             = "Декабрь"
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

var StatsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Back),
		tgbotapi.NewKeyboardButton(MonthStats),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(LastVoidShard),
		tgbotapi.NewKeyboardButton(LastSacredShard),
		tgbotapi.NewKeyboardButton(LastLegTome),
	),
)
var ChooseMonthKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Jan),
		tgbotapi.NewKeyboardButton(Feb),
		tgbotapi.NewKeyboardButton(Mar),
		tgbotapi.NewKeyboardButton(Apr),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(May),
		tgbotapi.NewKeyboardButton(Jun),
		tgbotapi.NewKeyboardButton(Jul),
		tgbotapi.NewKeyboardButton(Aug),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Sep),
		tgbotapi.NewKeyboardButton(Oct),
		tgbotapi.NewKeyboardButton(Nov),
		tgbotapi.NewKeyboardButton(Dec),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Back),
	),
)
