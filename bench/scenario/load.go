package scenario

import (
	"context"
	"math/rand"
	"time"
)

const (
	NumOfInitialEstateSearchWorker = 5
	NumOfInitialChairSearchWorker  = 5
)

func runEstateSearchWorker(ctx context.Context) {
	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		estateSearchScenario(ctx)
	}
}

func runChairSearchWorker(ctx context.Context) {
	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		chairSearchScenario(ctx)
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
