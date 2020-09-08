package main

import (
	"context"
	"flag"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/reporter"
	"github.com/isucon10-qualify/isucon10-qualify/bench/scenario"
	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
	"github.com/morikuni/failure"

	"github.com/isucon/isucon10-portal/bench-tool.go/benchrun"
)

type Config struct {
	TargetURLStr string
	TargetHost   string

	AllowedIPs []net.IP
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	conf := Config{}
	dataDir := ""
	fixtureDir := ""

	flags.StringVar(&conf.TargetURLStr, "target-url", benchrun.GetTargetAddress(), "target url")
	flags.StringVar(&dataDir, "data-dir", "initial-data", "data directory")
	flags.StringVar(&fixtureDir, "fixture-dir", "../webapp/fixture", "fixture directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		err = failure.Translate(err, fails.ErrBenchmarker, failure.Message("コマンドライン引数のパースに失敗しました"))
		fails.Add(err, fails.ErrorOfInitialize)
		reporter.SetFinished(true)
		reporter.Update(fails.GetMsgs(), 1, 0)
		reporter.Report()
	}

	err = client.SetShareTargetURLs(
		conf.TargetURLStr,
		conf.TargetHost,
	)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrBenchmarker), fails.ErrorOfInitialize)
		reporter.SetFinished(true)
		reporter.Update(fails.GetMsgs(), 1, 0)
		reporter.Report()
	}

	// 初期データの準備
	asset.Initialize(context.Background(), dataDir, fixtureDir)
	msgs := fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.SetFinished(true)
		reporter.Logf("asset initialize failed")
		reporter.Update(msgs, 1, 0)
		reporter.Report()
		return
	}

	reporter.Logf("=== initialize ===")
	// 初期化：/initialize にリクエストを送ることで、外部リソースのURLを指定する・DBのデータを初期データのみにする
	initRes := scenario.Initialize(context.Background())
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.SetFinished(true)
		reporter.Logf("initialize failed")
		reporter.Update(msgs, 1, 0)
		reporter.Report()
		return
	}

	reporter.SetLanguage(initRes.Language)

	reporter.Logf("=== verify ===")
	// 初期チェック：正しく動いているかどうかを確認する
	// 明らかにおかしいレスポンスを返しているアプリケーションはさっさと停止させることで、運営側のリソースを使い果たさない・他サービスへの攻撃に利用されるを防ぐ
	scenario.Verify(context.Background(), dataDir, fixtureDir)
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.SetFinished(true)
		reporter.Logf("verify failed")
		reporter.Update(msgs, 1, 0)
		reporter.Report()
		return
	}

	reporter.Logf("=== validation ===")
	// 一番大切なメイン処理：checkとloadの大きく2つの処理を行う
	// checkはアプリケーションが正しく動いているか常にチェックする
	// 理想的には全リクエストはcheckされるべきだが、それをやるとパフォーマンスが出し切れず、最適化されたアプリケーションよりも遅くなる
	// checkとloadは区別がつかないようにしないといけない。loadのリクエストはログアウト状態しかなかったので、ログアウト時のキャッシュを強くするだけでスコアがはねる問題が過去にあった
	// 今回はほぼ全リクエストがログイン前提になっているので、checkとloadの区別はできないはず
	scenario.Validation(context.Background())

	reporter.Logf("最終的な負荷レベル: %d", score.GetLevel())
	msgs, cCnt, aCnt, _ := fails.Get()
	reporter.Update(msgs, cCnt, aCnt)
	reporter.SetFinished(true)
	reporter.Report()
}
