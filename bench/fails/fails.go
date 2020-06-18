package fails

import (
	"context"
	"log"
	"sync"
	"time"

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
)

type ErrorLabel int

const (
	ErrorOfInitialize ErrorLabel = iota
	ErrorOfEstateSearchScenario
	ErrorOfChairSearchScenario
	ErrorOfEstateNazotteSearchScenario
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
	Msgs           []string
	lastErrorTimes map[ErrorLabel]time.Time

	critical    int
	application int
	trivial     int

	mu sync.RWMutex
}

func NewErrors() *Errors {
	msgs := make([]string, 0, 100)
	times := make(map[ErrorLabel]time.Time)
	return &Errors{
		Msgs:           msgs,
		lastErrorTimes: times,
	}
}

func (e *Errors) GetLastErrorTime(label ErrorLabel) time.Time {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lastErrorTimes[label]
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

	if err == context.DeadlineExceeded || err == context.Canceled {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	log.Printf("%+v", err)

	e.lastErrorTimes[label] = time.Now()

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
