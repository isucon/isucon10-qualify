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

func chairSearchScenario(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	var c *client.Client = client.NewClient("isucon-user")

	go func() {
		// Search Chairs with Query
		q := url.Values{}
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

		cr, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		var viewCount int64 = -1
		ok := true
		for i, chair := range cr.Chairs {
			_chair := asset.GetChairFromID(chair.ID)

			if _chair.GetStock() <= 0 {
				ok = false
				break
			}

			vc := _chair.GetViewCount()

			if i > 0 && viewCount < vc {
				ok = false
				break
			}
			viewCount = vc
		}

		if !ok {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		if len(cr.Chairs) == 0 {
			cancel()
			return
		}

		// Get detail of Chair
		randomPosition := rand.Intn(len(cr.Chairs))
		targetID := cr.Chairs[randomPosition].ID
		chair, err := c.GetChairDetailFromID(ctx, strconv.FormatInt(targetID, 10))

		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		ok = chair.Equal(asset.GetChairFromID(chair.ID))
		if !ok {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		chair = asset.GetChairFromID(targetID)

		// Buy Chair
		err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		// Get recommended Estates calculated with Chair
		er, err := c.GetRecommendedEstatesFromChair(ctx, targetID)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		viewCount = -1
		ok = true
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
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		// Get detail of Estate
		randomPosition = rand.Intn(len(er.Estates))
		targetID = er.Estates[randomPosition].ID
		e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		ok = e.Equal(asset.GetEstateFromID(e.ID))
		if !ok {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			cancel()
			return
		}

		// Request docs of Estate
		err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
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
