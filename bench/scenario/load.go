package scenario

import (
	"context"
	"math/rand"
	"time"
)

const (
	NumOfInitialEstateSearchWorker = 5
	NumOfInitialChairSearchWorker = 5
)

func runEstateSearchWorker(ctx context.Context) {
	for {
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		estateSearchScenario(ctx)
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func runChairSearchWorker(ctx context.Context) {
	for {
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		chairSearchScenario(ctx)
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func Load(ctx context.Context) {
	// 物件検索をして、資料請求をするシナリオ
	for i := 0; i < NumOfInitialEstateSearchWorker; i++ {
		go runEstateSearchWorker(ctx)
	}

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < NumOfInitialChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}
}
