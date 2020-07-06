package main

type Request struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	Query    string `json:"query"`
	Body     string `json:"body"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Snapshot struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Chair struct {
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
	Kind        string `json:"kind"`
}

type Estate struct {
	ID          int64   `json:"id"`
	Thumbnail   string  `json:"thumbnail"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DoorHeight  int64   `json:"doorHeight"`
	DoorWidth   int64   `json:"doorWidth"`
	Rent        int64   `json:"rent"`
	Features    string  `json:"features"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type NazotteRequestBody struct {
	Coordinates []Coordinate `json:"coordinates"`
}

type ChairsResponse struct {
	Count  int64   `json:"count"`
	Chairs []Chair `json:"chairs"`
}

type EstatesResponse struct {
	Count   int64    `json:"count"`
	Estates []Estate `json:"estates"`
}
