package main

type RequestFrame struct {
	RequectContent Request `json:"request"`
}
type Request struct {
	Method string      `json:"method"`
	Uri    string      `json:"uri"`
	ID     string      `json:"id"`
	Query  Queries     `json:"query"`
	Body   Coordinates `json:"body"`
}
type Coordinates struct {
	Coordinate []Coordinate `json:"coordinates"`
}
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type EstatesAnswerJson struct {
	Req Request     `json:"request"`
	Res EstatesBody `json:"response"`
}
type EstatesBody struct {
	Body EstateResponse `json:"body"`
}
type ChairsAnswerJson struct {
	Req Request    `json:"request"`
	Res ChairsBody `json:"response"`
}
type ChairsBody struct {
	Body ChairsResponse `json:"body"`
}
type EstateResponse struct {
	Count   int64    `json:"count"`
	Estates []Estate `json:"estates"`
}
type ChairsResponse struct {
	Count  int64   `json:"count"`
	Chairs []Chair `json:"chairs"`
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

type Queries struct {
	RentRangeId       string `json:"rentRangeId"`
	PriceRangeId      string `json:"priceRangeId"`
	DoorHeightRangeId string `json:"doorHeightRangeId"`
	DoorWidthRangeId  string `json:"doorWidthRangeId"`
	HeightRangeId     string `json:"heightRangeId"`
	WidthRangeId      string `json:"widthRangeId"`
	DepthRangeId      string `json:"depthRangeId"`
	Features          string `json:"features"`
	Kind              string `json:"kind"`
	Color             string `json:"color"`
	Page              string `json:"page"`
	PerPage           string `json:"perPage"`
}
