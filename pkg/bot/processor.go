package rslbot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/keyboards"
)

type Processor interface {
	Handle(state UserState, msg *tgbotapi.Message) (UserState, tgbotapi.Chattable, error)
}

type MainProcessor struct {
}

func (p *MainProcessor) Handle(state UserState, msg *tgbotapi.Message) (UserState, tgbotapi.Chattable, error) {
	switch msg.Text {
	case keyboards.Cb5:
		state.State = Cb5
		resp := tgbotapi.NewMessage(msg.Chat.ID, "Что упало с 5го КБ?")
		resp.ReplyMarkup = keyboards.NumericKeyboard
		return state, resp, nil
	case keyboards.Cb6:
		state.State = Cb6
		resp := tgbotapi.NewMessage(msg.Chat.ID, "Что упало с 6го КБ?")
		resp.ReplyMarkup = keyboards.NumericKeyboard
		return state, resp, nil
	case keyboards.Reject:
		state.State = MainMenu
		resp := tgbotapi.NewMessage(msg.Chat.ID, "До встречи")
		resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return state, resp, nil
	}

	resp := tgbotapi.NewMessage(msg.Chat.ID, "Привет")
	resp.ReplyMarkup = keyboards.HelloKeyboard
	return state, resp, nil
}

type CbProcessor struct {
	Level int
}

func (p *CbProcessor) Handle(state UserState, msg *tgbotapi.Message) (UserState, tgbotapi.Chattable, error) {
	log.Println("handling message from %v with CbProcessor", msg.From.UserName)

	switch msg.Text {
	case keyboards.Reject:
		state.State = MainMenu
		resp := tgbotapi.NewMessage(msg.Chat.ID, "До встречи")
		resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return state, resp, nil
	default:
		resp := tgbotapi.NewMessage(msg.Chat.ID, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)")
		return state, resp, nil
	}
}
