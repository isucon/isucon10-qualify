package scenario

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

type VerifySnapShot struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Method string             `json:"method"`
	URI    string             `json:"uri"`
	ID     string             `json:"id"`
	Query  Query              `json:"query"`
	Body   client.Coordinates `json:"coordinates"`
}

type Query struct {
	RentRangeID       string `json:"rentRangeId"`
	PriceRangeID      string `json:"priceRangeId"`
	DoorHeightRangeID string `json:"doorHeightRangeId"`
	DoorWidthRangeID  string `json:"doorWidthRangeId"`
	HeightRangeID     string `json:"heightRangeId"`
	WidthRangeID      string `json:"widthRangeId"`
	DepthRangeID      string `json:"depthRangeId"`
	Features          string `json:"features"`
	Kind              string `json:"kind"`
	Page              string `json:"page"`
	PerPage           string `json:"perPage"`
}

type Response struct {
	Body Body `json:"body"`
}

type Body struct {
	Count   int64          `json:"count"`
	Estates []asset.Estate `json:"estates"`
	Chairs  []asset.Chair  `json:"chairs"`
}

// Verify Initialize後のアプリケーションサーバーに対して、副作用のない検証を実行する
// 早い段階でベンチマークをFailさせて早期リターンさせるのが目的
// ex) recommended API や Search API を叩いて初期状態を確認する
func Verify(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := verifyChairStock(ctx)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := verifyChairViewCount(ctx)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := verifyEstateViewCount(ctx)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		dirName, _ := filepath.Abs("../initial-data/generate_verification")
		err := verifyWithSnapshot(ctx, dirName)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfVerify)
		}
		wg.Done()
	}()

	wg.Wait()

	return
}

func verifyChairStock(ctx context.Context) error {
	c := client.PickClient()
	err := c.BuyChair(ctx, "1")
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	chair, err := c.GetChairDetailFromID(ctx, "1")
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if chair != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.New(fails.ErrApplication, failure.Message("イスの在庫数が不正です"))
	}

	return nil
}

func verifyChairViewCount(ctx context.Context) error {
	c := client.PickClient()

	for i := 0; i < 2; i++ {
		_, err := c.GetChairDetailFromID(ctx, "2")
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			return failure.Translate(err, fails.ErrApplication)
		}
	}

	q := url.Values{}
	q.Add("features", "フットレスト")
	q.Add("page", "0")
	q.Add("perPage", "2")

	chairs, err := c.SearchChairsWithQuery(ctx, q)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if chairs.Chairs[0].ID != 2 || chairs.Chairs[1].ID != 3 {
		return failure.New(fails.ErrApplication, failure.Message("イスの閲覧数が不正です"))
	}

	return nil
}

func verifyEstateViewCount(ctx context.Context) error {
	c := client.PickClient()

	for i := 0; i < 2; i++ {
		_, err := c.GetEstateDetailFromID(ctx, "1")
		if err != nil {
			return failure.Translate(err, fails.ErrApplication)
		}
	}

	q := url.Values{}
	q.Add("features", "デザイナーズ物件")
	q.Add("page", "0")
	q.Add("perPage", "2")

	estates, err := c.SearchEstatesWithQuery(ctx, q)

	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication)
	}

	if estates.Estates[0].ID != 1 || estates.Estates[1].ID != 2 {
		return failure.New(fails.ErrApplication, failure.Message("物件の閲覧数が不正です"))
	}

	return nil
}

func verifyWithSnapshot(ctx context.Context, snapshotsDir string) error {
	snapshotFiles, err := ioutil.ReadDir(snapshotsDir)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("verify snapshot error snapshot file not found "))
	}

	for _, sf := range snapshotFiles {
		c := client.PickClient()
		raw, err := ioutil.ReadFile(path.Join(snapshotsDir, sf.Name()))
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("verify snapshot error ReadFile error "))
		}

		var verifyData VerifySnapShot
		err = json.Unmarshal(raw, &verifyData)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("verify snapshot error Unmarshal "))
		}

		params := url.Values{}
		elems := reflect.ValueOf(&(verifyData.Request.Query)).Elem()
		for i := 0; i < elems.NumField(); i++ {
			valueField := elems.Field(i)
			typeField := elems.Type().Field(i)
			tag := typeField.Tag
			if fmt.Sprintf("%v", valueField.Interface()) == "" {
				continue
			} else {
				params.Set(tag.Get("json"), fmt.Sprintf("%v", valueField.Interface()))
			}
		}

		prefixMsg := fmt.Sprintf("%v %v: ", verifyData.Request.Method, verifyData.Request.URI)

		switch verifyData.Request.URI {
		case "/api/estate/search":
			er, err := c.SearchEstatesWithQuery(ctx, params)
			if err != nil {
				return failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
			}

			if len(er.Estates) != int(verifyData.Response.Body.Count) {
				return failure.Translate(err, fails.ErrApplication, failure.Message(prefixMsg+" 不正なレスポンスです"))
			}

			if !reflect.DeepEqual(er.Estates, verifyData.Response.Body.Estates) {
				return failure.New(fails.ErrApplication, failure.Message("物件の検索結果が不正です"))
			}
		case "/api/estate/nazotte":
			er, err := c.SearchEstatesNazotte(ctx, &verifyData.Request.Body)
			if err != nil {
				return failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
			}

			if len(er.Estates) != len(verifyData.Response.Body.Estates) {
				return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" 物件の検索結果が不正です"))
			}

			if !reflect.DeepEqual(verifyData.Response.Body.Estates, er.Estates) {
				return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" 物件の検索結果が不正です"))
			}
		case "/api/chair/search":
			cr, err := c.SearchChairsWithQuery(ctx, params)
			if err != nil {
				failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
			}

			if len(cr.Chairs) != len(verifyData.Response.Body.Chairs) {
				return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" イスの検索結果が不正です"))
			}

			if !reflect.DeepEqual(verifyData.Response.Body.Chairs, cr.Chairs) {
				return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" イスの検索結果が不正です"))
			}
		case "/api/recommended_chair":
			cr, err := c.GetRecommendedChair(ctx)
			if err != nil {
				failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
			}
			if !reflect.DeepEqual(verifyData.Response.Body.Chairs, cr.Chairs) {
				return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" レスポンスが不正です"))
			}
		case "/api/recommended_estate":
			if verifyData.Request.ID != "" {
				id, err := strconv.ParseInt(verifyData.Request.ID, 10, 64)
				if err != nil {
					failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
				}
				er, err := c.GetRecommendedEstatesFromChair(ctx, id)
				if err != nil {
					failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
				}

				if !reflect.DeepEqual(verifyData.Response.Body.Estates, er.Estates) {
					return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" レスポンスが不正です"))
				}

			} else {
				er, err := c.GetRecommendedEstate(ctx)
				if err != nil {
					failure.Translate(err, fails.ErrBenchmarker, failure.Message(prefixMsg+" リクエストに失敗しました"))
				}

				if !reflect.DeepEqual(verifyData.Response.Body.Estates, er.Estates) {
					return failure.New(fails.ErrApplication, failure.Message(prefixMsg+" レスポンスが不正です"))
				}
			}
		default:
			return failure.New(fails.ErrBenchmarker, failure.Message("snapshot invalid check API endpoint"))
		}

	}

	return nil
}
