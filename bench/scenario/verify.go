package scenario

import (
	"context"
	"net/url"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

// Verify Initialize後のアプリケーションサーバーに対して、副作用のない検証を実行する
// 早い段階でベンチマークをFailさせて早期リターンさせるのが目的
// ex) popular API や Search API を叩いて初期状態を確認する
func Verify(ctx context.Context, snapshotsParentsDirPath, fixtureDir string) {
	c := client.NewClientForVerify()
	verifyWithSnapshot(ctx, c, snapshotsParentsDirPath)
	verifyWithScenario(ctx, c, fixtureDir)
}

func verifyChairStock(ctx context.Context, c *client.Client, chairFeatureForVerify string) error {
	err := c.BuyChair(ctx, "1")
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	chair, err := c.GetChairDetailFromID(ctx, "1")
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if chair != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.New(fails.ErrApplication, failure.Message("イスの在庫数が不正です"))
	}

	return nil
}

func verifyChairViewCount(ctx context.Context, c *client.Client, chairFeatureForVerify string) error {
	for i := 0; i < 2; i++ {
		_, err := c.GetChairDetailFromID(ctx, "2")
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			return failure.Translate(err, fails.ErrApplication)
		}
	}

	q := url.Values{}
	q.Add("features", chairFeatureForVerify)
	q.Add("page", "0")
	q.Add("perPage", "2")

	chairs, err := c.SearchChairsWithQuery(ctx, q)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if chairs.Chairs[0].ID != 2 || chairs.Chairs[1].ID != 3 {
		return failure.New(fails.ErrApplication, failure.Message("イスの閲覧数が不正です"))
	}

	return nil
}

func verifyEstateViewCount(ctx context.Context, c *client.Client, estateFeatureForVerify string) error {
	for i := 0; i < 2; i++ {
		_, err := c.GetEstateDetailFromID(ctx, "1")
		if err != nil {
			return failure.Translate(err, fails.ErrApplication)
		}
	}

	q := url.Values{}
	q.Add("features", estateFeatureForVerify)
	q.Add("page", "0")
	q.Add("perPage", "2")

	estates, err := c.SearchEstatesWithQuery(ctx, q)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if estates.Estates[0].ID != 1 || estates.Estates[1].ID != 2 {
		return failure.New(fails.ErrApplication, failure.Message("物件の閲覧数が不正です"))
	}

	return nil
}

func verifyWithScenario(ctx context.Context, c *client.Client, fixtureDir string) {
	chairFeatureForVerify, err := asset.GetChairFeatureForVerify()
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	}

	estateFeatureForVerify, err := asset.GetEstateFeatureForVerify()
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := verifyChairStock(ctx, c, *chairFeatureForVerify)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := verifyChairViewCount(ctx, c, *chairFeatureForVerify)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := verifyEstateViewCount(ctx, c, *estateFeatureForVerify)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Wait()
}
