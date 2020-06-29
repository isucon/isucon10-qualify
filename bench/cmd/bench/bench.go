package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/passes"
	"github.com/isucon10-qualify/isucon10-qualify/bench/scenario"
)

type Output struct {
	Pass     bool     `json:"pass"`
	Score    int      `json:"score"`
	Messages []string `json:"messages"`
}

type Config struct {
	TargetURLStr string
	TargetHost   string

	AllowedIPs []net.IP
}

func init() {
	rand.Seed(time.Now().UnixNano())

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	conf := Config{}
	dataDir := ""

	flags.StringVar(&conf.TargetURLStr, "target-url", "http://127.0.0.1:8000", "target url")
	flags.StringVar(&dataDir, "data-dir", "initial-data", "data directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = client.SetShareTargetURLs(
		conf.TargetURLStr,
		conf.TargetHost,
	)
	if err != nil {
		log.Fatal(err)
	}

	client.InitializeClients()

	// 初期データの準備
	asset.Initialize(dataDir)

	log.Print("=== initialize ===")
	// 初期化：/initialize にリクエストを送ることで、外部リソースのURLを指定する・DBのデータを初期データのみにする
	scenario.Initialize(context.Background())
	eMsgs := fails.ErrorsForCheck.GetMsgs()
	if len(eMsgs) > 0 {
		log.Print("initialize failed")

		output := Output{
			Pass:     false,
			Score:    0,
			Messages: eMsgs,
		}
		json.NewEncoder(os.Stdout).Encode(output)

		return
	}

	log.Print("=== verify ===")
	// 初期チェック：正しく動いているかどうかを確認する
	// 明らかにおかしいレスポンスを返しているアプリケーションはさっさと停止させることで、運営側のリソースを使い果たさない・他サービスへの攻撃に利用されるを防ぐ
	scenario.Verify(context.Background())
	eMsgs = fails.ErrorsForCheck.GetMsgs()
	if len(eMsgs) > 0 {
		log.Print("verify failed")

		output := Output{
			Pass:     false,
			Score:    0,
			Messages: eMsgs,
		}
		json.NewEncoder(os.Stdout).Encode(output)

		return
	}

	log.Print("=== validation ===")
	// 一番大切なメイン処理：checkとloadの大きく2つの処理を行う
	// checkはアプリケーションが正しく動いているか常にチェックする
	// 理想的には全リクエストはcheckされるべきだが、それをやるとパフォーマンスが出し切れず、最適化されたアプリケーションよりも遅くなる
	// checkとloadは区別がつかないようにしないといけない。loadのリクエストはログアウト状態しかなかったので、ログアウト時のキャッシュを強くするだけでスコアがはねる問題が過去にあった
	// 今回はほぼ全リクエストがログイン前提になっているので、checkとloadの区別はできないはず
	scenario.Validation(context.Background())

	// context.Canceledのエラーは直後に取れば基本的には入ってこない
	eMsgs, cCnt, aCnt, _ := fails.ErrorsForCheck.Get()
	// critical errorは1つでもあれば、application errorは10回以上で失格
	if cCnt > 0 || aCnt >= 10 {
		log.Print("cause error!")

		output := Output{
			Pass:     false,
			Score:    0,
			Messages: uniqMsgs(eMsgs),
		}
		json.NewEncoder(os.Stdout).Encode(output)

		return
	}

	score := 0
	score += 1 * passes.GetCount(passes.LabelOfGetChairDetailFromID)
	score += 1 * passes.GetCount(passes.LabelOfGetEstateDetailFromID)
	score += 1 * passes.GetCount(passes.LabelOfSearchChairsWithQuery)
	score += 1 * passes.GetCount(passes.LabelOfSearchEstatesWithQuery)
	score += 1 * passes.GetCount(passes.LabelOfSearchEstatesNazotte)
	score += 1 * passes.GetCount(passes.LabelOfGetRecommendedEstatesFromChair)
	score += 1 * passes.GetCount(passes.LabelOfBuyChair)
	score += 1 * passes.GetCount(passes.LabelOfRequestEstateDocument)

	// application errorは1回で10点減点
	penalty := 10 * aCnt
	log.Print(score, penalty)

	score -= penalty
	// 0点以下なら失格
	if score <= 0 {
		output := Output{
			Pass:     false,
			Score:    0,
			Messages: uniqMsgs(eMsgs),
		}
		json.NewEncoder(os.Stdout).Encode(output)

		return
	}

	output := Output{
		Pass:     true,
		Score:    score,
		Messages: uniqMsgs(eMsgs),
	}
	json.NewEncoder(os.Stdout).Encode(output)
}

func uniqMsgs(allMsgs []string) []string {
	sort.Strings(allMsgs)
	msgs := make([]string, 0, len(allMsgs))

	tmp := ""

	// 適当にuniqする
	for _, m := range allMsgs {
		if tmp != m {
			tmp = m
			msgs = append(msgs, m)
		}
	}

	return msgs
}
