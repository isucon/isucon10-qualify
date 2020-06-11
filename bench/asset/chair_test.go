package asset

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
)

func Test_ParallelStockDecrement(t *testing.T) {
	initialStock := int64(1000000)
	c := Chair{
		stock: initialStock,
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 100; j++ {
				c.DecrementStock()
			}
		}()
	}
	close(start)
	wg.Wait()
	got := c.GetStock()
	expected := initialStock - 100*100
	if got != expected {
		t.Errorf("unexpected stocks. expected: %v, but got: %v", expected, got)
	}
}

func Test_ParallelViewCountIncrement(t *testing.T) {
	type Incrementable interface {
		GetViewCount() int64
		IncrementViewCount()
	}

	initialViewCount := int64(100)

	testAsset := []struct {
		TestName string
		Asset    Incrementable
	}{
		{
			TestName: "Test Chair",
			Asset: &Chair{
				viewCount: initialViewCount,
			},
		},
		{
			TestName: "Test Estate",
			Asset: &Estate{
				viewCount: initialViewCount,
			},
		},
	}

	for _, tc := range testAsset {
		t.Run(tc.TestName, func(t *testing.T) {
			var wg sync.WaitGroup
			start := make(chan struct{})
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					for j := 0; j < 100; j++ {
						tc.Asset.IncrementViewCount()
					}
				}()
			}
			close(start)
			wg.Wait()
			got := tc.Asset.GetViewCount()
			expected := initialViewCount + 100*100
			if got != expected {
				t.Errorf("unexpected stocks. expected: %v, but got: %v", expected, got)
			}
		})
	}
}

func TestChair_MarshalJSON(t *testing.T) {
	chair := Chair{
		ID:          1,
		Name:        "name",
		Description: "description",
		Thumbnail:   "thumbnail",
		Price:       2,
		Height:      3,
		Width:       4,
		Depth:       5,
		Color:       "color",
		Features:    "features",
		Kind:        "kind",
		stock:       6,
		viewCount:   7,
	}
	b, err := json.Marshal(chair)
	if err != nil {
		t.Fatal("failed to marshal json:", err)
	}
	var got Chair
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal("failed to unmarshal json:", err)
	}
	if !reflect.DeepEqual(chair, got) {
		t.Errorf("unexpected chair. expected: %+v, but got: %+v", chair, got)
	}
}
