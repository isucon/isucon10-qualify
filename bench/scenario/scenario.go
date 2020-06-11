package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

var (
	ExecutionSeconds = 120
)

func Initialize(ctx context.Context) {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	// timeoutSeconds := 180

	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	err := initialize(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err)
	}
}

func Validation(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Load(ctx)
	}()
	select {
	case <-ctx.Done():
	}
}

func Load(ctx context.Context) {
	var wg sync.WaitGroup
	Num1 := 10
	Num2 := 20

	for i := 0; i < Num1; i++ {
		// 物件検索をして、資料請求をするシナリオ
		wg.Add(1)
		go func() {
			defer wg.Done()

			var c *client.Client
			var e *asset.Estate
			var estate *asset.Estate
			var er *client.EstatesResponse
			var viewCount int64
			var vc int64
			var ok bool
			var err error
			var randomPosition int
			var targetID int64
			var q url.Values

		MAIN:
			for j := 0; j < Num2; j++ {
				ch := time.After(1 * time.Second)
				c = newClient(ctx, "isucon-user")
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				// Search Estates with Query
				q = url.Values{}
				q.Set("doorHeightRangeId", "1")
				q.Set("doorWidthRangeId", "1")
				q.Set("rentRangeId", "1")
				q.Set("perPage", "20")
				q.Set("page", "2")

				er, err = c.SearchEstatesWithQuery(ctx, q)
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				viewCount = -1
				ok = true
				for i, estate := range er.Estates {
					if e = asset.GetEstateFromID(estate.ID); e == nil {
						ok = false
						break
					} else {
						vc = e.GetViewCount()
					}
					if i > 0 && viewCount < vc {
						ok = false
						break
					}
					viewCount = vc
				}

				if !ok {
					fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です")))
					goto Final
				}

				// Get Details with ID from previously searched list
				randomPosition = rand.Intn(len(er.Estates))
				targetID = er.Estates[randomPosition].ID
				e, err = c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				if estate = asset.GetEstateFromID(e.ID); estate == nil {
					ok = false
				} else {
					ok = e.Equal(estate)
				}
				if !ok {
					fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Messagef("GET /api/estate/%d: 物件情報が不正です", targetID)))
					goto Final
				}

				err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}
			Final:
				select {
				case <-ch:
				case <-ctx.Done():
					break MAIN
				}
			}
		}()
	}
}

func newClient(ctx context.Context, userAgent string) *client.Client {
	return client.NewClient(userAgent)
}
