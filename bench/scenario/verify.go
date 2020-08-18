package scenario

import (
	"context"
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

func verifyWithScenario(ctx context.Context, c *client.Client, fixtureDir string) {
	chairFeatureForVerify, err := asset.GetChairFeatureForVerify()
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

	wg.Wait()
}
