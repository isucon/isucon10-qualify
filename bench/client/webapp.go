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

const (
	timeoutOfGetChairDetailFromID           = 300 * time.Millisecond
	timeoutOfSearchChairsWithQuery          = 500 * time.Millisecond
	timeoutOfGetEstateDetailFromID          = 300 * time.Millisecond
	timeoutOfSearchEstatesWithQuery         = 500 * time.Millisecond
	timeoutOfSearchEstatesNazotte           = 1000 * time.Millisecond
	timeoutOfGetRecommendedEstatesFromChair = 300 * time.Millisecond
	timeoutOfBuyChair                       = 300 * time.Millisecond
	timeoutOfRequestEstateDocument          = 300 * time.Millisecond
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
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfGetChairDetailFromID)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/chair/:id: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	var chair asset.Chair

	err = json.NewDecoder(res.Body).Decode(&chair)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/chair/:id: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementChairViewCount(chair.ID)
	passes.IncrementCount(passes.LabelOfGetChairDetailFromID)

	return &chair, nil
}

func (c *Client) SearchChairsWithQuery(ctx context.Context, q url.Values) (*ChairsResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/chair/search", q)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfSearchChairsWithQuery)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/chair/search: リクエストがタイムアウトしました"))
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
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/chair/search: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: JSONデコードに失敗しました"))
	}

	passes.IncrementCount(passes.LabelOfSearchChairsWithQuery)

	return &chairs, nil
}

func (c *Client) SearchEstatesWithQuery(ctx context.Context, q url.Values) (*EstatesResponse, error) {
	req, err := c.newGetRequestWithQuery(ShareTargetURLs.AppURL, "/api/estate/search", q)

	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfSearchEstatesWithQuery)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/estate/search: リクエストがタイムアウトしました"))
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
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/estate/search: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search: JSONデコードに失敗しました"))
	}

	passes.IncrementCount(passes.LabelOfSearchEstatesWithQuery)

	return &estates, nil
}

func (c *Client) SearchEstatesNazotte(ctx context.Context, polygon *Coordinates) (*EstatesResponse, error) {
	b, err := json.Marshal(polygon)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: リクエストに失敗しました"))
	}

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/nazotte", bytes.NewBuffer(b))
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfSearchEstatesNazotte)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("POST /api/estate/nazotte: リクエストがタイムアウトしました"))
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
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("POST /api/estate/nazotte: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("POST /api/estate/nazotte: JSONデコードに失敗しました"))
	}

	passes.IncrementCount(passes.LabelOfSearchEstatesNazotte)

	return &estates, nil
}

func (c *Client) GetEstateDetailFromID(ctx context.Context, id string) (*asset.Estate, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/estate/"+id)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfGetEstateDetailFromID)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/estate/:id: リクエストがタイムアウトしました"))
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
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/estate/:id: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: JSONデコードに失敗しました"))
	}

	asset.IncrementEstateViewCount(estate.ID)
	passes.IncrementCount(passes.LabelOfGetEstateDetailFromID)

	return &estate, nil
}

func (c *Client) GetRecommendedEstatesFromChair(ctx context.Context, id int64) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/recommended_estate/"+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfGetRecommendedEstatesFromChair)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/recommended_estate/:id: リクエストがタイムアウトしました"))
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
		if timeoutCtx.Err() != nil {
			return nil, failure.New(fails.ErrTimeout, failure.Message("GET /api/recommended_estate/:id: リクエストがタイムアウトしました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: JSONデコードに失敗しました"))
	}

	passes.IncrementCount(passes.LabelOfGetRecommendedEstatesFromChair)

	return &estate, nil
}

func (c *Client) BuyChair(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair/buy/"+id, nil)
	if err != nil {
		return failure.Wrap(err, failure.Message("POST /api/chair/buy/:id: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfBuyChair)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if timeoutCtx.Err() != nil {
			return failure.New(fails.ErrTimeout, failure.Message("POST /api/chair/buy/:id: リクエストがタイムアウトしました"))
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
	passes.IncrementCount(passes.LabelOfBuyChair)

	return nil
}

func (c *Client) RequestEstateDocument(ctx context.Context, id string) error {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/req_doc/"+id, nil)
	if err != nil {
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutOfRequestEstateDocument)
	defer cancel()

	req = req.WithContext(timeoutCtx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if timeoutCtx.Err() != nil {
			return failure.New(fails.ErrTimeout, failure.Message("POST /api/estate/req_doc/:id: リクエストがタイムアウトしました"))
		}
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}
	defer res.Body.Close()

	err = checkStatusCode(res, 200)
	if err != nil {
		return failure.Wrap(err, failure.Message("POST /api/estate/req_doc/:id: リクエストに失敗しました"))
	}

	passes.IncrementCount(passes.LabelOfRequestEstateDocument)

	return nil
}
