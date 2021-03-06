package processor

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/formatting"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
	"vkokarev.com/rslbot/pkg/storage"
)

type CbProcessor struct {
	level   int
	stats   map[int64]entities.UserCbStat
	storage *storage.CbStatStorage
	m       sync.Mutex
}

func NewCbProcessor(level int, storage *storage.CbStatStorage) *CbProcessor {
	return &CbProcessor{
		level:   level,
		stats:   make(map[int64]entities.UserCbStat),
		storage: storage,
		m:       sync.Mutex{},
	}
}

func (p *CbProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	format := func(cbStat entities.UserCbStat) string {
		return formatting.VerticalCbStatWithHeader(cbStat, []formatting.TopFunc{}, "Твой дроп с *%d КБ*", p.level)
	}

	cbStat := p.getOrCreateStats(state.UserID)
	switch msg.Text {
	case messages.Reject:
		state.ProcType = entities.StateMainMenu
		p.deleteUserStat(state.UserID)

		state.Options.Clear()
		resp := chatutils.RemoveAndSendNew(msg, "До встречи", keyboards.MainMenuKeyboard)
		return state, resp, nil
	case messages.Approve:
		state.ProcType = entities.StateMainMenu
		cbStat := p.getOrCreateStats(state.UserID)
		err := p.storage.Save(ctx, &cbStat)
		if err != nil {
			return entities.UserState{}, nil, fmt.Errorf("cb state db update failed: %v", err)
		}

		p.updateStats(entities.NewCbUserState(state.UserID, p.level))

		state.Options.Clear()
		resp := chatutils.JoinResp(
			chatutils.EditTo(msg, format(cbStat), nil),
			chatutils.TextTo(msg, "Записано", keyboards.MainMenuKeyboard),
		)
		return state, resp, nil
	case messages.Nothing:
		state.ProcType = entities.StateMainMenu
		cbStat := entities.NewCbUserState(state.UserID, p.level)
		err := p.storage.Save(ctx, &cbStat)
		if err != nil {
			return entities.UserState{}, nil, fmt.Errorf("cb state db update failed: %v", err)
		}

		p.updateStats(entities.NewCbUserState(state.UserID, p.level))

		state.Options.Clear()
		resp := chatutils.JoinResp(
			chatutils.EditTo(msg, format(cbStat), nil),
			chatutils.TextTo(msg, "Записано", keyboards.MainMenuKeyboard),
		)
		return state, resp, nil
	case messages.Clear:
		cbStat = entities.NewCbUserState(state.UserID, p.level)
	case messages.LegTome:
		p.increment(&cbStat.LegTome)
	case messages.AncientShard:
		p.increment(&cbStat.AncientShard)
	case messages.VoidShard:
		p.increment(&cbStat.VoidShard)
	case messages.SacredShard:
		p.increment(&cbStat.SacredShard)
	case messages.EpicTome:
		p.increment(&cbStat.EpicTome)
	default:
		return state, nil, UnknownResuest
	}

	p.updateStats(cbStat)

	resp := chatutils.EditTo(msg, format(cbStat), keyboards.ChooseAddDropInlineKeyboard(state.Options.Levels[0]))
	return state, resp, nil

}

func (p *CbProcessor) CancelFor(userID int64) {
	p.deleteUserStat(userID)
}

func (p *CbProcessor) deleteUserStat(userID int64) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.stats, userID)
}

func (p *CbProcessor) updateStats(stat entities.UserCbStat) {
	p.m.Lock()
	defer p.m.Unlock()
	p.stats[stat.UserID] = stat
}
func (p *CbProcessor) getOrCreateStats(userID int64) entities.UserCbStat {
	p.m.Lock()
	defer p.m.Unlock()
	if s, ok := p.stats[userID]; ok && !s.Expired() {
		return s
	}
	s := entities.NewCbUserState(userID, p.level)
	p.stats[userID] = s
	return s
}

func (p *CbProcessor) increment(val *int) {
	*val = *val + 1
}
