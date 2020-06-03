package client

import (
	"context"
	"encoding/json"

	"github.com/morikuni/failure"
)

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

func (c *Client) GetChairDetailFromID(ctx context.Context, id string) (*Chair, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/chair/"+id)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	defer res.Body.Close()

	var chair Chair

	err = json.NewDecoder(res.Body).Decode(&chair)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}
	return &chair, nil
}
