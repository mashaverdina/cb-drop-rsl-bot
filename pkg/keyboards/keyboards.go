package keyboards

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/utils"
)

var MainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.Cb5),
		tgbotapi.NewKeyboardButton(messages.Cb6),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.Cb4),
		tgbotapi.NewKeyboardButton(messages.Stats),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.Help),
	),
)

func ChooseAddDropInlineKeyboard(level int) *tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, 0, 3)

	// 1st row: shards
	currentButtons := make([]tgbotapi.InlineKeyboardButton, 0, 3)
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.AncientShard, messages.AncientShard))
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.VoidShard, messages.VoidShard))
	if level != 4 {
		currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.SacredShard, messages.SacredShard))
	}
	markup = append(markup, currentButtons)

	// 2st row: tomes and clear
	currentButtons = make([]tgbotapi.InlineKeyboardButton, 0, 3)
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.EpicTome, messages.EpicTome))
	if level != 4 {
		currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.LegTome, messages.LegTome))
	}
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.Clear, messages.Clear))
	markup = append(markup, currentButtons)

	// 3rd row: ok, nothing, cancel
	currentButtons = make([]tgbotapi.InlineKeyboardButton, 0, 3)
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.Approve, messages.Approve))
	if level != 6 {
		currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.Nothing, messages.Nothing))
	}
	currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(messages.Reject, messages.Reject))
	markup = append(markup, currentButtons)

	r := tgbotapi.NewInlineKeyboardMarkup(markup...)
	return &r
}

var StatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.MonthStatsShort, messages.MonthStatsShort),
	),
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

var allMonth = []string{
	"",
	messages.Jan,
	messages.Feb,
	messages.Mar,
	messages.Apr,
	messages.May,
	messages.Jun,
	messages.Jul,
	messages.Aug,
	messages.Sep,
	messages.Oct,
	messages.Nov,
	messages.Dec,
}

func ChooseMonthKeyboard() *tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, 0, 6)

	markup = append(markup, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Days30, messages.Days30),
		tgbotapi.NewInlineKeyboardButtonData(messages.Days7, messages.Days7),
	))

	currentButtons := make([]tgbotapi.InlineKeyboardButton, 0, 3)

	_, realMonth, _ := utils.MskNow().Date()
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			currentButtons = append(currentButtons, tgbotapi.NewInlineKeyboardButtonData(allMonth[realMonth], allMonth[realMonth]))
			nextMonth(&realMonth)
		}
		markup = append(markup, currentButtons)
		currentButtons = make([]tgbotapi.InlineKeyboardButton, 0, 3)
	}

	markup = append(markup, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Back, messages.Back),
	))

	r := tgbotapi.NewInlineKeyboardMarkup(markup...)
	return &r
}

func nextMonth(m *time.Month) {
	if *m == time.December {
		*m = time.January
	} else {
		*m = *m + 1
	}

}
