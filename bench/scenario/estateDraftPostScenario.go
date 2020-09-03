package scenario

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func loadEstatesFromJSON(ctx context.Context, filePath string) ([]asset.Estate, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	estates := []asset.Estate{}
	decoder := json.NewDecoder(f)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var estate asset.Estate
		if err := decoder.Decode(&estate); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		estates = append(estates, estate)
	}

	return estates, nil
}

func estateDraftPostScenario(ctx context.Context, c *client.Client, filePath string) {
	estates, err := loadEstatesFromJSON(ctx, filePath)
	err = c.PostEstates(ctx, estates)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical), fails.ErrorOfEstateDraftPostScenario)
	}
}
