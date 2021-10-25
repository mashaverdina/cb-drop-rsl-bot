package entities

import (
	"time"

	"vkokarev.com/rslbot/pkg/utils"
)

type State string

const (
	StateMainMenu = "main-menu"
	StateCb6      = "cb-6"
	StateCb5      = "cb-5"
	StateCb4      = "cb-4"
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
		LastUpdate: utils.MskNow(),
		State:      StateMainMenu,
	}
}
