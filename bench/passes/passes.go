package passes

import (
	"sync"
	"time"
)

var (
	passMap map[Label]*pass
)

type pass struct {
	durations []time.Duration
	mu        sync.RWMutex
}

type Label = int

const (
	LabelOfGetChairDetailFromID Label = iota
	LabelOfSearchChairsWithQuery
	LabelOfGetEstateDetailFromID
	LabelOfSearchEstatesWithQuery
	LabelOfSearchEstatesNazotte
	LabelOfGetRecommendedChair
	LabelOfGetRecommendedEstate
	LabelOfGetRecommendedEstatesFromChair
	LabelOfBuyChair
	LabelOfRequestEstateDocument
)

func init() {
	passMap = map[Label]*pass{
		LabelOfGetChairDetailFromID:           {durations: make([]time.Duration, 0)},
		LabelOfSearchChairsWithQuery:          {durations: make([]time.Duration, 0)},
		LabelOfGetEstateDetailFromID:          {durations: make([]time.Duration, 0)},
		LabelOfSearchEstatesWithQuery:         {durations: make([]time.Duration, 0)},
		LabelOfSearchEstatesNazotte:           {durations: make([]time.Duration, 0)},
		LabelOfGetRecommendedChair:            {durations: make([]time.Duration, 0)},
		LabelOfGetRecommendedEstate:           {durations: make([]time.Duration, 0)},
		LabelOfGetRecommendedEstatesFromChair: {durations: make([]time.Duration, 0)},
		LabelOfBuyChair:                       {durations: make([]time.Duration, 0)},
		LabelOfRequestEstateDocument:          {durations: make([]time.Duration, 0)},
	}
}

func (p *pass) getDurations() []time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.durations
}

func (p *pass) getCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.durations)
}

func (p *pass) addDuration(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.durations = append(p.durations, d)
}

func GetDurations(label Label) []time.Duration {
	return passMap[label].getDurations()
}

func GetCount(label Label) int {
	return passMap[label].getCount()
}

func AddDuration(d time.Duration, label Label) {
	passMap[label].addDuration(d)
}
