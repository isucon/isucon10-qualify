package scenario

import (
	"context"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func chairDraftPostScenario(ctx context.Context, c *client.Client, filePath string) {
	err := c.PostChairs(ctx, filePath)
	if err != nil {
		fails.ErrorsForCheck.Add(failure.Translate(err, fails.ErrCritical), fails.ErrorOfChairDraftPostScenario)
	}
}
