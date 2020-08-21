package scenario

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

func chairSearchScenario(ctx context.Context, c *client.Client) error {
	t := time.Now()
	err := c.AccessTopPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	t = time.Now()
	err = c.AccessChairSearchPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	// Search Chairs with Query
	q, err := createRandomChairSearchQuery()
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	t = time.Now()
	cr, err := c.SearchChairsWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if len(cr.Chairs) == 0 {
		return nil
	}

	if !isChairsOrderedByPrice(cr.Chairs, t) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	numOfPages := int(cr.Count) / parameter.PerPageOfChairSearch

	if numOfPages != 0 {
		for i := 0; i < parameter.NumOfCheckChairSearchPaging; i++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			cr, err := c.SearchChairsWithQuery(ctx, q)
			if err != nil {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
				return failure.New(fails.ErrTimeout)
			}

			if len(cr.Chairs) == 0 {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if !isChairsOrderedByPrice(cr.Chairs, t) {
				err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}
			numOfPages = int(cr.Count) / parameter.PerPageOfChairSearch
			if numOfPages == 0 {
				break
			}
		}
	}

	// Get detail of Chair
	randomPosition := rand.Intn(len(cr.Chairs))
	targetID := cr.Chairs[randomPosition].ID
	t = time.Now()
	chair, er, err := c.AccessChairDetailPage(ctx, targetID)

	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if chair == nil {
		return nil
	}

	if !isChairEqualToAsset(chair) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: イス情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if !isEstatesOrderedByPopularity(er.Estates) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: おすすめ結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Buy Chair
	err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		if _chair, err := asset.GetChairFromID(targetID); err != nil || _chair.GetStock() > 0 {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
			return failure.New(fails.ErrApplication)
		}
	}

	// Get detail of Estate
	randomPosition = rand.Intn(len(er.Estates))
	targetID = er.Estates[randomPosition].ID
	t = time.Now()
	e, err := c.AccessEstateDetailPage(ctx, targetID)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if !isEstateEqualToAsset(e) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Request docs of Estate
	err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	return nil
}
