package scenario

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/morikuni/failure"
)

var loadLevel int64

func GetLoadLevel() int64 {
	return atomic.LoadInt64(&loadLevel)
}

func runEstateSearchWorker(ctx context.Context) {
	u, _ := uuid.NewRandom()
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()), false)

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
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()), false)

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
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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
	c := client.NewClient(fmt.Sprintf("isucon-user-%v", u.String()), false)

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
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
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

func runBotWorker(ctx context.Context) {
	u, _ := uuid.NewRandom()
	c := client.NewClient(fmt.Sprintf("isucon-bot-%v", u.String()), true)

	for {
		go botScenario(ctx, c)
		r := rand.Intn(parameter.SleepSwingOnBotInterval) - parameter.SleepSwingOnBotInterval*0.5
		s := parameter.SleepTimeOnBotInterval + time.Duration(r)*time.Millisecond
		t := time.NewTimer(s)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func checkWorkers(ctx context.Context) {
	t := time.NewTicker(parameter.IntervalForCheckWorkers)
	for {
		select {
		case <-t.C:
			cet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfChairSearchScenario)
			eet := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateSearchScenario)
			net := fails.ErrorsForCheck.GetLastErrorTime(fails.ErrorOfEstateNazotteSearchScenario)
			if time.Since(cet) > parameter.IntervalForCheckWorkers &&
				time.Since(eet) > parameter.IntervalForCheckWorkers &&
				time.Since(net) > parameter.IntervalForCheckWorkers {
				log.Println("負荷レベルが上昇しました。")
				level := atomic.AddInt64(&loadLevel, 1)
				incWorkers := parameter.ListOfIncWorkers[level]
				for i := 0; i < incWorkers.ChairSearchWorker; i++ {
					go runChairSearchWorker(ctx)
				}
				for i := 0; i < incWorkers.EstateSearchWorker; i++ {
					go runEstateSearchWorker(ctx)
				}
				for i := 0; i < incWorkers.EstateNazotteSearchWorker; i++ {
					go runEstateNazotteSearchWorker(ctx)
				}
				for i := 0; i < incWorkers.BotWorker; i++ {
					go runBotWorker(ctx)
				}
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
	level := GetLoadLevel()
	incWorkers := parameter.ListOfIncWorkers[level]

	// 物件検索をして、資料請求をするシナリオ
	for i := 0; i < incWorkers.ChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < incWorkers.EstateSearchWorker; i++ {
		go runEstateSearchWorker(ctx)
	}

	// なぞって検索をするシナリオ
	for i := 0; i < incWorkers.EstateNazotteSearchWorker; i++ {
		go runEstateNazotteSearchWorker(ctx)
	}

	// ボットによる検索シナリオ
	for i := 0; i < incWorkers.BotWorker; i++ {
		go runBotWorker(ctx)
	}

	go checkWorkers(ctx)
}
