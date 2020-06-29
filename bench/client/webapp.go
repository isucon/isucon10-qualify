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
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/passes"
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
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("Initialize client.newGetRequest error occured"))
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
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("GetChairDerailFromID client.newGetRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		io.Copy(ioutil.Discard, res.Body)
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
	passes.AddDuration(time.Since(t), passes.LabelOfGetChairDetailFromID)

	return &chair, nil
}

func (c *Client) SearchChairsWithQuery(ctx context.Context, q url.Values) (*ChairsResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/chair/search", q)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("SearchChairsWithQuery client.newGetRequestWithQuery error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfSearchChairsWithQuery)

	return &chairs, nil
}

func (c *Client) SearchEstatesWithQuery(ctx context.Context, q url.Values) (*EstatesResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/estate/search", q)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("SearchEstatesWithQuery client.newGetRequestWithQuery error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estates EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estates)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfSearchEstatesWithQuery)

	return &estates, nil
}

func (c *Client) SearchEstatesNazotte(ctx context.Context, polygon *Coordinates) (*EstatesResponse, error) {
	b, err := json.Marshal(polygon)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("SearchEstatesNazotte json.Marshal error occured"))
	}

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/nazotte", bytes.NewBuffer(b))
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("SearchEstatesNazotte client.newPostRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estates EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estates)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfSearchEstatesNazotte)

	return &estates, nil
}

func (c *Client) GetEstateDetailFromID(ctx context.Context, id string) (*asset.Estate, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/estate/"+id)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("GetEstateDetailFromID client.newGetRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estate asset.Estate

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementEstateViewCount(estate.ID)
	passes.AddDuration(time.Since(t), passes.LabelOfGetEstateDetailFromID)

	return &estate, nil
}

func (c *Client) GetRecommendedChair(ctx context.Context) (*ChairsResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_chair")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("GetRecommendedChair client.newGetRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_chair: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_chair: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfGetRecommendedChair)

	return &chairs, nil
}

func (c *Client) GetRecommendedEstate(ctx context.Context) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("GetRecommendedEstate client.newGetRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estate EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfGetRecommendedEstate)

	return &estate, nil
}

func (c *Client) GetRecommendedEstatesFromChair(ctx context.Context, id int64) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate/"+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker, failure.Message("GetRecommendedEstatesFromChair client.newGetRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var estate EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: JSONデコードに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfGetRecommendedEstatesFromChair)

	return &estate, nil
}

func (c *Client) BuyChair(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair/buy/"+id, nil)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("BuyChair client.newPostRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/chair/buy/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	err = checkStatusCode(res, 200)
	if err != nil {
		return failure.Wrap(err, failure.Message("POST /api/chair/buy/:id: リクエストに失敗しました"))
	}

	intid, _ := strconv.ParseInt(id, 10, 64)
	asset.DecrementChairStock(intid)
	passes.AddDuration(time.Since(t), passes.LabelOfBuyChair)

	return nil
}

func (c *Client) RequestEstateDocument(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/req_doc/"+id, nil)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("RequestEstateDocument client.newPostRequest error occured"))
	}

	req = req.WithContext(ctx)
	t := time.Now()
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	err = checkStatusCode(res, 200)
	if err != nil {
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}

	passes.AddDuration(time.Since(t), passes.LabelOfRequestEstateDocument)

	return nil
}
