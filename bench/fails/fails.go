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
	// ErrorsForCheck is 基本的にはこっちを使う
	ErrorsForCheck *Errors
	// ErrorsForFinal is 最後のFinal Checkで使う。これをしないとcontext.Canceledのエラーが混ざる
	ErrorsForFinal *Errors
)

func init() {
	ErrorsForCheck = NewErrors()
	ErrorsForFinal = NewErrors()
}

type Errors struct {
	Msgs []string

	critical    int
	application int
	trivial     int

	mu sync.RWMutex
}

func NewErrors() *Errors {
	msgs := make([]string, 0, 100)
	return &Errors{
		Msgs: msgs,
	}
}

func (e *Errors) GetMsgs() (msgs []string) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.Msgs[:]
}

func (e *Errors) Get() (msgs []string, critical, application, trivial int) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.Msgs[:], e.critical, e.application, e.trivial
}

func (e *Errors) Add(err error, label ErrorLabel) {
	if err == nil {
		return
	}

	cause := failure.CauseOf(err)
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	log.Printf("%+v", err)

	msg, ok := failure.MessageOf(err)
	code, _ := failure.CodeOf(err)

	if ok {
		switch code {
		case ErrCritical:
			msg += " (critical error)"
			e.critical++
		case ErrTimeout:
			msg += "（タイムアウトしました）"
			e.trivial++
		case ErrTemporary:
			msg += "（一時的なエラー）"
			e.trivial++
		case ErrApplication:
			e.application++
		case ErrBenchmarker:
			e.Msgs = append(e.Msgs, "運営に連絡してください")
			return
		default:
			e.application++
		}

		e.Msgs = append(e.Msgs, msg)
	} else {
		// 想定外のエラーなのでcritical扱いにしておく
		e.critical++
		e.Msgs = append(e.Msgs, "運営に連絡してください")
	}
}
