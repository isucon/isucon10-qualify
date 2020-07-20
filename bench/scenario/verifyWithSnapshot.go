package scenario

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

const (
	NumOfVerifyChairSearch                = 3
	NumOfVerifyEstateSearch               = 3
	NumOfVerifyRecommendedChair           = 1
	NumOfVerifyRecommendedEstate          = 1
	NumOfVerifyRecommendedEstateWithChair = 3
	NumOfVerifyEstateNazotte              = 3
)

type Request struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	Query    string `json:"query"`
	Body     string `json:"body"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Snapshot struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

func loadSnapshotFromFile(filePath string) (*Snapshot, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var snapshot *Snapshot
	err = json.Unmarshal(raw, &snapshot)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func verifyChairSearch(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Snapshotの読み込みに失敗しました"))
	}

	q, err := url.ParseQuery(snapshot.Request.Query)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Request QueryのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchChairsWithQuery(ctx, q)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search: イスの検索結果が不正です"))
		}

		var expected *client.ChairsResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: イスの検索結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: イスの検索結果が不正です"))
		}
	}

	return nil
}

func verifyEstateSearch(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Snapshotの読み込みに失敗しました"))
	}

	q, err := url.ParseQuery(snapshot.Request.Query)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Request QueryのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchEstatesWithQuery(ctx, q)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/search: 物件の検索結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 物件の検索結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 物件の検索結果が不正です"))
		}
	}

	return nil
}

func verifyRecommendedChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_chair: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetRecommendedChair(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/recommended_chair: イスのおすすめ結果が不正です"))
		}

		var expected *client.ChairsResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_chair: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_chair: イスのおすすめ結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_chair: イスのおすすめ結果が不正です"))
		}
	}

	return nil
}

func verifyRecommendedEstate(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetRecommendedEstate(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/recommended_estate: 物件のおすすめ結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate: 物件のおすすめ結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate: 物件のおすすめ結果が不正です"))
		}
	}

	return nil
}

func verifyRecommendedEstateWithChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate:id: Snapshotの読み込みに失敗しました"))
	}

	idx := strings.LastIndex(snapshot.Request.Resource, "/")
	if idx == -1 || idx == len(snapshot.Request.Resource)-1 {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate:id: 不正なSnapshotです"))
	}
	id, err := strconv.ParseInt(string([]rune(snapshot.Request.Resource)[idx+1:]), 10, 64)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate:id: 不正なSnapshotです"))
	}

	actual, err := c.GetRecommendedEstatesFromChair(ctx, id)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/recommended_estate:id: 物件のおすすめ結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate:id: Response BodyのUnmarshalでエラーが発生しました"))
		}
		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate:id: 物件のおすすめ結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate:id: 物件のおすすめ結果が不正です"))
		}
	}

	return nil
}

func verifyEstateNazotte(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Snapshotの読み込みに失敗しました"))
	}

	var coordinates *client.Coordinates
	err = json.Unmarshal([]byte(snapshot.Request.Body), &coordinates)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Request BodyのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchEstatesNazotte(ctx, coordinates)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 物件の検索結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !reflect.DeepEqual(expected, actual) {
			return failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 物件の検索結果が不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 物件の検索結果が不正です"))
		}
	}

	return nil
}

func verifyWithSnapshot(ctx context.Context, c *client.Client, snapshotsParentsDirPath string) {
	wg := sync.WaitGroup{}

	snapshotsDirPath := filepath.Join(snapshotsParentsDirPath, "chair_search")
	snapshots, err := ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyChairSearch; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyChairSearch(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_search")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyEstateSearch; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateSearch(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "recommended_chair")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_chair: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyRecommendedChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyRecommendedChair(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "recommended_estate")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyRecommendedEstate; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyRecommendedEstate(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "recommended_estate_with_chair")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate:id: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyRecommendedEstateWithChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyRecommendedEstateWithChair(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_nazotte")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyEstateNazotte; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateNazotte(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	wg.Wait()
}
