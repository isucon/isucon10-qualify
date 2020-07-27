package conversion

import (
	"sync/atomic"
)

var (
	count int64
)

func IncrementCount() {
	atomic.AddInt64(&count, 1)
}

func GetCount() int64 {
	return atomic.LoadInt64(&count)
}
