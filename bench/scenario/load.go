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

func Load(ctx context.Context) {
	var wg sync.WaitGroup
	WorkloadNum := 5
	Scenario1Num := 3
	Scenario2Num := 3

	for i := 0; i < WorkloadNum; i++ {
		// 物件検索をして、資料請求をするシナリオs
		wg.Add(1)
		go func() {
			defer wg.Done()

			var c *client.Client
			var e *asset.Estate
			var er *client.EstatesResponse
			var viewCount int64
			var vc int64
			var ok bool
			var err error
			var randomPosition int
			var targetID int64
			var q url.Values

		MAIN:
			for j := 0; j < Scenario1Num; j++ {
				ch := time.After(1 * time.Second)
				c = client.NewClient("isucon-user")
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
				q.Set("page", "0")

				er, err = c.SearchEstatesWithQuery(ctx, q)
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				viewCount = -1
				ok = true
				for i, estate := range er.Estates {
					e = asset.GetEstateFromID(estate.ID)
					vc = e.GetViewCount()
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

				ok = e.Equal(asset.GetEstateFromID(e.ID))
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

		// イス検索から物件ページに行き、資料請求をするまでのシナリオ
		wg.Add(1)
		go func() {
			defer wg.Done()

			var c *client.Client
			var chair *asset.Chair
			var _chair *asset.Chair
			var e *asset.Estate
			var cr *client.ChairsResponse
			var er *client.EstatesResponse
			var viewCount int64
			var vc int64
			var ok bool
			var err error
			var randomPosition int
			var targetID int64
			var q url.Values

		MAIN:
			for j := 0; j < Scenario2Num; j++ {
				ch := time.After(1 * time.Second)
				c = client.NewClient("isucon-user")
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				// Search Chairs with Query
				q = url.Values{}
				q.Set("priceRangeId", strconv.Itoa(rand.Intn(6)))
				if (rand.Intn(100) % 5) == 0 {
					q.Set("heightRangeId", strconv.Itoa(rand.Intn(4)))
				}
				if (rand.Intn(100) % 5) == 0 {
					q.Set("widthRangeId", strconv.Itoa(rand.Intn(4)))
				}
				if (rand.Intn(100) % 5) == 0 {
					q.Set("depthRangeId", strconv.Itoa(rand.Intn(4)))
				}

				if (rand.Intn(100) % 4) == 0 {
					// q.Set("kind", "エルゴノミクス")
				}
				if (rand.Intn(100) % 4) == 0 {
					// q.Set("color", "black")
				}
				if (rand.Intn(100) % 4) == 0 {
					// q.Set("features", "リクライニング,肘掛け")
				}

				q.Set("perPage", strconv.Itoa(rand.Intn(20)+30))
				q.Set("page", "0")

				cr, err = c.SearchChairsWithQuery(ctx, q)
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				viewCount = -1
				ok = true
				for i, chair := range cr.Chairs {
					_chair = asset.GetChairFromID(chair.ID)

					if _chair.GetStock() <= 0 {
						ok = false
						break
					}

					vc = _chair.GetViewCount()

					if i > 0 && viewCount < vc {
						ok = false
						break
					}
					viewCount = vc
				}

				if !ok {
					fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です")))
					goto Final
				}

				if len(cr.Chairs) == 0 {
					continue
				}

				// Get detail of Chair
				randomPosition = rand.Intn(len(cr.Chairs))
				targetID = cr.Chairs[randomPosition].ID
				chair, err = c.GetChairDetailFromID(ctx, strconv.FormatInt(targetID, 10))

				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				ok = chair.Equal(asset.GetChairFromID(chair.ID))
				if !ok {
					fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Messagef("GET /api/estate/%d: 物件情報が不正です", targetID)))
					goto Final
				}

				chair = asset.GetChairFromID(targetID)

				// Buy Chair
				err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				// Get recommended Estates calculated with Chair
				er, err = c.GetRecommendedEstatesFromChair(ctx, targetID)
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				viewCount = -1
				ok = true
				for i, estate := range er.Estates {
					e = asset.GetEstateFromID(estate.ID)
					vc = e.GetViewCount()
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

				// Get detail of Estate
				randomPosition = rand.Intn(len(er.Estates))
				targetID = er.Estates[randomPosition].ID
				e, err = c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
				if err != nil {
					fails.ErrorsForCheck.Add(err)
					goto Final
				}

				ok = e.Equal(asset.GetEstateFromID(e.ID))
				if !ok {
					fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Messagef("GET /api/estate/%d: 物件情報が不正です", targetID)))
					goto Final
				}

				// Request docs of Estate
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
