package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

func estateSearchScenario(ctx context.Context) error {
	passCtx, pass := context.WithCancel(ctx)
	failCtx, fail := context.WithCancel(ctx)

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
			fail()
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
			fail()
			return
		}

		// Get Details with ID from previously searched list
		randomPosition := rand.Intn(len(er.Estates))
		targetID := er.Estates[randomPosition].ID
		e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			fail()
			return
		}

		ok = e.Equal(asset.GetEstateFromID(e.ID))
		if !ok {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			fail()
			return
		}

		err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
			fail()
		}

		pass()
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-failCtx.Done():
		return failure.New(fails.ErrApplication)
	case <-passCtx.Done():
		return nil
	}
}
