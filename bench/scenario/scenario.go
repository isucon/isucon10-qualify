package scenario

import (
	"context"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/paramater"
)

func Initialize(ctx context.Context) (*client.InitializeResponse) {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	// timeoutSeconds := 180

	ctx, cancel := context.WithTimeout(ctx, paramater.InitializeTimeout)
	defer cancel()

	res, err := initialize(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfInitialize)
	}
	return res
}

func Validation(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, paramater.LoadTimeout)
	defer cancel()
	go Load(cancelCtx)
	<-cancelCtx.Done()
}
