package asset

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

var (
	chairMap  map[int64]*Chair
	estateMap map[int64]*Estate
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(ctx context.Context, dataDir, fixtureDir string) {
	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
		if err != nil {
			return err
		}
		defer f.Close()

		chairMap = map[int64]*Chair{}
		decoder := json.NewDecoder(f)
		for {
			if err := childCtx.Err(); err != nil {
				return err
			}

			var chair Chair
			if err := decoder.Decode(&chair); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			StoreChair(chair)
		}
		return nil
	})

	eg.Go(func() error {
		f, err := os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
		if err != nil {
			return err
		}
		defer f.Close()

		estateMap = map[int64]*Estate{}
		decoder := json.NewDecoder(f)
		for {
			if err := childCtx.Err(); err != nil {
				return err
			}

			var estate Estate
			if err := decoder.Decode(&estate); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			StoreEstate(estate)
		}

		return nil
	})

	eg.Go(func() error {
		err := loadChairSearchCondition(fixtureDir)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		err := loadEstateSearchCondition(fixtureDir)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		err = failure.Translate(err, fails.ErrBenchmarker, failure.Message("assetの初期化に失敗しました"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfInitialize)
	}
}

func ExistsChairInMap(id int64) bool {
	_, ok := chairMap[id]
	return ok
}

func GetChairFromID(id int64) (*Chair, error) {
	var c *Chair
	if ExistsChairInMap(id) {
		c, _ = chairMap[id]
		return c, nil
	}

	return nil, errors.New("requested chair not found")
}

func StoreChair(chair Chair) {
	chairMap[chair.ID] = &chair
}

func DecrementChairStock(id int64) {
	if ExistsChairInMap(id) {
		chairMap[id].DecrementStock()
	}
}

func ExistsEstateInMap(id int64) bool {
	_, ok := estateMap[id]
	return ok
}

func GetEstateFromID(id int64) (*Estate, error) {
	var e *Estate
	if ExistsEstateInMap(id) {
		e, _ = estateMap[id]
		return e, nil
	}
	return nil, errors.New("requested estate not found")
}

func StoreEstate(estate Estate) {
	estateMap[estate.ID] = &estate
}
