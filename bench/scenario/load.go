package scenario

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

const (
	SleepTimeOnFailScenario               = 1 * time.Second
	SleepSwingOnFailScenario              = 1000 // * time.Millisecond
	SleepTimeOnUserAway                   = 500 * time.Millisecond
	SleepSwingOnUserAway                  = 100 // * time.Millisecond
	IntervalForCheckWorkers               = 10 * time.Second
	NumOfInitialEstateSearchWorker        = 1
	NumOfInitialChairSearchWorker         = 1
	NumOfInitialEstateNazotteSearchWorker = 1
	DisengagementResponseTime             = 1 * time.Second
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
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(SleepSwingOnUserAway) - SleepSwingOnUserAway*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(SleepSwingOnFailScenario) - SleepSwingOnFailScenario*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
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
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(SleepSwingOnUserAway) - SleepSwingOnUserAway*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(SleepSwingOnFailScenario) - SleepSwingOnFailScenario*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func runEstateNazotteSearchWorker(ctx context.Context) {
	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := estateNazotteSearchScenario(ctx)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(SleepSwingOnUserAway) - SleepSwingOnUserAway*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(SleepSwingOnFailScenario) - SleepSwingOnFailScenario*0.5
				s := SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func checkWorkers(ctx context.Context) {
	t := time.NewTicker(IntervalForCheckWorkers)
	for {
		select {
		case <-t.C:
			cet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfChairSearchScenario)
			eet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateSearchScenario)
			net := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateNazotteSearchScenario)
			if time.Since(cet) > IntervalForCheckWorkers &&
				time.Since(eet) > IntervalForCheckWorkers &&
				time.Since(net) > IntervalForCheckWorkers {
				log.Println("負荷レベルが上昇しました。")
				go runChairSearchWorker(ctx)
				go runEstateSearchWorker(ctx)
				go runEstateNazotteSearchWorker(ctx)
			} else {
				log.Println("シナリオ内でエラーが発生したため負荷レベルを上げられませんでした。")
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

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < NumOfInitialChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}

	// なぞって検索をするシナリオ
	for i := 0; i < NumOfInitialEstateNazotteSearchWorker; i++ {
		go runEstateNazotteSearchWorker(ctx)
	}

	go checkWorkers(ctx)
}
