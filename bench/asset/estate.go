package asset

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type JSONEstate struct {
	ID          int64   `json:"id"`
	Thumbnail   string  `json:"thumbnail"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DoorHeight  int64   `json:"doorHeight"`
	DoorWidth   int64   `json:"doorWidth"`
	ViewCount   int64   `json:"viewCount"`
	Rent        int64   `json:"rent"`
	Features    string  `json:"features"`
}

type Estate struct {
	ID          int64
	Thumbnail   string
	Name        string
	Description string
	Address     string
	Latitude    float64
	Longitude   float64
	DoorHeight  int64
	DoorWidth   int64
	Rent        int64
	Features    string

	viewCount int64
}

func (e Estate) MarshalJSON() ([]byte, error) {

	m := JSONEstate{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Thumbnail:   e.Thumbnail,
		Rent:        e.Rent,
		Address:     e.Address,
		Latitude:    e.Latitude,
		Longitude:   e.Longitude,
		DoorHeight:  e.DoorHeight,
		DoorWidth:   e.DoorWidth,
		Features:    e.Features,
		ViewCount:   e.viewCount,
	}

	return json.Marshal(m)
}

func (e *Estate) UnmarshalJSON(data []byte) error {
	var je JSONEstate

	err := json.Unmarshal(data, &je)
	if err != nil {
		fmt.Println(err)
		return err
	}

	e.ID = je.ID
	e.Name = je.Name
	e.Description = je.Description
	e.Thumbnail = je.Thumbnail
	e.Rent = je.Rent
	e.Address = je.Address
	e.DoorHeight = je.DoorHeight
	e.DoorWidth = je.DoorWidth
	e.Latitude = je.Latitude
	e.Longitude = je.Longitude
	e.Features = je.Features
	e.viewCount = je.ViewCount

	return nil
}

func (e1 *Estate) Equal(e2 *Estate) bool {
	return e1.ID == e2.ID
}

func (e *Estate) GetViewCount() int64 {
	return atomic.LoadInt64(&(e.viewCount))
}

func (e *Estate) IncrementViewCount() {
	atomic.AddInt64(&(e.viewCount), 1)
}
