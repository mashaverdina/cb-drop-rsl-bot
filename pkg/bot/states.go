package rslbot

import (
	"time"
)

type State string

const (
	StateMainMenu = "main-menu"
	StateCb6      = "cb-6"
	StateCb5      = "cb-5"
	StateStats    = "stats"
	StateMonth    = "month"
)

type UserState struct {
	UserID     int64
	LastUpdate time.Time
	State      State
}

func NewUserState(userID int64) UserState {
	return UserState{
		UserID:     userID,
		LastUpdate: time.Now(),
		State:      StateMainMenu,
	}
}

type CbUserState struct {
	UserID     int64     `gorm:"primaryKey;index:planned_index"`
	RelatedTo  time.Time `gorm:"primaryKey;index:planned_index"`
	Level      int       `gorm:"primaryKey;index:planned_index"`
	LastUpdate time.Time `gorm:"index:planned_index"`

	AncientShard int `json:"AncientShard"`
	VoidShard    int `json:"VoidShard"`
	SacredShard  int `json:"SacredShard"`
	EpicTome     int `json:"EpicTome"`
	LegTome      int `json:"LegTome"`
}

func (s CbUserState) Expired() bool {
	return Related(s.RelatedTo) != Related(time.Now())
}

func NewCbUserState(userID int64, level int) CbUserState {
	return CbUserState{
		UserID:     userID,
		Level:      level,
		LastUpdate: time.Now(),
		RelatedTo:  Related(time.Now()),
	}
}

func Related(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
