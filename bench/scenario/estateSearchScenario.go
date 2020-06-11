package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

func estateSearchScenario(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	var c *client.Client = client.NewClient("isucon-user")

	go func() {
		// Search Estates with Query
		q := url.Values{}
		q.Set("doorHeightRangeId", "1")
		q.Set("doorWidthRangeId", "1")
		q.Set("rentRangeId", "1")
		q.Set("perPage", "20")
		q.Set("page", "0")

		er, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			fails.ErrorsForCheck.Add(err)
			cancel()
			return
		}

		var viewCount int64 = -1
		ok := true
		for i, estate := range er.Estates {
			e := asset.GetEstateFromID(estate.ID)
			vc := e.GetViewCount()
			if i > 0 && viewCount < vc {
				ok = false
				break
			}
			viewCount = vc
		}

		if !ok {
			fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です")))
			cancel()
			return
		}

		// Get Details with ID from previously searched list
		randomPosition := rand.Intn(len(er.Estates))
		targetID := er.Estates[randomPosition].ID
		e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err)
			cancel()
			return
		}

		ok = e.Equal(asset.GetEstateFromID(e.ID))
		if !ok {
			fails.ErrorsForCheck.Add(failure.New(fails.ErrApplication, failure.Messagef("GET /api/estate/%d: 物件情報が不正です", targetID)))
			cancel()
			return
		}

		err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

		if err != nil {
			fails.ErrorsForCheck.Add(err)
			cancel()
			return
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-timeoutCtx.Done():
		return
	}
}
