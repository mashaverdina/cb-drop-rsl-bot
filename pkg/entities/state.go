package entities

import (
	"time"
	"vkokarev.com/rslbot/pkg/utils"
)

type ProcType string

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
	ProcType   ProcType
	Options    Options
}

type Options struct {
	Levels []int
}

func (o *Options) Clear() {
	*o = Options{}
}

func (o *Options) DropLevels() {
	o.Levels = []int{}
}

func (o *Options) WithLevels(levels ...int) {
	o.Levels = levels
}

func NewUserState(userID int64) UserState {
	return UserState{
		UserID:     userID,
		LastUpdate: utils.MskNow(),
		ProcType:   StateMainMenu,
		Options:    Options{},
	}
}
