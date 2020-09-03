package fails

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/morikuni/failure"
)

const (
	// ErrCritical はクリティカルなエラー。少しでも大幅減点・失格になるエラー
	ErrCritical failure.StringCode = "error critical"
	// ErrApplication はアプリケーションの挙動でおかしいエラー。Verify時は1つでも失格。Validation時は一定数以上で失格
	ErrApplication failure.StringCode = "error application"
	// ErrTimeout はタイムアウトエラー。基本は大目に見る。
	ErrTimeout failure.StringCode = "error timeout"
	// ErrTemporary は一時的なエラー。基本は大目に見る。
	ErrTemporary failure.StringCode = "error temporary"
	// ErrBenchmarker はベンチマーカ側のエラー。基本的には運営に連絡してもらう
	ErrBenchmarker failure.StringCode = "error benchmarker"
	// ErrBot はBotによるリクエストによって発生したエラー。
	ErrBot failure.StringCode = "error bot"
)

type ErrorLabel int

const (
	ErrorOfInitialize ErrorLabel = iota
	ErrorOfVerify
	ErrorOfEstateSearchScenario
	ErrorOfChairSearchScenario
	ErrorOfEstateNazotteSearchScenario
	ErrorOfBotScenario
	ErrorOfEstateDraftPostScenario
	ErrorOfChairDraftPostScenario
)

var (
	msgs []string

	critical    int
	application int
	trivial     int

	mu sync.RWMutex
)

func init() {
	msgs = make([]string, 0, 100)
}

func GetMsgs() (msgs []string) {
	mu.RLock()
	defer mu.RUnlock()

	return msgs[:]
}

func Get() (msgs []string, critical, application, trivial int) {
	mu.RLock()
	defer mu.RUnlock()

	return msgs[:], critical, application, trivial
}

func Add(err error, label ErrorLabel) {
	if err == nil {
		return
	}

	cause := failure.CauseOf(err)
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	log.Printf("%+v", err)

	msg, ok := failure.MessageOf(err)
	code, _ := failure.CodeOf(err)

	if ok {
		switch code {
		case ErrCritical:
			msg += " (critical error)"
			critical++
		case ErrTimeout:
			msg += "（タイムアウトしました）"
			trivial++
		case ErrTemporary:
			msg += "（一時的なエラー）"
			trivial++
		case ErrApplication:
			application++
		case ErrBenchmarker:
			msgs = append(msgs, "運営に連絡してください")
			return
		default:
			application++
		}

		msgs = append(msgs, msg)
	} else {
		// 想定外のエラーなのでcritical扱いにしておく
		critical++
		msgs = append(msgs, "運営に連絡してください")
	}
}
