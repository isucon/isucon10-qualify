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

	q.Set("perPage", strconv.Itoa(CHAIR_CHECK_PER_PAGE))
	q.Set("page", "0")

	t := time.Now()
	cr, err := c.SearchChairsWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
		return failure.New(fails.ErrTimeout)
	}

	if len(cr.Chairs) == 0 {
		return nil
	}

	ok := checkSearchedChairViewCount(cr.Chairs)
	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	numOfPages := int(cr.Count) / CHAIR_CHECK_PER_PAGE
	if numOfPages != 0 {
		for i := 0; i <= CHAIR_CHECK_PAGE_COUNT; i++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			cr, err := c.SearchChairsWithQuery(ctx, q)
			if err != nil {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > DisengagementResponseTime {
				return failure.New(fails.ErrTimeout)
			}

			if len(cr.Chairs) == 0 {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}

			ok := checkSearchedChairViewCount(cr.Chairs)
			if !ok {
				err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: 検索結果が不正です"))
				fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
				return failure.New(fails.ErrApplication)
			}
		}
	}

	// Get detail of Chair
	randomPosition := rand.Intn(len(cr.Chairs))
	targetID := cr.Chairs[randomPosition].ID
	t = time.Now()
	chair, err := c.GetChairDetailFromID(ctx, strconv.FormatInt(targetID, 10))

	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
		return failure.New(fails.ErrTimeout)
	}

	if chair == nil {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: イス情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if _chair, err := asset.GetChairFromID(chair.ID); err != nil {
		ok = false
	} else {
		ok = chair.Equal(_chair)
	}
	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: イス情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Buy Chair
	err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Get recommended Estates calculated with Chair
	t = time.Now()
	er, err := c.GetRecommendedEstatesFromChair(ctx, targetID)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
		return failure.New(fails.ErrTimeout)
	}

	ok = checkSearchedEstateViewCount(er.Estates)

	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: おすすめ結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Get detail of Estate
	randomPosition = rand.Intn(len(er.Estates))
	targetID = er.Estates[randomPosition].ID
	t = time.Now()
	e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > DisengagementResponseTime {
		return failure.New(fails.ErrTimeout)
	}

	if estate, err := asset.GetEstateFromID(e.ID); err != nil {
		ok = false
	} else {
		ok = e.Equal(estate)
	}

	if !ok {
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
