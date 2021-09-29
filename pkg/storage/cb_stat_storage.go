package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.yandex/hasql"
	"gorm.io/gorm"

	"vkokarev.com/rslbot/pkg/entities"
	pg2 "vkokarev.com/rslbot/pkg/pg"
)

var NotFound = errors.New("missing entity")

type CbStatStorage struct {
	pg *pg2.PGClient
}

func NewCbStatStorage(pg *pg2.PGClient) *CbStatStorage {
	return &CbStatStorage{
		pg: pg,
	}
}

func (s *CbStatStorage) Create(ctx context.Context, state *entities.UserCbStat) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Create(state).Error
	})
}

func (s *CbStatStorage) Save(ctx context.Context, state *entities.UserCbStat) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(state).Error
	})
}

func (s *CbStatStorage) Load(ctx context.Context, userID int64, relatedTo time.Time, level int) (entities.UserCbStat, error) {
	state := &entities.UserCbStat{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.First(state, "user_id = ? and related_to = ? and level = ?", userID, relatedTo, level).Error
	})
	if err != nil {
		return entities.UserCbStat{}, err
	}
	return *state, nil
}

func (s *CbStatStorage) LastResource(ctx context.Context, userID int64, level int, resource string) (*time.Time, error) {
	var result *sql.NullTime = new(sql.NullTime)
	query := fmt.Sprintf("select max(related_to) as related_to from user_cb_stats where  user_id = ? and level = ? and %s > 0", resource)
	if err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw(query, userID, level).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(result); err != nil {
				// continue
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, nil
	}

	return &result.Time, nil
}

func (s *CbStatStorage) UserStatCombined(ctx context.Context, userID int64, levels []int, from time.Time, to time.Time) (entities.UserCbStat, error) {
	state := &entities.UserCbStat{UserID: userID}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw("select sum(ancient_shard) as ancient_shard, sum(void_shard) as void_shard, sum(sacred_shard) as sacred_shard, sum(epic_tome) as epic_tome, sum(leg_tome) as leg_tome from user_cb_stats where user_id = ? and level in ? and related_to >= ? and related_to <= ?", userID, levels, from, to).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			return rows.Scan(&state.AncientShard, &state.VoidShard, &state.SacredShard, &state.EpicTome, &state.LegTome)
		}
		return nil
	})
	if err != nil {
		return entities.UserCbStat{}, err
	}
	return *state, nil
}

func (s *CbStatStorage) UserStat(ctx context.Context, userID int64, levels []int, from time.Time, to time.Time) ([]entities.UserCbStat, error) {
	states := make([]entities.UserCbStat, 0)
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Model(&entities.UserCbStat{}).
			Select("*").
			Where("user_id = ? and level in ? and related_to >= ? and related_to <= ?", userID, levels, from, to).
			Scan(&states).
			Order("related_to").
			Error
	})

	return states, err
}
