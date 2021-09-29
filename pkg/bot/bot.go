package rslbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/bot/processor"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/pg"
	"vkokarev.com/rslbot/pkg/storage"
)

type Bot struct {
	ctx         context.Context
	botAPI      *tgbotapi.BotAPI
	states      map[int64]entities.UserState
	processors  map[entities.State]processor.Processor
	cancel      context.CancelFunc
	updates     tgbotapi.UpdatesChannel
	userStorage *storage.UserStorage
	msgQueue    chan []tgbotapi.Chattable
	done        chan interface{}
	m           sync.Mutex
	started     bool
	stopped     bool
	numWorkers  uint64
}

func NewBot(botAPI *tgbotapi.BotAPI, pg *pg.PGClient, numWorkers uint64) *Bot {
	bot := &Bot{
		ctx:        nil,
		botAPI:     botAPI,
		states:     make(map[int64]entities.UserState),
		processors: make(map[entities.State]processor.Processor),
		cancel:     nil,
		updates:    nil,
		msgQueue:   make(chan []tgbotapi.Chattable),
		done:       make(chan interface{}),
		m:          sync.Mutex{},
		started:    false,
		stopped:    false,
		numWorkers: numWorkers,
	}

	cbStatStorage := storage.NewCbStatStorage(pg)
	bot.processors = map[entities.State]processor.Processor{
		entities.StateMainMenu: &processor.MainProcessor{},
		entities.StateCb5:      processor.NewCbProcessor(5, cbStatStorage),
		entities.StateCb6:      processor.NewCbProcessor(6, cbStatStorage),
		entities.StateStats:    processor.NewStatsProcessor(cbStatStorage),
		entities.StateMonth:    processor.NewMonthProcessor(cbStatStorage),
	}
	bot.userStorage = storage.NewUserStorage(pg)

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

	for i := uint64(0); i < b.numWorkers; i++ {
		go b.worker(ctx)
	}

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
			go b.processUpdate(update)
		}
	}
}

func (b *Bot) getOrCreateState(userID int64) entities.UserState {
	b.m.Lock()
	defer b.m.Unlock()
	if s, ok := b.states[userID]; ok {
		return s
	}
	s := entities.NewUserState(userID)
	b.states[userID] = s
	return s
}

func (b *Bot) updateState(userID int64, newState entities.UserState) {
	b.m.Lock()
	defer b.m.Unlock()
	b.states[userID] = newState
}

func (b *Bot) processUserMessage(msg *processor.ProcessingMessage) ([]tgbotapi.Chattable, error) {
	userID := msg.User.UserID
	state := b.getOrCreateState(userID)

	proc, err := b.findProcessor(state.State)
	if err != nil {
		return nil, err
	}
	newState, response, err := proc.Handle(b.ctx, state, msg)
	if err != nil {
		return nil, err
	}
	b.updateState(userID, newState)
	return response, nil
}

func (b *Bot) findProcessor(state entities.State) (processor.Processor, error) {
	if p, ok := b.processors[state]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}

func (b *Bot) processCommand(user entities.User, command string, arguments string) {
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

func (b *Bot) NotifySudo(user entities.User) {
	msg := tgbotapi.NewMessage(user.UserID, "Для данной команды требуются супер права")
	_, _ = b.botAPI.Send(msg)
}

func (b *Bot) NotifyNotFound(user entities.User) {
	msg := tgbotapi.NewMessage(user.UserID, "Команда не найдена")
	_, _ = b.botAPI.Send(msg)
}

func (b *Bot) NotifyAll(ctx context.Context, arguments string) error {
	allUsers, err := b.userStorage.All(b.ctx)
	if err != nil {
		return err
	}
	for _, user := range allUsers {
		select {
		case b.msgQueue <- []tgbotapi.Chattable{tgbotapi.NewMessage(user.UserID, arguments)}:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func (b *Bot) processUpdate(update tgbotapi.Update) {
	pm := processor.ProcessingMessage{}
	if update.Message != nil {
		user, err := b.userStorage.Load(b.ctx, update.Message.Chat.ID)
		if err != nil {
			user, err = b.userStorage.Create(b.ctx, &entities.User{
				UserID:       update.Message.Chat.ID,
				FirstName:    update.Message.From.FirstName,
				LastName:     update.Message.From.LastName,
				UserName:     update.Message.From.UserName,
				LanguageCode: update.Message.From.LanguageCode,
				Clan:         "",
				Nickname:     "",
			})
			if err != nil {
				user = entities.User{UserID: update.Message.Chat.ID}
			}
		}
		pm = processor.ProcessingMessage{
			User:    user,
			ChatID:  update.Message.Chat.ID,
			Text:    update.Message.Text,
			Message: update.Message,
		}
	} else if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := b.botAPI.Request(callback); err != nil {
			log.Print(fmt.Errorf("got error while processint callback request", err))
		}
		user, err := b.userStorage.Load(b.ctx, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			user = entities.User{UserID: update.CallbackQuery.Message.Chat.ID}
		}
		pm = processor.ProcessingMessage{
			User:    user,
			ChatID:  update.CallbackQuery.Message.Chat.ID,
			Text:    callback.Text,
			Message: update.CallbackQuery.Message,
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

func (b *Bot) worker(ctx context.Context) {
	for {
		select {
		case msgs := <-b.msgQueue:
			for _, msg := range msgs {
				_, err := b.botAPI.Send(msg)
				if err != nil {
					log.Println(fmt.Errorf("error while sending msg", err).Error())
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
