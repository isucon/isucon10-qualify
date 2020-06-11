package client

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/morikuni/failure"
)

func (c *Client) Initialize(ctx context.Context) error {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/initialize")
	if err != nil {
		return failure.Wrap(err, failure.Message("GET /initialize: リクエストに失敗しました"))
	}

	// T/O付きのコンテキストが入る
	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil {
		return failure.Wrap(err, failure.Message("GET /initialize: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	// MEMO: /initializeの成功ステータスによって第二引数が変わる可能性がある
	err = checkStatusCode(res, http.StatusOK)
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, res.Body)

	return nil
}

type ChairsResponse struct {
	Chairs []asset.Chair
}

type EstatesResponse struct {
	Estates []asset.Estate
}

func (c *Client) GetChairDetailFromID(ctx context.Context, id string) (*asset.Chair, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/chair/"+id)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	defer res.Body.Close()

	var chair asset.Chair
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&chair)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementChairViewCount(chair.ID)

	return &chair, nil
}

func (c *Client) SearchChairsWithQuery(ctx context.Context, q url.Values) (*ChairsResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/chair/search", q)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/search: Query: リクエストに失敗しました"))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil || res == nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/search: Query: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: JSONデコードに失敗しました"))
	}
	return &chairs, nil
}

func (c *Client) SearchEstatesWithQuery(ctx context.Context, q url.Values) (*EstatesResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/estate/search", q)

	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/estate/search: Query: リクエストに失敗しました"))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil || res == nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/chair/search: Query: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estates EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estates)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: JSONデコードに失敗しました"))
	}
	return &estates, nil
}

func (c *Client) GetEstateDetailFromID(ctx context.Context, id string) (*asset.Estate, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/estate/"+id)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/estate/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	defer res.Body.Close()

	var estate asset.Estate

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementEstateViewCount(estate.ID)

	return &estate, nil
}

func (c *Client) GetRecommendedEstatesFromChair(ctx context.Context, id int64) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate/"+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("GET /api/recommended_estate/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	defer res.Body.Close()

	var estate EstatesResponse
	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: JSONデコードに失敗しました"))
	}

	return &estate, nil
}

func (c *Client) BuyChair(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair/buy/"+id, nil)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/chair/buy/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/chair/buy/%v: リクエストに失敗しました", id))
	}

	err = checkStatusCode(res, 200)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/chair/buy/%v: リクエストに失敗しました", id))
	}

	intid, _ := strconv.ParseInt(id, 10, 64)
	asset.DecrementChairStock(intid)

	return nil
}

func (c *Client) RequestEstateDocument(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/req_doc/"+id, nil)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/estate/req_doc/%v: リクエストに失敗しました", id))
	}

	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/estate/req_doc/%v: リクエストに失敗しました", id))
	}

	err = checkStatusCode(res, 200)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("POST /api/estate/req_doc/%v: リクエストに失敗しました", id))
	}
	return nil
}
