package scenario

import (
	"context"
	"math/rand"
	"sort"
	"strconv"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

const MAX_NAZOTTE_RESPONSE_LENGTH = 200

type point struct {
	Latitude  float64
	Longitude float64
}

func convexHull(p []point) []point {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Latitude == p[j].Latitude {
			return p[i].Longitude < p[i].Longitude
		}
		return p[i].Latitude < p[j].Latitude
	})

	var h []point

	// Lower hull
	for _, pt := range p {
		for len(h) >= 2 && !ccw(h[len(h)-2], h[len(h)-1], pt) {
			h = h[:len(h)-1]
		}
		h = append(h, pt)
	}

	// Upper hull
	for i, t := len(p)-2, len(h)+1; i >= 0; i-- {
		pt := p[i]
		for len(h) >= t && !ccw(h[len(h)-2], h[len(h)-1], pt) {
			h = h[:len(h)-1]
		}
		h = append(h, pt)
	}

	return h[:len(h)-1]
}

// ccw returns true if the three points make a counter-clockwise turn
func ccw(a, b, c point) bool {
	return ((b.Latitude - a.Latitude) * (c.Longitude - a.Longitude)) > ((b.Longitude - a.Longitude) * (c.Latitude - a.Latitude))
}

func ToCoordinates(po []point) *client.Coordinates {
	r := make([]*client.Coordinate, 0, len(po)+1)

	for _, p := range po {
		r = append(r, &client.Coordinate{Latitude: p.Latitude, Longitude: p.Longitude})
	}

	// 始点と終点を一致させる
	r = append(r, r[0])

	return &client.Coordinates{Coordinates: r}
}

const errorDistance = 1E-6

// 点Pの周りの4点を返す
func getPointNeighbors(p point) []point {
	return []point{
		{Latitude: p.Latitude - errorDistance, Longitude: p.Longitude + errorDistance},
		{Latitude: p.Latitude + errorDistance, Longitude: p.Longitude + errorDistance},
		{Latitude: p.Latitude - errorDistance, Longitude: p.Longitude - errorDistance},
		{Latitude: p.Latitude + errorDistance, Longitude: p.Longitude - errorDistance},
	}
}

func contains(s []int64, e int64) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func getBoundingBox(points []point) []point {
	boundingBox := []point{
		{
			// TopLeftCorner
			Latitude: points[0].Latitude, Longitude: points[0].Longitude,
		},
		{
			// BottomRightCorner
			Latitude: points[0].Latitude, Longitude: points[0].Longitude,
		},
	}

	po := points[1:]

	for _, p := range po {
		if boundingBox[0].Latitude > p.Latitude {
			boundingBox[0].Latitude = p.Latitude
		}
		if boundingBox[0].Longitude > p.Longitude {
			boundingBox[0].Longitude = p.Longitude
		}

		if boundingBox[1].Latitude < p.Latitude {
			boundingBox[1].Latitude = p.Latitude
		}
		if boundingBox[1].Longitude < p.Longitude {
			boundingBox[1].Longitude = p.Longitude
		}
	}
	return boundingBox
}

func estateNazotteSearchScenario(ctx context.Context) error {
	var c *client.Client = client.PickClient()

	// Nazotte Search
	// create nazotte data randomly
	polygon := &client.Coordinates{}
	// corners 3 <= N <= 8
	polygonCorners := rand.Intn(6) + 3

	estateNeighborsPoint := make([]point, 0, 4*polygonCorners)
	choosedEstateIDs := make([]int64, polygonCorners)

	for i := 0; i < polygonCorners; i++ {
		target := rand.Int63n(10000) + 1
		e, _ := asset.GetEstateFromID(target)
		if !contains(choosedEstateIDs, e.ID) {
			p := point{Latitude: e.Latitude, Longitude: e.Longitude}
			estateNeighborsPoint = append(estateNeighborsPoint, getPointNeighbors(p)...)
			choosedEstateIDs[i] = e.ID
		} else {
			i--
		}
	}

	convexHulled := convexHull(estateNeighborsPoint)
	polygon = ToCoordinates(convexHulled)
	boundingBox := getBoundingBox(convexHulled)

	er, err := c.SearchEstatesNazotte(ctx, polygon)
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	if len(er.Estates) > MAX_NAZOTTE_RESPONSE_LENGTH {
		err = failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	ok := true
	for _, estate := range er.Estates {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil || !e.Equal(&estate) {
			ok = false
			break
		}

		if !(boundingBox[0].Latitude <= e.Latitude && boundingBox[1].Latitude >= e.Latitude) {
			ok = false
			break
		}
		if !(boundingBox[0].Longitude <= e.Longitude && boundingBox[1].Longitude >= e.Longitude) {
			ok = false
			break
		}

		if !e.Equal(&estate) {
			ok = false
			break
		}

	}

	if !ok {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/nazotte: 検索結果が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	randomPosition := rand.Intn(len(er.Estates))
	targetID := er.Estates[randomPosition].ID
	e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	estate, err := asset.GetEstateFromID(e.ID)
	if err != nil || !e.Equal(estate) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
		return failure.New(fails.ErrApplication)
	}

	return nil
}
