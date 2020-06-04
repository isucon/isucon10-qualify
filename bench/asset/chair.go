package asset

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
	Stock       int64  `json:"stock"`
}

func (c1 *Chair) Equal(c2 *Chair) bool {
	return c1.ID == c2.ID
}
