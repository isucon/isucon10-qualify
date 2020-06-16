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
		q.Set("rentRangeId", strconv.Itoa(rand.Intn(4)))
		if (rand.Intn(100) % 20) == 0 {
			q.Set("doorHeightRangeId", strconv.Itoa(rand.Intn(4)))
		}
		if (rand.Intn(100) % 20) == 0 {
			q.Set("doorWidthRangeId", strconv.Itoa(rand.Intn(4)))
		}
		q.Set("perPage", strconv.Itoa(rand.Intn(20)+30))
		q.Set("page", "0")

		er, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
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
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			cancel()
			return
		}

		// Get Details with ID from previously searched list
		randomPosition := rand.Intn(len(er.Estates))
		targetID := er.Estates[randomPosition].ID
		e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			cancel()
			return
		}

		ok = e.Equal(asset.GetEstateFromID(e.ID))
		if !ok {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			cancel()
			return
		}

		err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		}
		cancel()
	}()

	select {
	case <-ctx.Done():
		return
	case <-timeoutCtx.Done():
		return
	}
}
