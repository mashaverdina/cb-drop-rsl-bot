package storage

import (
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
		return db.Model(&entities.Notification{}).Select("*").Scan(&notifications).Error
	})
	return notifications, err
}

func (s *NotificationStorage) GetUsersFor(notification entities.Notification) ([]int64, error) {
	user_ids := make([]int64, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw("select u.user_id from users as u where (u.user_id not in (select user_id from disabled_notifications where notification_id=?))", notification.ID).Rows()
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
			user_ids = append(user_ids, uid)
		}
		return nil
	})
	return user_ids, err
}

func (s *NotificationStorage) DisableNotification(disabled entities.DisabledNotifications) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(disabled).Error
	})
}

func (s *NotificationStorage) EnableNotification(disabled entities.DisabledNotifications) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Delete(disabled).Error
	})
}

func (s *NotificationStorage) Update(notification entities.Notification) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(notification).Error
	})
}

func (s *NotificationStorage) GetByAlias(alias string) (entities.Notification, error) {
	n := entities.Notification{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Model(&entities.Notification{}).Select("*").Where("alias=?", alias).Scan(&n).Error
	})
	return n, err
}
