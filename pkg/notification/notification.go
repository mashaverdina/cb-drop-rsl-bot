package notification

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/bot/command"
	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/storage"
	"vkokarev.com/rslbot/pkg/utils"
)

type NotificationManager struct {
	msgQueue             chan<- []tgbotapi.Chattable
	notificationStorage  *storage.NotificationStorage
	cbStatStorage        *storage.CbStatStorage
	m                    sync.Mutex
	started              bool
	cancel               context.CancelFunc
	ctx                  context.Context
	defaultNotifications []entities.Notification
}

func NewNotificationManager(msgQueue chan<- []tgbotapi.Chattable, notificationStorage *storage.NotificationStorage, cbStatStorage *storage.CbStatStorage) *NotificationManager {
	return &NotificationManager{
		msgQueue:            msgQueue,
		notificationStorage: notificationStorage,
		cbStatStorage:       cbStatStorage,
		m:                   sync.Mutex{},
		started:             false,
		cancel:              nil,
		ctx:                 nil,
	}
}

func (nm *NotificationManager) Start(ctx context.Context) error {
	nm.m.Lock()
	defer nm.m.Unlock()
	if nm.started {
		return errors.New("started already")
	}
	nm.ctx, nm.cancel = context.WithCancel(ctx)

	defaultNotifications := make([]entities.Notification, 0)
	fillDropNotification, err := nm.notificationStorage.GetByAlias("fill_drop")
	if err != nil {
		return err
	}
	fillDropNotification, err = nm.notificationStorage.LoadFire(fillDropNotification.NotificationID, 13, 7)
	if err != nil {
		return err
	}
	nm.defaultNotifications = append(defaultNotifications, fillDropNotification)

	go nm.loop()
	nm.started = true
	return nil
}

func (nm *NotificationManager) Stop() error {
	nm.m.Lock()
	defer nm.m.Unlock()
	if !nm.started {
		return errors.New("not started")
	}
	nm.cancel()
	// todo done chan
	return nil
}

func (nm *NotificationManager) loop() {
	// ticker := time.NewTicker(time.Minute)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			ns, err := nm.notificationStorage.LoadAllNotifications()
			if err != nil {
				log.Println(fmt.Sprintf("error, while getting notifications: %v", err))
				continue
			}
			for _, notification := range ns {
				if notification.ShouldBeStarted() {
					nm.fireNotification(notification)
				}
			}
		case <-nm.ctx.Done():
			return
		}
	}
}

func (nm *NotificationManager) fireNotification(notification entities.Notification) {
	log.Println("notification", notification.Alias, "started")

	notification.LastFireTime = utils.MskNow()
	err := nm.notificationStorage.UpdateFire(notification)
	if err != nil {
		log.Println(fmt.Sprintf("error, can't set fire on notification: %v", err))
		return
	}

	users, err := nm.notificationStorage.GetUsersFor(notification)
	if err != nil {
		log.Println(fmt.Sprintf("error, can't load notification's users: %v", err))
		return
	}

	// todo fix that hack!
	if notification.RemoveActiveUsers {
		users, err = nm.removeActiveUsers(users)
		if err != nil {
			log.Println(fmt.Sprintf("error, can't remove active users: %v", err))
			return
		}
	}

	for _, user := range users {
		select {
		case nm.msgQueue <- chatutils.TextToNoMarkdown(
			&chatutils.SimpleMessage{user},
			fmt.Sprintf("%s\nЧто бы больше не получать данное уведомление введи (нажми) /%s%s\nДля изменения времени скопируй: /%s%s 13:30", notification.Text, command.NotificationOff, notification.Alias, command.NotificationOff, notification.Alias),
			keyboards.MainMenuKeyboard):
		case <-nm.ctx.Done():
			log.Println(fmt.Sprintf("notification %s was canceled", notification.Alias))
			return
		}
	}
	log.Println("notification", notification.Alias, "finished")
}

func (nm *NotificationManager) removeActiveUsers(users []int64) ([]int64, error) {
	active, err := nm.cbStatStorage.ActiveUsersAt(utils.MskNow())
	if err != nil {
		return nil, err
	}
	activeMap := make(map[int64]bool)
	for _, uid := range active {
		activeMap[uid] = true
	}

	result := make([]int64, 0, len(users))
	for _, uid := range users {
		if !activeMap[uid] {
			result = append(result, uid)
		}
	}
	return result, nil
}

func (nm *NotificationManager) AssignDefaultNotifications(user entities.User) error {
	for _, n := range nm.defaultNotifications {
		if err := nm.notificationStorage.DisableNotification(user, n); err != nil {
			return err
		}
		if err := nm.notificationStorage.EnableNotification(user, n); err != nil {
			return err
		}
	}
	return nil
}
