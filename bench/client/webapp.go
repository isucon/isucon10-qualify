package client

import (
	"context"
	"encoding/json"

	"github.com/morikuni/failure"
	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func (c *Client) GetChairDetailFromID(ctx context.Context, id string) (*asset.Chair, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/chair/"+id)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	defer res.Body.Close()

	var chair asset.Chair

	err = json.NewDecoder(res.Body).Decode(&chair)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}
	return &chair, nil
}
