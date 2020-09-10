package scenario

import (
	"context"
	"log"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

func Initialize(ctx context.Context) *client.InitializeResponse {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	// timeoutSeconds := 180

	ctx, cancel := context.WithTimeout(ctx, parameter.InitializeTimeout)
	defer cancel()

	res, err := initialize(ctx)
	if err != nil {
		fails.Add(err, fails.ErrorOfInitialize)
	}
	return res
}

func Validation(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, parameter.LoadTimeout)
	defer cancel()
	go Load(cancelCtx)

	select {
	case <-fails.Fail():
		log.Println("fail条件を満たしました")
	case <-cancelCtx.Done():
	}
}
