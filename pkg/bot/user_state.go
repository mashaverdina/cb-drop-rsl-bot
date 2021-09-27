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
	UserID     int
	LastUpdate time.Time
	State      State
}

func NewUserState(userID int) UserState {
	return UserState{
		UserID:     userID,
		LastUpdate: time.Now(),
		State:      MainMenu,
	}
}