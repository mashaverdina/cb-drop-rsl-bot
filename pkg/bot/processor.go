package rslbot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/keyboards"
)

var monthMap = map[string]time.Month{
	keyboards.Jan: time.January,
	keyboards.Feb: time.February,
	keyboards.Mar: time.March,
	keyboards.Apr: time.April,
	keyboards.May: time.May,
	keyboards.Jun: time.June,
	keyboards.Jul: time.July,
	keyboards.Aug: time.August,
	keyboards.Sep: time.September,
	keyboards.Oct: time.October,
	keyboards.Nov: time.November,
	keyboards.Dec: time.December,
}

const (
	dateFormat = "02.01.2006"
)

type ProcessingMessage struct {
	User      User
	ChatID    int64
	MessageID int
	Text      string
}

type Processor interface {
	Handle(ctx context.Context, state UserState, msg *ProcessingMessage) (UserState, []tgbotapi.Chattable, error)
}

type MainProcessor struct {
}

func editTo(chatID int64, msgID int, text string, markup *tgbotapi.InlineKeyboardMarkup) []tgbotapi.Chattable {
	resp := tgbotapi.NewEditMessageText(chatID, msgID, text)
	if markup != nil {
		resp.ReplyMarkup = markup
	}
	resp.ParseMode = "Markdown"
	return []tgbotapi.Chattable{resp}
}

func textTo(chatID int64, text string, markup interface{}) []tgbotapi.Chattable {
	resp := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		resp.ReplyMarkup = markup
	}
	return []tgbotapi.Chattable{resp}
}

func (p *MainProcessor) Handle(ctx context.Context, state UserState, msg *ProcessingMessage) (UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case keyboards.Cb5:
		state.State = StateCb5
		resp := textTo(msg.ChatID, "Что упало с 5го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case keyboards.Cb6:
		state.State = StateCb6
		resp := textTo(msg.ChatID, "Что упало с 6го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case keyboards.Stats:
		state.State = StateStats
		resp := textTo(msg.ChatID, "Что тебе показать?", keyboards.StatsKeyboard)
		return state, resp, nil
	}

	resp := textTo(msg.ChatID, "Привет", keyboards.MainMenuKeyboard)
	return state, resp, nil
}

type CbProcessor struct {
	level   int
	stats   map[int64]CbUserState
	storage *CbStatStorage
}

func NewCbProcessor(level int, storage *CbStatStorage) *CbProcessor {
	return &CbProcessor{
		level:   level,
		stats:   make(map[int64]CbUserState),
		storage: storage,
	}
}

func (p *CbProcessor) Handle(ctx context.Context, state UserState, msg *ProcessingMessage) (UserState, []tgbotapi.Chattable, error) {
	cbState := p.getOrCreateStats(state.UserID)
	switch msg.Text {
	case keyboards.Reject:
		state.State = StateMainMenu
		delete(p.stats, state.UserID)
		resp := joinResp(
			[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.ChatID, msg.MessageID)},
			textTo(msg.ChatID, "До встречи", keyboards.MainMenuKeyboard),
		)
		return state, resp, nil
	case keyboards.Approve:
		state.State = StateMainMenu

		cbState := p.stats[state.UserID]
		err := p.storage.Save(ctx, &cbState)
		if err != nil {
			return UserState{}, nil, err
		}

		p.stats[state.UserID] = NewCbUserState(state.UserID, p.level)
		resp := joinResp(
			editTo(msg.ChatID, msg.MessageID, msgFromStat(cbState, p.level), nil),
			textTo(msg.ChatID, "Записано", keyboards.MainMenuKeyboard),
		)
		return state, resp, nil
	case keyboards.Clear:
		cbState = NewCbUserState(state.UserID, p.level)
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
		resp := textTo(msg.ChatID, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)", nil)
		return state, resp, nil
	}

	p.stats[state.UserID] = cbState

	resp := editTo(msg.ChatID, msg.MessageID, msgFromStat(cbState, p.level), &keyboards.AddDropInlineKeyboard)
	return state, resp, nil

}

func joinResp(resps ...[]tgbotapi.Chattable) []tgbotapi.Chattable {
	result := make([]tgbotapi.Chattable, 0)
	for _, arr := range resps {
		result = append(result, arr...)
	}
	return result
}

func msgFromStat(state CbUserState, level int) string {
	lines := []string{}
	if level > 0 {
		lines = append(lines, fmt.Sprintf("Стата по *%d КБ*", level))
	}

	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.AncientShard, state.AncientShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.VoidShard, state.VoidShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.SacredShard, state.SacredShard))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.EpicTome, state.EpicTome))
	lines = append(lines, fmt.Sprintf("%s -- %d", keyboards.LegTome, state.LegTome))

	return strings.Join(lines, "\n")
}

func (p *CbProcessor) getOrCreateStats(userID int64) CbUserState {
	if s, ok := p.stats[userID]; ok && !s.Expired() {
		return s
	}
	s := NewCbUserState(userID, p.level)
	p.stats[userID] = s
	return s
}

