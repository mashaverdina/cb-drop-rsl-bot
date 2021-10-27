package globalstat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"vkokarev.com/rslbot/pkg/storage"
	"vkokarev.com/rslbot/pkg/utils"
)

type GlobalStatManager struct {
	cbStatStorage *storage.CbStatStorage
	m             sync.Mutex
	started       bool
	cancel        context.CancelFunc
	ctx           context.Context
	monthStat     map[time.Month]*Stat
	last7dStat    *Stat
	last30dStat   *Stat
}

func NewGlobalStatManager(cbStatStorage *storage.CbStatStorage) *GlobalStatManager {
	return &GlobalStatManager{
		cbStatStorage: cbStatStorage,
		m:             sync.Mutex{},
		started:       false,
		cancel:        nil,
		ctx:           nil,
		monthStat:     make(map[time.Month]*Stat),
	}
}

func (gs *GlobalStatManager) Start(ctx context.Context) error {
	gs.m.Lock()
	defer gs.m.Unlock()
	if gs.started {
		return errors.New("started already")
	}
	gs.ctx, gs.cancel = context.WithCancel(ctx)

	if err := gs.update(false); err != nil {
		return err
	}

	go gs.loop()
	gs.started = true
	return nil
}

func (gs *GlobalStatManager) Stop() error {
	gs.m.Lock()
	defer gs.m.Unlock()
	if !gs.started {
		return errors.New("not started")
	}
	gs.cancel()
	// todo done chan
	return nil
}

func (gs *GlobalStatManager) loop() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			if err := gs.update(true); err != nil {
				log.Println(fmt.Sprintf("can't update global stat: %v", err))
			} else {
				log.Println("global stat updated")
			}
		case <-gs.ctx.Done():
			return
		}
	}
}

func (gs *GlobalStatManager) update(doLock bool) error {
	if doLock {
		gs.m.Lock()
		defer gs.m.Unlock()
	}

	for _, month := range []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December} {
		from, to := utils.MonthInterval(month)
		minDays := 20
		if time.Now().Month() == month {
			minDays = int(0.8 * float64(time.Now().Day()))
		}
		stats, err := gs.cbStatStorage.FullStat(from, to, minDays)
		if err != nil {
			return err
		}
		gs.monthStat[month] = NewStat(stats)
	}
	{
		from, to := utils.LastDaysInterval(7)
		stats, err := gs.cbStatStorage.FullStat(from, to, 5)
		if err != nil {
			return err
		}
		gs.last7dStat = NewStat(stats)
	}

	{
		from, to := utils.LastDaysInterval(30)
		stats, err := gs.cbStatStorage.FullStat(from, to, 25)
		if err != nil {
			return err
		}
		gs.last30dStat = NewStat(stats)
	}
	return nil
}

func (gs *GlobalStatManager) TopFor7Days(level int, itemType string, count int) (float64, error) {
	return gs.last7dStat.GetCumSum(itemType, level).GetRating(count), nil
}

func (gs *GlobalStatManager) TopFor30Days(level int, itemType string, count int) (float64, error) {
	return gs.last30dStat.GetCumSum(itemType, level).GetRating(count), nil
}

func (gs *GlobalStatManager) TopForMonth(month time.Month, level int, itemType string, count int) (float64, error) {
	return gs.monthStat[month].GetCumSum(itemType, level).GetRating(count), nil
}
