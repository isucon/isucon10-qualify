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

type point struct {
	X float64
	Y float64
}

func convexHull(p []point) []point {
	sort.Slice(p, func(i, j int) bool {
		if p[i].X == p[j].X {
			return p[i].Y < p[i].Y
		}
		return p[i].X < p[j].X
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
	return ((b.X - a.X) * (c.Y - a.Y)) > ((b.Y - a.Y) * (c.X - a.X))
}

func ToCoordinates(po []point) *client.Coordinates {
	r := make([]*client.Coordinate, 0, len(po)+1)

	for _, p := range po {
		r = append(r, &client.Coordinate{Latitude: p.X, Longitude: p.Y})
	}

	// 始点と終点を一致させる
	r = append(r, r[0])

	return &client.Coordinates{Coordinates: r}
}

const errorDistance = 1E-6

// 点Pの周りの4点を返す
func getPointNeighbors(p point) []point {
	return []point{
		{X: p.X - errorDistance, Y: p.Y + errorDistance},
		{X: p.X + errorDistance, Y: p.Y + errorDistance},
		{X: p.X - errorDistance, Y: p.Y - errorDistance},
		{X: p.X + errorDistance, Y: p.Y - errorDistance},
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

func estateNazotteSearchScenario(ctx context.Context) error {
	passCtx, pass := context.WithCancel(ctx)
	failCtx, fail := context.WithCancel(ctx)

	var c *client.Client = client.NewClient("isucon-user")
	defer c.CloseIdleConnections()

	go func() {
		// Nazotte Search
		// create nazotte data randomly
		polygon := &client.Coordinates{}
		// corners 3 <= N <= 8
		polygonCorners := rand.Intn(6) + 3

		estateNeighborsPoint := make([]point, 0, 4*polygonCorners)
		choosedEstateIDs := make([]int64, polygonCorners)

		for i := 0; i < polygonCorners; i++ {
			target := rand.Int63n(10000)
			e := asset.GetEstateFromID(target)
			if !contains(choosedEstateIDs, e.ID) {
				p := point{X: e.Latitude, Y: e.Longitude}
				estateNeighborsPoint = append(estateNeighborsPoint, getPointNeighbors(p)...)
				choosedEstateIDs[i] = e.ID
			} else {
				i--
			}
		}

		convexHulled := convexHull(estateNeighborsPoint)
		polygon = ToCoordinates(convexHulled)

		er, err := c.SearchEstatesNazotte(ctx, polygon)
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		if len(er.Estates) < polygonCorners {
			err = failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: 検索結果が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		for _, estate := range er.Estates {
			if contains(choosedEstateIDs, estate.ID) {
				polygonCorners--
				if polygonCorners == 0 {
					break
				}
			}
		}

		if polygonCorners != 0 {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/nazotte: 検索結果が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		randomPosition := rand.Intn(len(er.Estates))
		targetID := er.Estates[randomPosition].ID
		e, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		if !e.Equal(asset.GetEstateFromID(e.ID)) {
			err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: 物件情報が不正です"))
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))
		if err != nil {
			fails.ErrorsForCheck.Add(err, fails.ErrorOfEstateNazotteSearchScenario)
			fail()
			return
		}

		pass()
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-failCtx.Done():
		return failure.New(fails.ErrApplication)
	case <-passCtx.Done():
		return nil
	}
}
