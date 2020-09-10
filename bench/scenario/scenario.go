package scenario

import (
	"context"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/isucon10-qualify/isucon10-qualify/bench/reporter"
	"github.com/morikuni/failure"
)

func Initialize(ctx context.Context) *client.InitializeResponse {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	// timeoutSeconds := 180

	ctx, cancel := context.WithTimeout(ctx, parameter.InitializeTimeout)
	defer cancel()

	res, err := initialize(ctx)
	if err != nil {
		if ctx.Err() != nil {
			err = failure.New(fails.ErrCritical, failure.Message("POST /initialize: リクエストがタイムアウトしました"))
			fails.Add(err, fails.ErrorOfInitialize)
		} else {
			fails.Add(err, fails.ErrorOfInitialize)
		}
	}
	return res
}

func Validation(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, parameter.LoadTimeout)
	defer cancel()
	go Load(cancelCtx)

	for {
		t := time.NewTimer(parameter.ReportInterval)
		select {
		case <-t.C:
			reporter.Report(fails.Get())
		case <-fails.Fail():
			reporter.Logf("fail条件を満たしました")
			return
		case <-cancelCtx.Done():
			return
		}
	}
}
