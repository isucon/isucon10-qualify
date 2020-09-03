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

func loadChairsFromJSON(ctx context.Context, filePath string) ([]asset.Chair, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	chairs := []asset.Chair{}
	decoder := json.NewDecoder(f)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var chair asset.Chair
		if err := decoder.Decode(&chair); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chairs = append(chairs, chair)
	}

	return chairs, nil
}

func chairDraftPostScenario(ctx context.Context, c *client.Client, filePath string) {
	chairs, err := loadChairsFromJSON(ctx, filePath)
	err = c.PostChairs(ctx, chairs)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical), fails.ErrorOfChairDraftPostScenario)
	}
}
