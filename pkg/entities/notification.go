package entities

import (
	"time"

	"vkokarev.com/rslbot/pkg/utils"
)

type Notification struct {
	ID           int64 `gorm:"primaryKey;index:planned_index;autoIncrement:true"`
	Alias        string
	Text         string
	FireAt       time.Time
	LastFireTime time.Time
}

type DisabledNotifications struct {
	UserID         int64 `gorm:"primaryKey"`
	NotificationID int64 `gorm:"primaryKey"`
}

func (n *Notification) ShouldBeStarted() bool {
	now := utils.MskNow()
	nowYear, nowMonth, nowDay := now.Date()

	fireAt := time.Date(nowYear, nowMonth, nowDay, n.FireAt.Hour(), n.FireAt.Minute(), n.FireAt.Second(), n.FireAt.Nanosecond(), now.Location())
	if !isInsideTimeInterval(now, fireAt, fireAt.Add(5*time.Minute)) {
		return false
	}

	lastFireYear, lastFireMonth, lastFireDay := n.LastFireTime.Date()

	return !(lastFireDay == nowDay && lastFireMonth == nowMonth && lastFireYear == nowYear)

}

func isInsideTimeInterval(val, left, right time.Time) bool {
	return val.After(left) && val.Before(right)
}
