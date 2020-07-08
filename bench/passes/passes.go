package passes

import (
	"sync"
)

var (
	passMap map[Label]int
	mu      sync.RWMutex
)

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
	LabelOfStaticFiles
)

func init() {
	passMap = map[Label]int{}
}

func GetCount(label Label) int {
	mu.RLock()
	defer mu.RUnlock()
	return passMap[label]
}

func IncrementCount(label Label) {
	mu.Lock()
	defer mu.Unlock()
	passMap[label]++
}
