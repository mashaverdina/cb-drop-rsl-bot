package rslbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/keyboards"
)

type ProcessingMessage struct {
	UserID int64
	ChatID int64
	Text   string
}

type Processor interface {
	Handle(state UserState, msg *ProcessingMessage) (UserState, tgbotapi.Chattable, error)
}

type MainProcessor struct {
}

func (p *MainProcessor) Handle(state UserState, msg *ProcessingMessage) (UserState, tgbotapi.Chattable, error) {
	switch msg.Text {
	case keyboards.Cb5:
		state.State = Cb5
		resp := tgbotapi.NewMessage(msg.ChatID, "Что упало с 5го КБ?")
		resp.ReplyMarkup = keyboards.NumericKeyboard
		return state, resp, nil
	case keyboards.Cb6:
		state.State = Cb6
		resp := tgbotapi.NewMessage(msg.ChatID, "Что упало с 6го КБ?")
		resp.ReplyMarkup = keyboards.NumericKeyboard
		return state, resp, nil
	case keyboards.Reject:
		state.State = MainMenu
		resp := tgbotapi.NewMessage(msg.ChatID, "До встречи")
		resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return state, resp, nil
	}

	resp := tgbotapi.NewMessage(msg.ChatID, "Привет")
	resp.ReplyMarkup = keyboards.HelloKeyboard
	return state, resp, nil
}

type CbProcessor struct {
	level int
	stats map[int64]CbUserState
}

func NewCbProcessor(level int) *CbProcessor {
	return &CbProcessor{
		level: level,
		stats: make(map[int64]CbUserState),
	}
}

func (p *CbProcessor) Handle(state UserState, msg *ProcessingMessage) (UserState, tgbotapi.Chattable, error) {
	cbState := p.getOrCreateStats(state.UserID)
	switch msg.Text {
	case keyboards.Reject:
		state.State = MainMenu
		resp := tgbotapi.NewMessage(msg.ChatID, "До встречи")
		resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return state, resp, nil
	case keyboards.LegTome:
		p.increment(&cbState.LegTome)
	case keyboards.AncientShard:
		p.increment(&cbState.SacredShard)
	case keyboards.VoidShard:
		p.increment(&cbState.VoidShard)
	case keyboards.SacredShard:
		p.increment(&cbState.SacredShard)
	case keyboards.EpicTome:
		p.increment(&cbState.EpicTome)
	default:
		resp := tgbotapi.NewMessage(msg.ChatID, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)")
		return state, resp, nil
	}

	resp := tgbotapi.NewMessage(msg.ChatID, msgFromStat(cbState))
	p.stats[state.UserID] = cbState
	return state, resp, nil

}

func msgFromStat(state CbUserState) string {
	return fmt.Sprintf("Итого:\n - Синих шародов - %d\n - Лег книжек - %d", state.SacredShard, state.LegTome)
}

func (p *CbProcessor) getOrCreateStats(userID int64) CbUserState {
	if s, ok := p.stats[userID]; ok {
		return s
	}
	s := NewCbUserState(userID)
	p.stats[userID] = s
	return s
}

func (p *CbProcessor) increment(val *int) {
	*val = *val + 1
}
