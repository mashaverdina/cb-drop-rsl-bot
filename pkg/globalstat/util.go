package globalstat

import (
	"fmt"
	"sort"

	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/messages"
)

type statPair struct {
	itemsCount int
	usersCount int
}

type CumSum struct {
	Total int
	Users []int
}

func NewCumSum(pairs []statPair) *CumSum {
	if len(pairs) == 0 {
		return &CumSum{
			Total: 0,
			Users: make([]int, 0),
		}
	}
	// assert ordered by items
	maxItems := pairs[len(pairs)-1].itemsCount
	users := make([]int, maxItems+1)
	for _, pair := range pairs {
		users[pair.itemsCount] = pair.usersCount
	}
	t := 0
	for i := range users {
		t = users[i] + t
		users[i] = t
	}
	return &CumSum{
		Total: t,
		Users: users,
	}
}

func (c CumSum) GetRating(count int) float64 {
	if count >= len(c.Users) {
		return 1.
	}
	usersHasLess := c.Users[count]
	if count > 0 {
		diff := (usersHasLess - c.Users[count-1]) / 2
		usersHasLess = usersHasLess - diff
	} else {
		usersHasLess = usersHasLess / 2
	}
	// return float64(c.Total - usersHasLess) / float64(c.Total)
	return float64(usersHasLess) / float64(c.Total)
}

type Stat struct {
	data map[string]CumSum
}

func NewStat(stats []entities.UserCbStat) *Stat {
	result := make(map[string]CumSum)

	for level := 4; level <= 6; level++ {
		for _, itemType := range []string{
			messages.AncientShard, messages.VoidShard,
			messages.SacredShard, messages.EpicTome, messages.LegTome,
		} {
			data := make(map[int]int)
			for _, item := range stats {
				if item.Level == level {
					itemsCount := extract(item, itemType)
					if _, ok := data[itemsCount]; !ok {
						data[itemsCount] = 0
					}
					data[itemsCount] = data[itemsCount] + 1
				}
			}
			pairs := make([]statPair, 0, len(data))
			for itemsCount, usersCount := range data {
				pairs = append(pairs, statPair{
					itemsCount: itemsCount,
					usersCount: usersCount,
				})
			}
			sort.Slice(pairs, func(i, j int) bool {
				return pairs[i].itemsCount < pairs[j].itemsCount
			})
			key := makeKey(itemType, level)
			result[key] = *NewCumSum(pairs)
		}
	}
	return &Stat{data: result}
}

func (s Stat) GetCumSum(itemType string, level int) CumSum {
	return s.data[makeKey(itemType, level)]
}

func makeKey(itemType string, level int) string {
	return fmt.Sprintf("%s:%d", itemType, level)
}

func extract(item entities.UserCbStat, itemType string) int {
	switch itemType {
	case messages.AncientShard:
		return item.AncientShard
	case messages.VoidShard:
		return item.VoidShard
	case messages.SacredShard:
		return item.SacredShard
	case messages.EpicTome:
		return item.EpicTome
	case messages.LegTome:
		return item.LegTome
	default:
		// todo
		panic("unknown item type " + itemType)
	}
}
