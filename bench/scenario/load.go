package scenario

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/paramater"
	"github.com/morikuni/failure"
	"github.com/google/uuid"
)

func runEstateSearchWorker(ctx context.Context) {
	u, _ := uuid.NewRandom()
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()))

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := estateSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(paramater.SleepSwingOnUserAway) - paramater.SleepSwingOnUserAway*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(paramater.SleepSwingOnFailScenario) - paramater.SleepSwingOnFailScenario*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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
	u, _ := uuid.NewRandom()
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()))

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := chairSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(paramater.SleepSwingOnUserAway) - paramater.SleepSwingOnUserAway*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(paramater.SleepSwingOnFailScenario) - paramater.SleepSwingOnFailScenario*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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
	u, _ := uuid.NewRandom()
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()))

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := estateNazotteSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(paramater.SleepSwingOnUserAway) - paramater.SleepSwingOnUserAway*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(paramater.SleepSwingOnFailScenario) - paramater.SleepSwingOnFailScenario*0.5
				s := paramater.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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
	t := time.NewTicker(paramater.IntervalForCheckWorkers)
	for {
		select {
		case <-t.C:
			cet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfChairSearchScenario)
			eet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateSearchScenario)
			net := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateNazotteSearchScenario)
			if time.Since(cet) > paramater.IntervalForCheckWorkers &&
				time.Since(eet) > paramater.IntervalForCheckWorkers &&
				time.Since(net) > paramater.IntervalForCheckWorkers {
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
	for i := 0; i < paramater.NumOfInitialEstateSearchWorker; i++ {
		go runEstateSearchWorker(ctx)
	}

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < paramater.NumOfInitialChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}

	// なぞって検索をするシナリオ
	for i := 0; i < paramater.NumOfInitialEstateNazotteSearchWorker; i++ {
		go runEstateNazotteSearchWorker(ctx)
	}

	go checkWorkers(ctx)
}
