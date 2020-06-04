package asset

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
	ViewCount   int64   `json:"viewCount"`
	Rent        int64   `json:"rent"`
	Features    string  `json:"features"`
}

func (e1 *Estate) Equal(e2 *Estate) bool {
	return e1.ID == e2.ID
}
