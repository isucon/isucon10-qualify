package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

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
	q.Set("perPage", strconv.Itoa(rand.Intn(20)+30))
	q.Set("page", "0")

	er, err := c.SearchEstatesWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if len(er.Estates) == 0 {
		return nil
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
		return failure.New(fails.ErrApplication)
	}

	// Get Details with ID from previously searched list
	randomPosition := rand.Intn(len(er.Estates))
	targetID := er.Estates[randomPosition].ID
	e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	ok = e.Equal(asset.GetEstateFromID(e.ID))
	if !ok {
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
