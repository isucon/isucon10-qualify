package main

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

var estateFeatureList = []string{
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
