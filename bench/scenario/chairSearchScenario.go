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

var chairKindList = []string{
	"ゲーミングチェア",
	"座椅子",
	"エルゴノミクス",
	"ハンモック",
}

var chairColorList = []string{
	"黒",
	"白",
	"赤",
	"青",
	"緑",
	"黄",
	"紫",
	"ピンク",
	"オレンジ",
	"水色",
	"ネイビー",
	"ベージュ",
}

var chairFeatureList = []string{
	"折りたたみ可",
	"肘掛け",
	"キャスター",
	"リクライニング",
	"高さ調節可",
	"フットレスト",
}

func chairSearchScenario(ctx context.Context) error {
	var c *client.Client = client.PickClient()

	t := time.Now()
	err := c.AccessTopPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	t = time.Now()
	err = c.AccessChairSearchPage(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

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

	if (rand.Intn(100) % 20) == 0 {
		q.Set("kind", chairKindList[rand.Intn(len(chairKindList))])
	}
	if (rand.Intn(100) % 20) == 0 {
		q.Set("color", chairColorList[rand.Intn(len(chairColorList))])
	}
	if (rand.Intn(100) % 20) == 0 {
		features := make([]string, len(chairFeatureList))
		copy(features, chairFeatureList)
		rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
		featureLength := rand.Intn(3) + 1
		q.Set("features", strings.Join(features[:featureLength], ","))
	}

	q.Set("perPage", strconv.Itoa(paramater.PerPageOfChairSearch))
	q.Set("page", "0")

	t = time.Now()
	cr, err := c.SearchChairsWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if len(cr.Chairs) == 0 {
		return nil
	}

	if !isChairsOrderedByViewCount(cr.Chairs) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	numOfPages := int(cr.Count) / paramater.PerPageOfChairSearch

	if numOfPages != 0 {
		for i := 0; i < paramater.NumOfCheckChairSearchPaging; i++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			cr, err := c.SearchChairsWithQuery(ctx, q)
			if err != nil {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
				return failure.New(fails.ErrTimeout)
			}

			if len(cr.Chairs) == 0 {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if !isChairsOrderedByViewCount(cr.Chairs) {
				err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}
			numOfPages = int(cr.Count) / paramater.PerPageOfChairSearch
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

	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if chair == nil {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: イス情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if !isChairEqualToAsset(chair) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: イス情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if !isEstatesOrderedByViewCount(er.Estates) {
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

	if time.Since(t) > paramater.ThresholdTimeOfAbandonmentPage {
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
