package entities

import (
	"time"
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
	return Related(s.RelatedTo) != Related(time.Now())
}

func NewCbUserState(userID int64, level int) UserCbStat {
	return UserCbStat{
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
