package asset

import (
	"encoding/json"
	"fmt"
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
	Popularity  int64   `json:"popularity"`
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

	popularity int64
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
		Popularity:  e.popularity,
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
	e.popularity = je.Popularity

	return nil
}

func (e1 *Estate) Equal(e2 *Estate) bool {
	return e1.ID == e2.ID &&
		e1.Name == e2.Name &&
		e1.Description == e2.Description &&
		e1.Thumbnail == e2.Thumbnail &&
		e1.Rent == e2.Rent &&
		e1.Address == e2.Address &&
		e1.DoorHeight == e2.DoorHeight &&
		e1.DoorWidth == e2.DoorWidth &&
		e1.Latitude == e2.Latitude &&
		e1.Longitude == e2.Longitude &&
		e1.Features == e2.Features
}

func (e *Estate) GetPopularity() int64 {
	return e.popularity
}
