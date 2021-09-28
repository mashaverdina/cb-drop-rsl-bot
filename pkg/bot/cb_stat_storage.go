package rslbot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.yandex/hasql"
	"gorm.io/gorm"

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

func (s *CbStatStorage) Create(ctx context.Context, state *CbUserState) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Create(state).Error
	})
}

func (s *CbStatStorage) Save(ctx context.Context, state *CbUserState) error {
	return s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.Save(state).Error
	})
}

func (s *CbStatStorage) Load(ctx context.Context, userID int64, relatedTo time.Time, level int) (CbUserState, error) {
	state := &CbUserState{}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		return db.First(state, "user_id = ? and related_to = ? and level = ?", userID, relatedTo, level).Error
	})
	if err != nil {
		return CbUserState{}, err
	}
	return *state, nil
}

func (s *CbStatStorage) LastResource(ctx context.Context, userID int64, level int, resource string) (*time.Time, error) {
	var result *sql.NullTime = new(sql.NullTime)
	query := fmt.Sprintf("select max(related_to) as related_to from cb_user_states where  user_id = ? and level = ? and %s > 0", resource)
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

func (s *CbStatStorage) UserStat(ctx context.Context, userID int64, levels []int, from time.Time, to time.Time) (CbUserState, error) {
	state := &CbUserState{UserID: userID}
	err := s.pg.ExecuteInTransaction(hasql.Primary, func(db *gorm.DB) error {
		rows, err := db.Raw("select sum(ancient_shard) as ancient_shard, sum(void_shard) as void_shard, sum(sacred_shard) as sacred_shard, sum(epic_tome) as epic_tome, sum(leg_tome) as leg_tome from cb_user_states where user_id = ? and level in ? and related_to >= ? and related_to <= ?", userID, levels, from, to).Rows()
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
		return CbUserState{}, err
	}
	return *state, nil
}
