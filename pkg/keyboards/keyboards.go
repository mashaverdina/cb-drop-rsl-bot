package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// generic controls
const (
	Approve = "✅ OK"
	Reject  = "❌ Закрыть"
	Back    = "🔙 Назад"
)

// main menu
const (
	Cb5   = "😈 Добавить дроп с 5 КБ"
	Cb6   = "👹 Добавить дроп с 6 КБ"
	Stats = "📈 Cтатистика"
)

// add drop inline menu
const (
	Clear        = "🔄 Заново"
	AncientShard = "💙 Древний"
	VoidShard    = "💜 Темный"
	SacredShard  = "💛 Сакрал"
	EpicTome     = "📘 Эпик том"
	LegTome      = "📙 Лег том"
)

// stats menu
const (
	LastVoidShard   = "💜 Последний темный"
	LastSacredShard = "💛 Последний сакрал"
	LastLegTome     = "📙 Последний лег том"
	MonthStats      = "📅 Дроп за месяц"
)

//month menu
const (
	Jan = "Январь"
	Feb = "Февраль"
	Mar = "Март"
	Apr = "Апрель"
	May = "Май"
	Jun = "Июнь"
	Jul = "Июль"
	Aug = "Август"
	Sep = "Сентябрь"
	Oct = "Октябрь"
	Nov = "Ноябрь"
	Dec = "Декабрь"
)

var MainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Cb5),
		tgbotapi.NewKeyboardButton(Cb6),
		tgbotapi.NewKeyboardButton(Stats),
	),
)

var AddDropInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
