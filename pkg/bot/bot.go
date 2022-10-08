package rslbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"vkokarev.com/rslbot/pkg/globalstat"

	"vkokarev.com/rslbot/pkg/bot/command"
	"vkokarev.com/rslbot/pkg/bot/processor"
	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/notification"
	"vkokarev.com/rslbot/pkg/pg"
	"vkokarev.com/rslbot/pkg/storage"
)

type Bot struct {
	ctx                 context.Context
	botAPI              *tgbotapi.BotAPI
	states              map[int64]entities.UserState
	processors          map[entities.ProcType]processor.Processor
	commands            []command.BotCommand
	cancel              context.CancelFunc
	updates             tgbotapi.UpdatesChannel
	userStorage         *storage.UserStorage
	msgQueue            chan []tgbotapi.Chattable
	done                chan interface{}
	m                   sync.Mutex
	started             bool
	stopped             bool
	numWorkers          uint64
	notificationManager *notification.NotificationManager
	globalStatManager   *globalstat.GlobalStatManager
}

func NewBot(botAPI *tgbotapi.BotAPI, pg *pg.PGClient, numWorkers uint64) *Bot {
	bot := &Bot{
		ctx:        nil,
		botAPI:     botAPI,
		states:     make(map[int64]entities.UserState),
		processors: make(map[entities.ProcType]processor.Processor),
		commands:   make([]command.BotCommand, 0),
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
	bot.globalStatManager = globalstat.NewGlobalStatManager(cbStatStorage)

	bot.processors = map[entities.ProcType]processor.Processor{
		entities.StateMainMenu: &processor.MainProcessor{},
		entities.StateCb4:      processor.NewCbProcessor(4, cbStatStorage),
		entities.StateCb5:      processor.NewCbProcessor(5, cbStatStorage),
		entities.StateCb6:      processor.NewCbProcessor(6, cbStatStorage),
		entities.StateStats:    processor.NewStatsProcessor(cbStatStorage),
		entities.StateClans:    processor.NewClansProcessor(cbStatStorage),
		entities.StateMonth:    processor.NewMonthProcessor(cbStatStorage, bot.globalStatManager),
	}
	bot.userStorage = storage.NewUserStorage(pg)

	notificationStorage := storage.NewNotificationStorage(pg)
	bot.notificationManager = notification.NewNotificationManager(bot.msgQueue, notificationStorage, cbStatStorage)

	// order is important
	bot.commands = []command.BotCommand{
		&command.StartCommand{},
		&command.HelpCommand{},
		&command.FAQCommand{},
		&command.SupportCommand{},
		command.NewClanCommand(bot.userStorage),
		command.NewNotificationCommand(notificationStorage),
		command.NewNotifyAllCommand(bot.userStorage, bot.msgQueue),
		command.NewMigrateCommand(notificationStorage, bot.notificationManager),
		&command.NotFoundCommand{},
	}

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

	if err := b.notificationManager.Start(b.ctx); err != nil {
		return err
	}

	if err := b.globalStatManager.Start(b.ctx); err != nil {
		return err
	}
	b.started = true

	return nil
}

func (b *Bot) Stop(timeout time.Duration) error {
	b.m.Lock()
	defer b.m.Unlock()

	if !b.started || b.stopped {
		return errors.New("invalid state")
	}

	// todo
	_ = b.notificationManager.Stop()

	_ = b.globalStatManager.Stop()

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

	proc, err := b.findProcessor(state.ProcType)
	if err != nil {
		return nil, err
	}

	newState, response, err := proc.Handle(b.ctx, state, msg)
	if err == processor.UnknownResuest {
		mainProc := b.processors[entities.StateMainMenu].(*processor.MainProcessor)

		if proc == mainProc {
			return chatutils.TextTo(msg, "Не тыкай куда попало, используй кнопки!", keyboards.MainMenuKeyboard), nil
		}

		if mainProc.CanHandle(msg) {
			proc.CancelFor(msg.User.UserID)
			b.updateState(msg.User.UserID, entities.NewUserState(msg.User.UserID))
			return b.processUserMessage(msg)
		} else {
			return chatutils.TextTo(msg, "Не тыкай куда попало, используй кнопки!", keyboards.MainMenuKeyboard), nil
		}
	} else if err != nil {
		return nil, err
	}
	b.updateState(userID, newState)
	return response, nil
}

func (b *Bot) findProcessor(state entities.ProcType) (processor.Processor, error) {
	if p, ok := b.processors[state]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}

func (b *Bot) processCommand(user entities.User, commandName string, arguments string) ([]tgbotapi.Chattable, error) {
	for _, cmd := range b.commands {
		if cmd.CanHandle(commandName) {
			return cmd.Handle(b.ctx, user, commandName, arguments)
		}
	}
	return (&command.NotFoundCommand{}).Handle(b.ctx, user, commandName, arguments)
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
			if err := b.notificationManager.AssignDefaultNotifications(user.UserID); err != nil {
				log.Printf(fmt.Sprintf("failed to assign default notifications: %v", err))
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

	var msgs []tgbotapi.Chattable
	var err error
	if update.Message != nil && update.Message.Command() != "" {
		msgs, err = b.processCommand(pm.User, update.Message.Command(), update.Message.CommandArguments())
	} else {
		msgs, err = b.processUserMessage(&pm)
	}

	if err != nil {
		log.Printf("got error, while processing message: %v", err)
		msgs = chatutils.TextTo(&pm, "Произошла отвратительная ошибка, попробуй еще раз", keyboards.MainMenuKeyboard)
	}

	for _, msg := range msgs {
		if _, err := b.botAPI.Send(msg); err != nil {
			log.Printf("error while sending message: %v\n", err)
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
