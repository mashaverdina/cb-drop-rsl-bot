package storage

import (
	"time"

	"golang.yandex/hasql"
	"gorm.io/gorm"

	"vkokarev.com/rslbot/pkg/entities"
	pg2 "vkokarev.com/rslbot/pkg/pg"
)

type NotificationStorage struct {
	pg *pg2.PGClient
}

func NewNotificationStorage(pg *pg2.PGClient) *NotificationStorage {
	return &NotificationStorage{
		pg: pg,
	}
}

func (s *NotificationStorage) LoadAllNotifications() ([]entities.Notification, error) {
	notifications := make([]entities.Notification, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		// return db.Model(&entities.Notification{}).Select("*").Scan(&notifications).Error
		return db.Model(&entities.Notification{}).Select("*").Joins("left join notification_fires on notification_fires.notification_id = notifications.notification_id").Scan(&notifications).Error
	})
	return notifications, err
}

func (s *NotificationStorage) GetUsersFor(notification entities.Notification) ([]int64, error) {
	userIDs := make([]int64, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw("select user_id from notification_users where fire_id=?", notification.FireID).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var uid int64 = 0
			err := rows.Scan(&uid)
			if err != nil {
				return err
			}
			userIDs = append(userIDs, uid)
		}
		return nil
	})
	return userIDs, err
}

func (s *NotificationStorage) DisableNotification(userID int64, n entities.Notification) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Exec("delete from notification_users where user_id = ? and fire_id in (select fire_id from notification_fires where notification_fires.notification_id=?)", userID, n.NotificationID).Error
	})
}

func (s *NotificationStorage) EnableNotification(userID int64, n entities.Notification) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Create(entities.NotificationUsers{
			UserID: userID,
			FireID: n.FireID,
		}).Error
	})
}

func (s *NotificationStorage) UpdateFire(notification entities.Notification) error {
	nf := entities.NotificationFires{
		FireID:       notification.FireID,
		FireAt:       notification.FireAt,
		LastFireTime: notification.LastFireTime,
	}
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(nf).Error
	})
}

func (s *NotificationStorage) GetByAlias(alias string) (entities.Notification, error) {
	n := entities.Notification{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Model(&entities.Notification{}).Select("*").Where("alias=?", alias).Scan(&n).Error
	})
	return n, err
}

func (s *NotificationStorage) LoadFire(notificationID int64, hour, minutes int) (entities.Notification, error) {
	fireAt := time.Date(2021, 1, 1, hour, minutes, 0, 0, time.UTC)
	notification := entities.Notification{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Model(&entities.Notification{}).Select("*").Joins("left join notification_fires on notification_fires.notification_id = notifications.notification_id").Where("notifications.notification_id = ? and notification_fires.fire_at = ?", notificationID, fireAt).Scan(&notification).Error
	})
	return notification, err
}

func (s *NotificationStorage) NonDisabledUsers() ([]int64, error) {
	userIDs := make([]int64, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw("select user_id from users where user_id not in (select user_id from disabled_notifications)").Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var uid int64 = 0
			err := rows.Scan(&uid)
			if err != nil {
				return err
			}
			userIDs = append(userIDs, uid)
		}
		return nil
	})
	return userIDs, err
}
