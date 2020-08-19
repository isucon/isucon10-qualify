package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

type InitializeResponse struct {
	Language string `json:"language"`
}

func (c *Client) Initialize(ctx context.Context) (*InitializeResponse, error) {
	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/initialize", nil)
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	// T/O付きのコンテキストが入る
	req = req.WithContext(ctx)

	res, err := c.Do(req)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("POST /initialize: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	// MEMO: /initializeの成功ステータスによって第二引数が変わる可能性がある
	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("POST /initialize: レスポンスコードが不正です"))
	}

	var initRes InitializeResponse

	err = json.NewDecoder(res.Body).Decode(&initRes)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("POST /initialize: JSONデコードに失敗しました"))
	}

	if initRes.Language == "" {
		return nil, failure.New(fails.ErrApplication, failure.Message("POST /initialize: 実装言語が設定されていません"))
	}

	return &initRes, nil
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/chair/:id: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/:id: JSONデコードに失敗しました"))
	}

	return &chair, nil
}

func loadChairsFromJSON(ctx context.Context, filePath string) ([]asset.Chair, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	chairs := []asset.Chair{}
	decoder := json.NewDecoder(f)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var chair asset.Chair
		if err := decoder.Decode(&chair); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chairs = append(chairs, chair)
	}

	return chairs, nil
}

func (c *Client) PostChairs(ctx context.Context, filePath string) error {
	chairs, err := loadChairsFromJSON(ctx, filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	var (
		b  bytes.Buffer
		fw io.Writer
	)
	w := multipart.NewWriter(&b)
	csv := ""
	for _, chair := range chairs {
		if err := ctx.Err(); err != nil {
			return err
		}
		csv += fmt.Sprintf("%s\n", chair.ToCSV())
		asset.StoreChair(chair)
	}
	r := strings.NewReader(csv)

	if fw, err = w.CreateFormFile("chairs", filePath); err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	if _, err := io.Copy(fw, r); err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	w.Close()

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair", &b)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/chair: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusCreated})
	if err != nil {
		if c.isBot {
			return failure.Translate(err, fails.ErrBot)
		}
		return failure.Wrap(err, failure.Message("POST /api/chair: リクエストに失敗しました"))
	}

	return nil
}

func (c *Client) GetChairSearchCondition(ctx context.Context) (*asset.ChairSearchCondition, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/chair/search/condition")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search/condition: リクエストに失敗しました"))
	}

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search/condition: レスポンスコードが不正です"))
	}

	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	var condition asset.ChairSearchCondition

	err = json.NewDecoder(res.Body).Decode(&condition)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/chair/search/condition: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search/condition: JSONデコードに失敗しました"))
	}

	return &condition, nil
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/chair/search: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/chair/search: JSONデコードに失敗しました"))
	}

	return &chairs, nil
}

func (c *Client) GetEstateSearchCondition(ctx context.Context) (*asset.EstateSearchCondition, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/estate/search/condition")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search/condition: リクエストに失敗しました"))
	}

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search/condition: レスポンスコードが不正です"))
	}

	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	var condition asset.EstateSearchCondition

	err = json.NewDecoder(res.Body).Decode(&condition)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/estate/search/condition: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/search/condition: JSONデコードに失敗しました"))
	}

	return &condition, nil
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/estate/search: リクエストに失敗しました"))
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("POST /api/estate/nazotte: リクエストに失敗しました"))
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/estate/:id: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/estate/:id: JSONデコードに失敗しました"))
	}

	return &estate, nil
}

func loadEstatesFromJSON(ctx context.Context, filePath string) ([]asset.Estate, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	estates := []asset.Estate{}
	decoder := json.NewDecoder(f)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var estate asset.Estate
		if err := decoder.Decode(&estate); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		estates = append(estates, estate)
	}

	return estates, nil
}

func (c *Client) PostEstates(ctx context.Context, filePath string) error {
	estates, err := loadEstatesFromJSON(ctx, filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	var (
		b  bytes.Buffer
		fw io.Writer
	)
	w := multipart.NewWriter(&b)
	csv := ""
	for _, estate := range estates {
		if err := ctx.Err(); err != nil {
			return err
		}
		csv += fmt.Sprintf("%s\n", estate.ToCSV())
		asset.StoreEstate(estate)
	}
	r := strings.NewReader(csv)

	if fw, err = w.CreateFormFile("estates", filePath); err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	if _, err := io.Copy(fw, r); err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	w.Close()

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate", &b)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Wrap(err, failure.Message("POST /api/estate: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusCreated})
	if err != nil {
		if c.isBot {
			return failure.Translate(err, fails.ErrBot)
		}
		return failure.Wrap(err, failure.Message("POST /api/estate: リクエストに失敗しました"))
	}

	return nil
}

func (c *Client) GetPopularChair(ctx context.Context) (*ChairsResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/popular_chair")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_chair: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_chair: レスポンスコードが不正です"))
	}

	var chairs ChairsResponse

	err = json.NewDecoder(res.Body).Decode(&chairs)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/popular_chair: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_chair: JSONデコードに失敗しました"))
	}

	return &chairs, nil
}

func (c *Client) GetPopularEstate(ctx context.Context) (*EstatesResponse, error) {
	req, err := c.newGetRequest(ShareTargetURLs.AppURL, "/api/popular_estate")
	if err != nil {
		return nil, failure.Translate(err, fails.ErrBenchmarker)
	}

	req = req.WithContext(ctx)
	res, err := c.Do(req)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_estate: リクエストに失敗しました"))
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	err = checkStatusCode(res, []int{http.StatusOK})
	if err != nil {
		if c.isBot {
			return nil, failure.Translate(err, fails.ErrBot)
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_estate: レスポンスコードが不正です"))
	}

	var estate EstatesResponse

	err = json.NewDecoder(res.Body).Decode(&estate)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/popular_estate: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/popular_estate: JSONデコードに失敗しました"))
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
		if nerr, ok := err.(interface{ Timeout() bool }); ok && nerr.Timeout() {
			return nil, failure.Translate(err, fails.ErrTimeout, failure.Message("GET /api/recommended_estate/:id: リクエストに失敗しました"))
		}
		return nil, failure.Wrap(err, failure.Message("GET /api/recommended_estate/:id: JSONデコードに失敗しました"))
	}

	return &estate, nil
}

type EmailRequest struct {
	Email string `json:"email"`
}

func (c *Client) BuyChair(ctx context.Context, id string) error {
	jsonStr, err := json.Marshal(EmailRequest{Email: c.GetEmail()})
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/chair/buy/"+id, bytes.NewBuffer(jsonStr))
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
	jsonStr, err := json.Marshal(EmailRequest{Email: c.GetEmail()})
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker)
	}

	req, err := c.newPostRequest(ShareTargetURLs.AppURL, "/api/estate/req_doc/"+id, bytes.NewBuffer(jsonStr))
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
