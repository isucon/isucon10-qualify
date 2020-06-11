package scenario

import (
	"context"
	"sync"
)

func Load(ctx context.Context) {
	var wg sync.WaitGroup
	WorkloadNum := 5
	Scenario1Num := 3
	Scenario2Num := 3

	for i := 0; i < WorkloadNum; i++ {
		// 物件検索をして、資料請求をするシナリオ
		for j := 0; j < Scenario1Num; j++ {
			wg.Add(1)
			go func() {
				estateSearchScenario(ctx)
				wg.Done()
			}()
		}

		// イス検索から物件ページに行き、資料請求をするまでのシナリオ
		for j := 0; j < Scenario2Num; j++ {
			wg.Add(1)
			go func() {
				chairSearchScenario(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
