package scenario

import (
	"context"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/paramater"
)

func botScenario(ctx context.Context, c *client.Client) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
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

		_, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
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
		q.Set("perPage", strconv.Itoa(paramater.PerPageOfEstateSearch))
		q.Set("page", "0")

		_, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfBotScenario)
		}
	}()

	wg.Wait()
}
