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

	q.Set("perPage", strconv.Itoa(rand.Intn(20)+30))
	q.Set("page", "0")

	cr, err := c.SearchChairsWithQuery(ctx, q)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if len(cr.Chairs) == 0 {
		return nil
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
		return failure.New(fails.ErrApplication)
	}


	// Get detail of Chair
	randomPosition := rand.Intn(len(cr.Chairs))
	targetID := cr.Chairs[randomPosition].ID
	chair, err := c.GetChairDetailFromID(ctx, strconv.FormatInt(targetID, 10))

	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	ok = chair.Equal(asset.GetChairFromID(chair.ID))
	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	chair = asset.GetChairFromID(targetID)

	// Buy Chair
	err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	// Get recommended Estates calculated with Chair
	er, err := c.GetRecommendedEstatesFromChair(ctx, targetID)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
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
		return failure.New(fails.ErrApplication)
	}

	// Get detail of Estate
	randomPosition = rand.Intn(len(er.Estates))
	targetID = er.Estates[randomPosition].ID
	e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfChairSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	ok = e.Equal(asset.GetEstateFromID(e.ID))
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
