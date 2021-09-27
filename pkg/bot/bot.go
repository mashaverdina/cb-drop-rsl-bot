package rslbot

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	ctx        context.Context
	botAPI     *tgbotapi.BotAPI
	states     map[int64]UserState
	processors map[State]Processor
	cancel     context.CancelFunc
	updates    tgbotapi.UpdatesChannel
	done       chan interface{}
	m          sync.Mutex
	started    bool
	stopped    bool
}

func NewBot(botAPI *tgbotapi.BotAPI) *Bot {
	bot := &Bot{
		ctx:        nil,
		botAPI:     botAPI,
		states:     make(map[int64]UserState),
		processors: make(map[State]Processor),
		cancel:     nil,
		updates:    nil,
		done:       make(chan interface{}),
		m:          sync.Mutex{},
		started:    false,
		stopped:    false,
	}

	bot.processors = map[State]Processor{
		StateMainMenu: &MainProcessor{},
		StateCb5:      NewCbProcessor(5),
		StateCb6:      NewCbProcessor(6),
		StateStats:    NewStatsProcessor(),
		StateMonth:    NewMonthProcessor(),
	}

	bot.processors[StateMainMenu] = &MainProcessor{}
	bot.processors[StateCb5] = NewCbProcessor(5)
	bot.processors[StateCb6] = NewCbProcessor(6)
	bot.processors[StateStats] = NewStatsProcessor()
	bot.processors[StateMonth] = NewMonthProcessor()
	return bot

}

func (b *Bot) Start(ctx context.Context) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.started {
		return errors.New("started already")
	}

	b.ctx, b.cancel = context.WithCancel(ctx)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)
	go b.loop(updates)

	b.started = true

	return nil
}

func (b *Bot) Stop(timeout time.Duration) error {
	b.m.Lock()
	defer b.m.Unlock()

	if !b.started || b.stopped {
		return errors.New("invalid state")
	}

	b.cancel()
	select {
	case <-b.done:
		return nil
	case <-time.After(timeout):
		return errors.New("can't close bot due to timeout")
	}
}

func (b *Bot) loop(updates tgbotapi.UpdatesChannel) {
	log.Println("starting bot loop")
	for update := range updates {
		select {
		case <-b.ctx.Done():
			b.done <- true
			return
		default:
		}

		pm := ProcessingMessage{}

		if update.Message != nil {
			pm = ProcessingMessage{
				UserID:    update.Message.Chat.ID,
				ChatID:    update.Message.Chat.ID,
				Text:      update.Message.Text,
				MessageID: update.Message.MessageID,
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.botAPI.Request(callback); err != nil {
				panic(err)
			}
			pm = ProcessingMessage{
				UserID:    update.CallbackQuery.Message.Chat.ID,
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				Text:      callback.Text,
				MessageID: update.CallbackQuery.Message.MessageID,
			}
		}

		msg, err := b.processUserMessage(&pm)
		if err != nil {
			log.Fatalf("got error, while processing message: %v", err)
		}

		// Send the message.
		if _, err := b.botAPI.Send(msg); err != nil {
			// todo не паниковать
			log.Printf("error while sending message: %v\n", err)
		}
	}
	log.Println("exiting loop")
}

func (b *Bot) getOrCreateState(userID int64) UserState {
	if s, ok := b.states[userID]; ok {
		return s
	}
	s := NewUserState(userID)
	b.states[userID] = s
	return s
}

func (b *Bot) processUserMessage(msg *ProcessingMessage) (tgbotapi.Chattable, error) {
	userID := msg.UserID
	state := b.getOrCreateState(userID)

	processor, err := b.findProcessor(state.State)
	if err != nil {
		return nil, err
	}
	newState, response, err := processor.Handle(state, msg)
	if err != nil {
		return nil, err
	}
	b.states[userID] = newState
	return response, nil
}

func (b *Bot) findProcessor(state State) (Processor, error) {
	if p, ok := b.processors[state]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}
