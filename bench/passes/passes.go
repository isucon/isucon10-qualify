package passes

import (
	"sync"
)

var (
	passCounterMap map[Label]*counter
)

type counter struct {
	count int64
	mu    sync.RWMutex
}

type Label = int

const (
	LabelOfGetChairDetailFromID Label = iota
	LabelOfSearchChairsWithQuery
	LabelOfGetEstateDetailFromID
	LabelOfSearchEstatesWithQuery
	LabelOfSearchEstatesNazotte
	LabelOfGetRecommendedEstatesFromChair
	LabelOfBuyChair
	LabelOfRequestEstateDocument
)

func init() {
	passCounterMap = map[Label]*counter{
		LabelOfGetChairDetailFromID:           {},
		LabelOfSearchChairsWithQuery:          {},
		LabelOfGetEstateDetailFromID:          {},
		LabelOfSearchEstatesWithQuery:         {},
		LabelOfSearchEstatesNazotte:           {},
		LabelOfGetRecommendedEstatesFromChair: {},
		LabelOfBuyChair:                       {},
		LabelOfRequestEstateDocument:          {},
	}
}

func (c *counter) getCount() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func (c *counter) incrementCount() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func GetCount(label Label) int64 {
	return passCounterMap[label].getCount()
}

func IncrementCount(label Label) {
	passCounterMap[label].incrementCount()
}
