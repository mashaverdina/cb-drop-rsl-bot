package storage

import (
	"context"

	"golang.yandex/hasql"
	"gorm.io/gorm"

	"vkokarev.com/rslbot/pkg/entities"
	pg2 "vkokarev.com/rslbot/pkg/pg"
)

type UserStorage struct {
	pg *pg2.PGClient
}

func NewUserStorage(pg *pg2.PGClient) *UserStorage {
	return &UserStorage{
		pg: pg,
	}
}

func (s *UserStorage) Create(ctx context.Context, user *entities.User) (entities.User, error) {
	return *user, s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Create(user).Error
	})

}

func (s *UserStorage) Save(ctx context.Context, user *entities.User) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

func (s *UserStorage) Load(ctx context.Context, userID int64) (entities.User, error) {
	user := &entities.User{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.First(user, "user_id = ?", userID).Error
	})
	if err != nil {
		return entities.User{}, err
	}
	return *user, nil
}

func (s *UserStorage) All(ctx context.Context) ([]entities.User, error) {
	users := make([]entities.User, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Model(&entities.User{}).Select("*").Scan(&users).Error
	})
	return users, err
}
