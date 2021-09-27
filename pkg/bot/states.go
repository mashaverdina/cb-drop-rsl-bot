package rslbot

import (
	"time"
)

type State string

const (
	MainMenu = "main-menu"
	Cb6      = "cb-6"
	Cb5      = "cb-5"
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
		State:      MainMenu,
	}
}

type CbUserState struct {
	UserID     int64
	LastUpdate time.Time

	AncientShard int `json:"AncientShard"`
	VoidShard    int `json:"VoidShard"`
	SacredShard  int `json:"SacredShard"`
	EpicTome     int `json:"EpicTome"`
	LegTome      int `json:"LegTome"`
}

func NewCbUserState(userID int64) CbUserState {
	return CbUserState{
		UserID:     userID,
		LastUpdate: time.Now(),
	}
}
