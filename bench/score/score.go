package score

import (
	"sync"
)

var (
	score           int64 = 0
	level           int64 = 0
	levelChan       chan int64
	boundaryOfLevel []int64 = []int64{
		400, 800, 1200, 1600, 2000,
		2400, 2800, 3200, 3600, 4000,
		4400, 4800, 5200, 5600, 6000,
		6400, 6800, 7200, 7600, 8000,
		8400, 8800, 9200, 9600, 10000,
	}
	mu sync.RWMutex
)

func init() {
	levelChan = make(chan int64, 1)
}

func IncrementScore() {
	mu.Lock()
	defer mu.Unlock()
	score++
	if score >= boundaryOfLevel[level] {
		level++
		levelChan <- level
	}
}

func GetScore() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return score
}

func GetLevel() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return level
}

func LevelUp() chan int64 {
	return levelChan
}
