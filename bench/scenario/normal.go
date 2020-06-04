package scenario

import (
	"context"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
)

func initialize(ctx context.Context) error {
	c, err := client.NewClientForInitialize()
	if err != nil {
		return err
	}
	err = c.Initialize(ctx)
	if err != nil {
		return err
	}
	return nil
}
