package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/paramater"
)

var estateFeatureList = []string{
	"2階以上",
	"駐車場あり",
	"ロフトあり",
	"バストイレ別",
	"DIY可能",
	"ペット飼育可能",
	"インターネット無料",
	"オートロック",
	"駅から徒歩5分",
}

func generateRandomQueryForSearchEstates() url.Values {
	q := url.Values{}
	if (rand.Intn(100) % 10) == 0 {
		q.Set("rentRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 10) == 0 {
		q.Set("doorHeightRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 10) == 0 {
		q.Set("doorWidthRangeId", strconv.Itoa(rand.Intn(4)))
	}

	features := make([]string, len(estateFeatureList))
	copy(features, estateFeatureList)
	rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
	featureLength := rand.Intn(5) + 1
	q.Set("features", strings.Join(features[:featureLength], ","))
	q.Set("perPage", strconv.Itoa(paramater.PerPageOfEstateSearch))
	q.Set("page", "0")

	return q
}

func estateSearchScenario(ctx context.Context, c *client.Client) error {

	t := time.Now()
	err := c.AccessTopPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	t = time.Now()
	err = c.AccessEstateSearchPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	// Search Estates with Query
	q := generateRandomQueryForSearchEstates()

	t = time.Now()
	er, err := c.SearchEstatesWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if len(er.Estates) == 0 {
		return nil
	}

	if !isEstatesOrderedByViewCount(er.Estates) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	numOfPages := int(er.Count) / paramater.PerPageOfEstateSearch
	if numOfPages != 0 {
		for i := 0; i < paramater.NumOfCheckEstateSearchPaging; i++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			er, err := c.SearchEstatesWithQuery(ctx, q)
			if err != nil {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
				return failure.New(fails.ErrTimeout)
			}

			if len(er.Estates) == 0 {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if !isEstatesOrderedByViewCount(er.Estates) {
				err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です"))
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}
			numOfPages = int(er.Count) / paramater.PerPageOfEstateSearch
			if numOfPages == 0 {
				break
			}
		}
	}

	// Get Details with ID from previously searched list
	randomPosition := rand.Intn(len(er.Estates))
	targetID := er.Estates[randomPosition].ID
	t = time.Now()
	e, err := c.AccessEstateDetailPage(ctx, targetID)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	estate, err := asset.GetEstateFromID(e.ID)
	if err != nil || !e.Equal(estate) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	return nil
}
