package asset

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type JSONChair struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Price       int64  `json:"price"`
	Height      int64  `json:"height"`
	Width       int64  `json:"width"`
	Depth       int64  `json:"depth"`
	Color       string `json:"color"`
	Features    string `json:"features"`
	ViewCount   int64  `json:"view_count"`
	Kind        string `json:"kind"`
	Stock       int64  `json:"stock"`
}

type Chair struct {
	ID          int64
	Name        string
	Description string
	Thumbnail   string
	Price       int64
	Height      int64
	Width       int64
	Depth       int64
	Color       string
	Features    string
	Kind        string

	viewCount   int64
	stock       int64
	soldOutTime atomic.Value
}

func (c Chair) MarshalJSON() ([]byte, error) {

	m := JSONChair{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Thumbnail:   c.Thumbnail,
		Price:       c.Price,
		Height:      c.Height,
		Width:       c.Width,
		Depth:       c.Depth,
		Color:       c.Color,
		Features:    c.Features,
		ViewCount:   c.viewCount,
		Kind:        c.Kind,
		Stock:       c.stock,
	}

	return json.Marshal(m)
}

func (c *Chair) UnmarshalJSON(data []byte) error {
	var jc JSONChair

	err := json.Unmarshal(data, &jc)
	if err != nil {
		return err
	}

	c.ID = jc.ID
	c.Name = jc.Name
	c.Description = jc.Description
	c.Thumbnail = jc.Thumbnail
	c.Price = jc.Price
	c.Height = jc.Height
	c.Width = jc.Width
	c.Depth = jc.Depth
	c.Color = jc.Color
	c.Features = jc.Features
	c.viewCount = jc.ViewCount
	c.Kind = jc.Kind
	c.stock = jc.Stock
	c.soldOutTime = atomic.Value{}

	return nil
}

func (c1 *Chair) Equal(c2 *Chair) bool {
	return c1.ID == c2.ID &&
		c1.Name == c2.Name &&
		c1.Description == c2.Description &&
		c1.Thumbnail == c2.Thumbnail &&
		c1.Price == c2.Price &&
		c1.Height == c2.Height &&
		c1.Width == c2.Width &&
		c1.Depth == c2.Depth &&
		c1.Color == c2.Color &&
		c1.Features == c2.Features &&
		c1.Kind == c2.Kind
}

func (c *Chair) GetViewCount() int64 {
	return atomic.LoadInt64(&(c.viewCount))
}

func (c *Chair) GetStock() int64 {
	return atomic.LoadInt64(&(c.stock))
}

func (c *Chair) IncrementViewCount() {
	atomic.AddInt64(&(c.viewCount), 1)
}

func (c *Chair) DecrementStock() {
	stock := atomic.AddInt64(&(c.stock), -1)
	if stock == 0 {
		c.soldOutTime.Store(time.Now())
	}
}

func (c *Chair) GetSoldOutTime() *time.Time {
	value := c.soldOutTime.Load()
	if value == nil {
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return nil
	}
	return &t
}
