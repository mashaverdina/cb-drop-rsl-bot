package rslbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/pg"
)

type Bot struct {
	ctx         context.Context
	botAPI      *tgbotapi.BotAPI
	states      map[int64]UserState
	processors  map[State]Processor
	cancel      context.CancelFunc
	updates     tgbotapi.UpdatesChannel
	userStorage *UserStorage
	done        chan interface{}
	m           sync.Mutex
	started     bool
	stopped     bool
}

func NewBot(botAPI *tgbotapi.BotAPI, pg *pg.PGClient) *Bot {
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

	cbStatStorage := NewCbStatStorage(pg)
	bot.processors = map[State]Processor{
		StateMainMenu: &MainProcessor{},
		StateCb5:      NewCbProcessor(5, cbStatStorage),
		StateCb6:      NewCbProcessor(6, cbStatStorage),
		StateStats:    NewStatsProcessor(cbStatStorage),
		StateMonth:    NewMonthProcessor(cbStatStorage),
	}
	bot.userStorage = NewUserStorage(pg)

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
	for {
		select {
		case <-b.ctx.Done():
			b.done <- true
			return
		case update := <-updates:
			pm := ProcessingMessage{}
			if update.Message != nil {
				user, err := b.userStorage.Load(b.ctx, update.Message.Chat.ID)
				if err != nil {
					user, err = b.userStorage.Create(b.ctx, &User{
						UserID:       update.Message.Chat.ID,
						FirstName:    update.Message.From.FirstName,
						LastName:     update.Message.From.LastName,
						UserName:     update.Message.From.UserName,
						LanguageCode: update.Message.From.LanguageCode,
						Clan:         "",
						Nickname:     "",
					})
					if err != nil {
						user = User{UserID: update.Message.Chat.ID}
					}
				}
				pm = ProcessingMessage{
					User:      user,
					ChatID:    update.Message.Chat.ID,
					Text:      update.Message.Text,
					MessageID: update.Message.MessageID,
				}
			} else if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := b.botAPI.Request(callback); err != nil {
					panic(err)
				}
				user, err := b.userStorage.Load(b.ctx, update.CallbackQuery.Message.Chat.ID)
				if err != nil {
					user = User{UserID: update.CallbackQuery.Message.Chat.ID}
				}
				pm = ProcessingMessage{
					User:      user,
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					Text:      callback.Text,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
			}

			if update.Message != nil && update.Message.Command() != "" {
				b.processCommand(pm.User, update.Message.Command(), update.Message.CommandArguments())
			} else {
				msgs, err := b.processUserMessage(&pm)
				if err != nil {
					log.Fatalf("got error, while processing message: %v", err)
				}
				for _, msg := range msgs {
					if _, err := b.botAPI.Send(msg); err != nil {
						log.Printf("error while sending message: %v\n", err)
					}
				}
			}
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

func (b *Bot) processUserMessage(msg *ProcessingMessage) ([]tgbotapi.Chattable, error) {
	userID := msg.User.UserID
	state := b.getOrCreateState(userID)

	processor, err := b.findProcessor(state.State)
	if err != nil {
		return nil, err
	}
	newState, response, err := processor.Handle(b.ctx, state, msg)
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

func (b *Bot) processCommand(user User, command string, arguments string) {
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(user.UserID, "Добро пожаловать в RSL.CB бот. Используй клавиатуру внизу")
		msg.ReplyMarkup = keyboards.MainMenuKeyboard
		_, _ = b.botAPI.Send(msg)
	case "notifyall":
		if !user.HasSudo {
			b.NotifySudo(user)
			return
		}
		if err := b.NotifyAll(b.ctx, arguments); err != nil {
			msg := tgbotapi.NewMessage(user.UserID, fmt.Sprintf("Ошибка: %v", err))
			_, _ = b.botAPI.Send(msg)
		}
	default:
		b.NotifyNotFound(user)
	}
}

func (b *Bot) NotifySudo(user User) {
	msg := tgbotapi.NewMessage(user.UserID, "Для данной команды требуются супер права")
	_, _ = b.botAPI.Send(msg)
}

func (b *Bot) NotifyNotFound(user User) {
	msg := tgbotapi.NewMessage(user.UserID, "Команда не найдена")
	_, _ = b.botAPI.Send(msg)
}

func (b *Bot) NotifyAll(ctx context.Context, arguments string) error {
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	users := make(chan User)

	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case u := <-users:
					_, _ = b.botAPI.Send(tgbotapi.NewMessage(u.UserID, arguments))
				case <-workerCtx.Done():
					return
				}
			}
		}()
	}

	allUsers, err := b.userStorage.All(b.ctx)
	if err != nil {
		return err
	}
	for _, user := range allUsers {
		select {
		case users <- user:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
