package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/conversion"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

type Coordinates struct {
	Coordinates []*Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *Client) Initialize(ctx context.Context) error {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/initialize")
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	// T/O付きのコンテキストが入る
	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil {
		return failure.Wrap(err, failure.Message("GET /initialize: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	// MEMO: /initializeの成功ステータスによって第二引数が変わる可能性がある
	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		return failure.Wrap(err, failure.Message("GET /initialize: レスポンスコードが不正です"))
	}

	return nil
}

type ChairsResponse struct {
	Count  int64         `json:"count"`
	Chairs []asset.Chair `json:"chairs"`
}

type EstatesResponse struct {
	Count   int64          `json:"count"`
	Estates []asset.Estate `json:"estates"`
}

func (c *Client) GetChairDetailFromID(ctx context.Context, id string) (*asset.Chair, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/chair/"+id)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK, http.StatusNotFound})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: レスポンスコードが不正です"))
	}

	if res.StatusCode == http.StatusNotFound {
		_, err = io.Copy(ioutil.Discard, res.Body)
		if err != nil {
			return nil, failure.Translate(err, fails.ErrBenchmarker)
		}
		return nil, nil
	}

	var chair asset.Chair

	err = json.NewDecoder(res.Body).Decode(&chair)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementChairViewCount(chair.ID)

	return &chair, nil
}

func (c *Client) SearchChairsWithQuery(ctx context.Context, q url.Values) (*ChairsResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/chair/search", q)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: レスポンスコードが不正です"))
	}

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: JSONデコードに失敗しました"))
	}

	return &chairs, nil
}

func (c *Client) SearchEstatesWithQuery(ctx context.Context, q url.Values) (*EstatesResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/estate/search", q)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: レスポンスコードが不正です"))
	}

	var estates EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estates)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: JSONデコードに失敗しました"))
	}

	return &estates, nil
}

func (c *Client) SearchEstatesNazotte(ctx context.Context, polygon *Coordinates) (*EstatesResponse, error) {
	b, err := json.Marshal(polygon)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/nazotte", bytes.NewBuffer(b))
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: レスポンスコードが不正です"))
	}

	var estates EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estates)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: JSONデコードに失敗しました"))
	}

	return &estates, nil
}

func (c *Client) GetEstateDetailFromID(ctx context.Context, id string) (*asset.Estate, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/estate/"+id)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: レスポンスコードが不正です"))
	}

	var estate asset.Estate

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementEstateViewCount(estate.ID)

	return &estate, nil
}

func (c *Client) GetRecommendedChair(ctx context.Context) (*ChairsResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_chair")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_chair: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_chair: レスポンスコードが不正です"))
	}

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_chair: JSONデコードに失敗しました"))
	}

	return &chairs, nil
}

func (c *Client) GetRecommendedEstate(ctx context.Context) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate: レスポンスコードが不正です"))
	}

	var estate EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate: JSONデコードに失敗しました"))
	}

	return &estate, nil
}

func (c *Client) GetRecommendedEstatesFromChair(ctx context.Context, id int64) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate/"+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: レスポンスコードが不正です"))
	}

	var estate EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: JSONデコードに失敗しました"))
	}

	return &estate, nil
}

func (c *Client) BuyChair(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair/buy/"+id, nil)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/chair/buy/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return failure.Translate(err, fails.ErrBot)
		}
		return failure.Wrap(err, failure.Message("POST /api/chair/buy/:id: リクエストに失敗しました"))
	}

	intid, _ := strconv.ParseInt(id, 10, 64)
	asset.DecrementChairStock(intid)
	if !c.isBot {
		conversion.IncrementCount()
	}

	return nil
}

func (c *Client) RequestEstateDocument(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/req_doc/"+id, nil)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return failure.Translate(err, fails.ErrBot)
		}
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}

	if !c.isBot {
		conversion.IncrementCount()
	}

	return nil
}
