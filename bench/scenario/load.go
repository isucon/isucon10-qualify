package scenario

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

const (
	SleepTimeOnFailScenario        = 1 * time.Second
	IntervalForCheckWorkers        = 10 * time.Second
	NumOfInitialEstateSearchWorker = 1
	NumOfInitialChairSearchWorker  = 1
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
		err := estateSearchScenario(ctx)
		if err != nil {
			t = time.NewTimer(SleepTimeOnFailScenario)
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func checkEstateSearchWorker(ctx context.Context) {
	t := time.NewTicker(IntervalForCheckWorkers)
	for {
		select {
		case <-t.C:
			et := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateSearchScenario)
			if time.Since(et) > IntervalForCheckWorkers {
				log.Println("物件検索シナリオの負荷レベルが上昇しました。")
				go runEstateSearchWorker(ctx)
			} else {
				log.Println("物件検索シナリオでエラーが発生したため負荷レベルを上げられませんでした。")
			}
		case <-ctx.Done():
			t.Stop()
			return
		}
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
		err := chairSearchScenario(ctx)
		if err != nil {
			t = time.NewTimer(SleepTimeOnFailScenario)
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func checkChairSearchWorker(ctx context.Context) {
	t := time.NewTicker(IntervalForCheckWorkers)
	for {
		select {
		case <-t.C:
			et := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfChairSearchScenario)
			if time.Since(et) > IntervalForCheckWorkers {
				log.Println("イス検索シナリオの負荷レベルが上昇しました。")
				go runChairSearchWorker(ctx)
			} else {
				log.Println("イス検索シナリオでエラーが発生したため負荷レベルを上げられませんでした。")
			}
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func Load(ctx context.Context) {
	// 物件検索をして、資料請求をするシナリオ
	for i := 0; i < NumOfInitialEstateSearchWorker; i++ {
		go runEstateSearchWorker(ctx)
	}
	go checkEstateSearchWorker(ctx)

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < NumOfInitialChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}
	go checkChairSearchWorker(ctx)
}
