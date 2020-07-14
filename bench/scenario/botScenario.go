package scenario

import (
	"context"
	"strconv"
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
		chairs, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
			}
			return
		}

		for _, chair := range chairs.Chairs {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_, err := c.GetChairDetailFromID(ctx, id)
				code, _ := failure.CodeOf(err)
				if code != fails.ErrBot {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
				}
			}(strconv.FormatInt(chair.ID, 10))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		q := generateRandomQueryForSearchEstates()
		q.Set("perPage", "10")
		estates, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
			}
			return
		}

		for _, estate := range estates.Estates {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_, err := c.GetEstateDetailFromID(ctx, id)
				code, _ := failure.CodeOf(err)
				if code != fails.ErrBot {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
				}
			}(strconv.FormatInt(estate.ID, 10))
		}
	}()

	wg.Wait()
}
