package scenario

import (
	"context"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func botScenario(ctx context.Context, c *client.Client) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		q := generateRandomQueryForSearchChairs()
		q.Set("perPage", "10")
		_, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
			}
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		q := generateRandomQueryForSearchEstates()
		q.Set("perPage", "10")
		_, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
			}
			return
		}
	}()

	wg.Wait()
}
