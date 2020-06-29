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
)

var estateFeatureList = []string{
	"バストイレ別",
	"駅から徒歩5分",
	"ペット飼育可能",
	"デザイナーズ物件",
}

func estateSearchScenario(ctx context.Context) error {
	var c *client.Client = client.PickClient()

	// Search Estates with Query
	q := url.Values{}
	q.Set("rentRangeId", strconv.Itoa(rand.Intn(4)))
	if (rand.Intn(100) % 20) == 0 {
		q.Set("doorHeightRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 20) == 0 {
		q.Set("doorWidthRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 20) == 0 {
		features := make([]string, len(estateFeatureList))
		copy(features, estateFeatureList)
		rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
		featureLength := rand.Intn(3) + 1
		q.Set("features", strings.Join(features[:featureLength], ","))
	}
	q.Set("perPage", strconv.Itoa(ESTATE_CHECK_PER_PAGE))
	q.Set("page", "0")

	t := time.Now()
	er, err := c.SearchEstatesWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
		return failure.New(fails.ErrTimeout)
	}

	if len(er.Estates) == 0 {
		return nil
	}

	ok := checkSearchedEstateViewCount(er.Estates)

	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	numOfPages := int(er.Count) / ESTATE_CHECK_PER_PAGE
	if numOfPages != 0 {
		for i := 0; i < ESTATE_CHECK_PAGE_COUNT; i++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			er, err := c.SearchEstatesWithQuery(ctx, q)
			if err != nil {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > DisengagementResponseTime {
				return failure.New(fails.ErrTimeout)
			}

			if len(er.Estates) == 0 {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			ok := checkSearchedEstateViewCount(er.Estates)
			if !ok {
				err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 検索結果が不正です"))
				fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
				return failure.New(fails.ErrApplication)
			}
			numOfPages = int(er.Count) / ESTATE_CHECK_PER_PAGE
			if numOfPages == 0 {
				break
			}
		}
	}

	// Get Details with ID from previously searched list
	randomPosition := rand.Intn(len(er.Estates))
	targetID := er.Estates[randomPosition].ID
	t = time.Now()
	e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
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
