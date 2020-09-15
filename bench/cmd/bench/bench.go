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
	defer func() {
		reporter.SetFinished(true)
		reporter.Report(fails.Get())
	}()

	defer func() {
		err := recover()
		if err, ok := err.(error); ok {
			err = failure.Translate(err, fails.ErrBenchmarker)
			fails.Add(err)
		}
	}()

	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	conf := Config{}
	dataDir := ""
	fixtureDir := ""

	flags.StringVar(&conf.TargetURLStr, "target-url", "http://" + benchrun.GetTargetAddress(), "target url")
	flags.StringVar(&dataDir, "data-dir", "../initial-data", "data directory")
	flags.StringVar(&fixtureDir, "fixture-dir", "../webapp/fixture", "fixture directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		err = failure.Translate(err, fails.ErrBenchmarker, failure.Message("コマンドライン引数のパースに失敗しました"))
		fails.Add(err)
		reporter.SetPassed(false)
		reporter.SetReason("コマンドライン引数のパースに失敗しました")
		return
	}

	err = client.SetShareTargetURLs(
		conf.TargetURLStr,
		conf.TargetHost,
	)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrBenchmarker))
		reporter.SetPassed(false)
		reporter.SetReason("ベンチ対象サーバーのURLが不正です")
		return
	}

	// 初期データの準備
	asset.Initialize(context.Background(), dataDir, fixtureDir)
	msgs := fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.Logf("asset initialize failed")
		reporter.SetPassed(false)
		reporter.SetReason("ベンチマーカーの初期化に失敗しました")
		return
	}

	reporter.Logf("=== initialize ===")
	// 初期化：/initialize にリクエストを送ることで、外部リソースのURLを指定する・DBのデータを初期データのみにする
	initRes := scenario.Initialize(context.Background())
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.Logf("initialize failed")
		reporter.SetPassed(false)
		reporter.SetReason("POST /initializeに失敗しました")
		return
	}

	reporter.SetLanguage(initRes.Language)

	reporter.Logf("=== verify ===")
	// 初期チェック：正しく動いているかどうかを確認する
	// 明らかにおかしいレスポンスを返しているアプリケーションはさっさと停止させることで、運営側のリソースを使い果たさない・他サービスへの攻撃に利用されるを防ぐ
	scenario.Verify(context.Background(), dataDir, fixtureDir)
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		reporter.Logf("verify failed")
		reporter.SetPassed(false)
		reporter.SetReason("アプリケーション互換性チェックに失敗しました")
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

	// ベンチマーク終了時にcritical errorが1つ以上、もしくはapplication errorが10回以上で失格
	msgs, critical, application, _ := fails.Get()
	isPassed := true

	if critical > 0 {
		isPassed = false
		reporter.SetReason("致命的なエラーが発生しました")
	} else if application >= 10 {
		isPassed = false
		reporter.SetReason("アプリケーションエラーが10回以上発生しました")
	} else {
		reporter.SetReason("OK")
	}

	reporter.SetPassed(isPassed)
}
