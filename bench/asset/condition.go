package asset

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

var (
	chairSearchCondition   *ChairSearchCondition
	chairFeatureForVerify  *string
	estateSearchCondition  *EstateSearchCondition
	estateFeatureForVerify *string
)

type Range struct {
	ID  int64 `json:"id"`
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type RangeCondition struct {
	Prefix string   `json:"prefix"`
	Suffix string   `json:"suffix"`
	Ranges []*Range `json:"ranges"`
}

type ListCondition struct {
	List []string `json:"list"`
}

type EstateSearchCondition struct {
	DoorWidth  RangeCondition `json:"doorWidth"`
	DoorHeight RangeCondition `json:"doorHeight"`
	Rent       RangeCondition `json:"rent"`
	Feature    ListCondition  `json:"feature"`
}

type ChairSearchCondition struct {
	Width   RangeCondition `json:"width"`
	Height  RangeCondition `json:"height"`
	Depth   RangeCondition `json:"depth"`
	Price   RangeCondition `json:"price"`
	Color   ListCondition  `json:"color"`
	Feature ListCondition  `json:"feature"`
	Kind    ListCondition  `json:"kind"`
}

func loadChairSearchCondition(fixtureDir string) error {
	jsonText, err := ioutil.ReadFile(filepath.Join(fixtureDir, "chair_condition.json"))
	if err != nil {
		return err
	}

	var condition *ChairSearchCondition
	json.Unmarshal(jsonText, &condition)
	// condition.Featureの最後の1つはVerify用で該当件数が少ないため、Validationのシナリオ内では使用しない
	chairFeatureForVerify = &condition.Feature.List[len(condition.Feature.List)-1]
	condition.Feature.List = condition.Feature.List[:len(condition.Feature.List)-1]
	chairSearchCondition = condition
	return nil
}

func loadEstateSearchCondition(fixtureDir string) error {
	jsonText, err := ioutil.ReadFile(filepath.Join(fixtureDir, "estate_condition.json"))
	if err != nil {
		return err
	}
	var condition *EstateSearchCondition
	json.Unmarshal(jsonText, &condition)
	// condition.Featureの最後の1つはVerify用で該当件数が少ないため、Validationのシナリオ内では使用しない
	estateFeatureForVerify = &condition.Feature.List[len(condition.Feature.List)-1]
	condition.Feature.List = condition.Feature.List[:len(condition.Feature.List)-1]
	estateSearchCondition = condition
	return nil
}

func GetChairSearchCondition() (*ChairSearchCondition, error) {
	if chairSearchCondition == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("イスの検索条件が読み込まれていません"))
	}
	return chairSearchCondition, nil
}

func GetChairFeatureForVerify() (*string, error) {
	if chairFeatureForVerify == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("イスの検索条件が読み込まれていません"))
	}
	return chairFeatureForVerify, nil
}

func GetEstateSearchCondition() (*EstateSearchCondition, error) {
	if estateSearchCondition == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("物件の検索条件が読み込まれていません"))
	}
	return estateSearchCondition, nil
}

func GetEstateFeatureForVerify() (*string, error) {
	if estateSearchCondition == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("物件の検索条件が読み込まれていません"))
	}
	return estateFeatureForVerify, nil
}
