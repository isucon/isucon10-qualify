package scenario

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

const (
	NumOfVerifyChairSearchCondition   = 3
	NumOfVerifyChairSearch            = 3
	NumOfVerifyEstateSearchCondition  = 3
	NumOfVerifyEstateSearch           = 3
	NumOfVerifyPopularChair           = 1
	NumOfVerifyPopularEstate          = 1
	NumOfVerifyPopularEstateWithChair = 3
	NumOfVerifyEstateNazotte          = 3
)

var (
	ignoreChairUnexported  = cmpopts.IgnoreUnexported(asset.Chair{})
	ignoreEstateUnexported = cmpopts.IgnoreUnexported(asset.Estate{})
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

func verifyChairSearchCondition(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetChairSearchCondition(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search/condition: イスの検索条件が不正です"))
		}

		var expected *asset.ChairSearchCondition
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreChairUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search/condition: イスの検索条件が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search/condition: イスの検索条件が不正です"))
		}
	}

	return nil
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

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreChairUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: イスの検索結果が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: イスの検索結果が不正です"))
		}
	}

	return nil
}

func verifyEstateSearchCondition(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetEstateSearchCondition(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/search/condition: 物件の検索条件が不正です"))
		}

		var expected *asset.EstateSearchCondition
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !cmp.Equal(*expected, *actual) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search/condition: 物件の検索条件が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search/condition: 物件の検索条件が不正です"))
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

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreEstateUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 物件の検索結果が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: 物件の検索結果が不正です"))
		}
	}

	return nil
}

func verifyPopularChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_chair: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetPopularChair(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/popular_chair: イスのおすすめ結果が不正です"))
		}

		var expected *client.ChairsResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_chair: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreChairUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_chair: イスのおすすめ結果が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_chair: イスのおすすめ結果が不正です"))
		}
	}

	return nil
}

func verifyPopularEstate(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetPopularEstate(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/popular_estate: 物件のおすすめ結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate: Response BodyのUnmarshalでエラーが発生しました"))
		}

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreEstateUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_estate: 物件のおすすめ結果が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_estate: 物件のおすすめ結果が不正です"))
		}
	}

	return nil
}

func verifyPopularEstateWithChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate:id: Snapshotの読み込みに失敗しました"))
	}

	idx := strings.LastIndex(snapshot.Request.Resource, "/")
	if idx == -1 || idx == len(snapshot.Request.Resource)-1 {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate:id: 不正なSnapshotです"))
	}
	id, err := strconv.ParseInt(string([]rune(snapshot.Request.Resource)[idx+1:]), 10, 64)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate:id: 不正なSnapshotです"))
	}

	actual, err := c.GetPopularEstatesFromChair(ctx, id)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/popular_estate:id: 物件のおすすめ結果が不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate:id: Response BodyのUnmarshalでエラーが発生しました"))
		}
		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreEstateUnexported))
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_estate:id: 物件のおすすめ結果が不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/popular_estate:id: 物件のおすすめ結果が不正です"))
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

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported) {
			log.Printf("%s\n%s\n", filePath, cmp.Diff(*expected, *actual, ignoreEstateUnexported))
			return failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 物件の検索結果が不正です"), failure.Messagef("snapshot: %s", filePath))
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

	snapshotsDirPath := filepath.Join(snapshotsParentsDirPath, "chair_search_condition")
	snapshots, err := ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyChairSearchCondition; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyChairSearchCondition(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "chair_search")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
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

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_search_condition")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyEstateSearchCondition; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateSearchCondition(ctx, c, filePath)
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

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "popular_chair")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_chair: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyPopularChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyPopularChair(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "popular_estate")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyPopularEstate; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyPopularEstate(ctx, c, filePath)
				if err != nil {
					fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "popular_estate_with_chair")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/popular_estate:id: Snapshotディレクトリがありません"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
	} else {
		for i := 0; i < NumOfVerifyPopularEstateWithChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyPopularEstateWithChair(ctx, c, filePath)
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
