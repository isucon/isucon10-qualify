package scenario

import (
	"context"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

const (
	initializeTimeout = 180 * time.Second
	loadTimeout       = 60 * time.Second
)

func Initialize(ctx context.Context) {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	// timeoutSeconds := 180

	ctx, cancel := context.WithTimeout(ctx, initializeTimeout)
	defer cancel()

	err := initialize(ctx)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfInitialize)
	}
}

func Validation(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, loadTimeout)
	defer cancel()
	go Load(cancelCtx)
	<-cancelCtx.Done()
}
