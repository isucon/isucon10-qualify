package main

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

var estateFeatureList = []string{
	"2階以上",
	"駐車場あり",
	"室内洗浄機置き場",
	"エアコン付き",
	"オートロック",
	"洗面所独立",
	"ロフトあり",
	"ガスコンロ対応",
	"インターネット無料",
	"ユニバーサルデザイン",
	"DIY可",
	"即入居可",
	"楽器相談可",
	"保証人不要",
	"角部屋",
	"床下収納",
	"バストイレ別",
	"駅から徒歩5分",
	"ペット飼育可能",
}

func createRandomEstateSearchQuery() url.Values {
	q := url.Values{}
	q.Set("rentRangeId", strconv.Itoa(rand.Intn(4)))
	if (rand.Intn(100) % 20) == 0 {
		q.Set("doorHeightRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 20) == 0 {
		q.Set("doorWidthRangeId", strconv.Itoa(rand.Intn(4)))
	}
	if (rand.Intn(100) % 20) == 0 {
		features := make([]string, len(estateFeatureList))
		copy(features, estateFeatureList)
		rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
		featureLength := rand.Intn(3) + 1
		q.Set("features", strings.Join(features[:featureLength], ","))
	}
	q.Set("perPage", strconv.Itoa(rand.Intn(30)+20))
	q.Set("page", "0")

	return q
}
