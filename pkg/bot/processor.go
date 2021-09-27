package rslbot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/keyboards"
)

type ProcessingMessage struct {
	UserID    int64
	ChatID    int64
	MessageID int
	Text      string
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
		resp.ReplyMarkup = keyboards.AddDropInlineKeyboard
		return state, resp, nil
	case keyboards.Cb6:
		state.State = Cb6
		resp := tgbotapi.NewMessage(msg.ChatID, "Что упало с 6го КБ?")
		resp.ReplyMarkup = keyboards.AddDropInlineKeyboard
		return state, resp, nil
	case keyboards.Reject:
		state.State = MainMenu
		resp := tgbotapi.NewMessage(msg.ChatID, "До встречи")
		resp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return state, resp, nil
	}

	resp := tgbotapi.NewMessage(msg.ChatID, "Привет")
	resp.ReplyMarkup = keyboards.MainMenuKeyboard
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
		p.stats[state.UserID] = NewCbUserState(state.UserID)
		return state, resp, nil
	case keyboards.Clear:
		p.stats[state.UserID] = NewCbUserState(state.UserID)
	case keyboards.LegTome:
		p.increment(&cbState.LegTome)
	case keyboards.AncientShard:
		p.increment(&cbState.AncientShard)
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

	resp := tgbotapi.NewEditMessageText(msg.ChatID, msg.MessageID, p.msgFromStat(cbState))
	resp.ReplyMarkup = &keyboards.AddDropInlineKeyboard
	resp.ParseMode = "markdown"
	p.stats[state.UserID] = cbState
	return state, resp, nil

}

func (p *CbProcessor) msgFromStat(state CbUserState) string {
	lines := []string{
		fmt.Sprintf("Стата по *%d КБ*", p.level),
	}

	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.AncientShard, state.AncientShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.VoidShard, state.VoidShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.SacredShard, state.SacredShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.EpicTome, state.EpicTome))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.LegTome, state.LegTome))

	return strings.Join(lines, "\n")
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
