package main

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

var chairKindList = []string{
	"ゲーミングチェア",
	"座椅子",
	"エルゴノミクス",
	"ハンモック",
}

var chairColorList = []string{
	"黒",
	"白",
	"赤",
	"青",
	"緑",
	"黄",
	"紫",
	"ピンク",
	"オレンジ",
	"水色",
	"ネイビー",
	"ベージュ",
}

var chairFeatureList = []string{
	"折りたたみ可",
	"肘掛け",
	"キャスター",
	"リクライニング",
	"高さ調節可",
	"フットレスト",
}

func createRandomChairSearchQuery() url.Values {
	q := url.Values{}
	q.Set("priceRangeId", strconv.Itoa(rand.Intn(6)))
	if (rand.Intn(100) % 5) == 0 {
		q.Set("heightRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 5) == 0 {
		q.Set("widthRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 5) == 0 {
		q.Set("depthRangeId", strconv.Itoa(rand.Intn(4)))
	}

	if (rand.Intn(100) % 20) == 0 {
		q.Set("kind", chairKindList[rand.Intn(len(chairKindList))])
	}
	if (rand.Intn(100) % 20) == 0 {
		q.Set("color", chairColorList[rand.Intn(len(chairColorList))])
	}
	if (rand.Intn(100) % 20) == 0 {
		features := make([]string, len(chairFeatureList))
		copy(features, chairFeatureList)
		rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
		featureLength := rand.Intn(3) + 1
		q.Set("features", strings.Join(features[:featureLength], ","))
	}

	q.Set("perPage", strconv.Itoa(rand.Intn(30)+20))
	q.Set("page", "0")

	return q
}
