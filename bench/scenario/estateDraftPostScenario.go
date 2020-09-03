package scenario

import (
	"context"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func estateDraftPostScenario(ctx context.Context, c *client.Client, filePath string) {
	err := c.PostEstates(ctx, filePath)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical), fails.ErrorOfEstateDraftPostScenario)
	}
}