func (p *CbProcessor) increment(val *int) {
	*val = *val + 1
}

type StatsProcessor struct {
	cbStatStorage *CbStatStorage
}

func NewStatsProcessor(cbStatStorage *CbStatStorage) *StatsProcessor {
	return &StatsProcessor{
		cbStatStorage: cbStatStorage,
	}
}

func (p *StatsProcessor) LastStat(ctx context.Context, msg *ProcessingMessage, resource string, header string) ([]tgbotapi.Chattable, error) {
	lastFrom5, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 5, resource)
	if err != nil {
		return nil, err
	}
	lastFrom6, err := p.cbStatStorage.LastResource(ctx, msg.User.UserID, 6, resource)
	if err != nil {
		return nil, err
	}

	resp := textTo(msg.ChatID, strings.Join([]string{
		header,
		fmt.Sprintf("С 5го -- %s", timePast(lastFrom5)),
		fmt.Sprintf("С 6го -- %s", timePast(lastFrom6)),
	}, "\n"), keyboards.MainMenuKeyboard)
	return resp, nil
}

func (p *StatsProcessor) Handle(ctx context.Context, state UserState, msg *ProcessingMessage) (UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case keyboards.Back:
		state.State = StateMainMenu
		resp := textTo(msg.ChatID, "До встречи", keyboards.MainMenuKeyboard)
		return state, resp, nil
	case keyboards.LastVoidShard:
		state.State = StateMainMenu
		resp, err := p.LastStat(ctx, msg, "void_shard", keyboards.LastVoidShard+" осколок")
		if err != nil {
			return UserState{}, nil, err
		}
		return state, joinResp(
			[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.ChatID, msg.MessageID)},
			resp,
		), nil
	case keyboards.LastSacredShard:
		state.State = StateMainMenu
		resp, err := p.LastStat(ctx, msg, "sacred_shard", keyboards.LastVoidShard+" осколок")
		if err != nil {
			return UserState{}, nil, err
		}
		return state, joinResp(
			[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.ChatID, msg.MessageID)},
			resp,
		), nil
	case keyboards.LastLegTome:
		state.State = StateMainMenu
		resp, err := p.LastStat(ctx, msg, "leg_tome", keyboards.LastLegTome)
		if err != nil {
			return UserState{}, nil, err
		}
		return state, joinResp(
			[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.ChatID, msg.MessageID)},
			resp,
		), nil
	case keyboards.MonthStats:
		state.State = StateMonth
		return state, joinResp(
			[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.ChatID, msg.MessageID)},
			textTo(msg.ChatID, "📅 Выбери месяц", keyboards.ChooseMonthKeyboard),
		), nil
	default:
		resp := textTo(msg.ChatID, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)", nil)
		return state, resp, nil
	}
}

type MonthProcessor struct {
	cbStatStorage *CbStatStorage
}

func NewMonthProcessor(cbStatStorage *CbStatStorage) *MonthProcessor {
	return &MonthProcessor{
		cbStatStorage: cbStatStorage,
	}
}

func (p *MonthProcessor) Handle(ctx context.Context, state UserState, msg *ProcessingMessage) (UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case keyboards.Back:
		state.State = StateStats
		resp := textTo(msg.ChatID, "Что тебе показать?", keyboards.StatsKeyboard)
		return state, resp, nil
	case keyboards.Jan, keyboards.Feb, keyboards.Mar, keyboards.Apr, keyboards.May, keyboards.Jun, keyboards.Jul, keyboards.Aug, keyboards.Sep, keyboards.Oct, keyboards.Nov, keyboards.Dec:
		state.State = StateMainMenu
		// TODO: get stats for month from DB
		from, to := mothInterval(msg.Text)
		monthStat, err := p.cbStatStorage.UserStatCombined(ctx, msg.User.UserID, []int{5, 6}, from, to)
		var replyMsg = ""
		if err != nil {
			replyMsg = "Статистики пока нет"
		} else {
			statText := msgFromStat(monthStat, 0)
			replyMsg = fmt.Sprintf("Вот твоя статистика за %s:\n%s", msg.Text, statText)
		}

		resp := textTo(msg.ChatID, replyMsg, keyboards.MainMenuKeyboard)
		return state, resp, nil
	default:
		resp := textTo(msg.ChatID, "АХАХАХХАА ТЫТ ТУТ ЗАВИС (Нажми закрыть)", nil)
		return state, resp, nil
	}
}

func timePast(t *time.Time) string {
	if t == nil {
		return "никогда"
	}
	delta := time.Now().Sub(*t)
	if delta.Hours() < 24 {
		return t.Format(dateFormat) + " (сегодня)"
	}
	return t.Format(dateFormat) + " (" + strconv.FormatInt(int64(delta.Hours()/24), 10) + " д. назад)"
}

func mothInterval(month string) (time.Time, time.Time) {
	mn := monthMap[month]
	cy, cm, _ := time.Now().Date()
	if cm < mn {
		cy = cy - 1
	}

	from := time.Date(cy, mn, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return from, to
}
