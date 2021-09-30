package entities

import (
	"time"

	"vkokarev.com/rslbot/pkg/utils"
)

type UserCbStat struct {
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

func (s UserCbStat) Expired() bool {
	return Related(s.RelatedTo) != Related(utils.MskNow())
}

func NewCbUserState(userID int64, level int) UserCbStat {
	return UserCbStat{
		UserID:     userID,
		Level:      level,
		LastUpdate: utils.MskNow(),
		RelatedTo:  Related(utils.MskNow()),
	}
}

func Related(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
