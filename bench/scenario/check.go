package scenario

import (
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func isEstateEqualToAsset(e *asset.Estate) bool {
	estate, err := asset.GetEstateFromID(e.ID)
	if err != nil {
		return false
	}
	return estate.Equal(e)
}

func isEstatesOrderedByPopularity(e []asset.Estate) bool {
	var popularity int64 = -1
	for i, estate := range e {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil {
			return false
		}
		p := e.GetPopularity()
		if i > 0 && popularity < p {
			return false
		}
		popularity = p
	}
	return true
}

func isChairEqualToAsset(c *asset.Chair) bool {
	chair, err := asset.GetChairFromID(c.ID)
	if err != nil {
		return false
	}
	return chair.Equal(c)
}

func isChairSoldOutAt(c *asset.Chair, t time.Time) bool {
	if c.GetStock() > 0 {
		return false
	}

	soldOutTime := c.GetSoldOutTime()
	if soldOutTime == nil {
		return false
	}
	return t.After(*soldOutTime)
}

func isChairsOrderedByPopularity(c []asset.Chair, t time.Time) bool {
	var popularity int64 = -1
	for i, chair := range c {
		_chair, err := asset.GetChairFromID(chair.ID)
		if err != nil {
			return false
		}

		if isChairSoldOutAt(_chair, t) {
			return false
		}

		p := _chair.GetPopularity()

		if i > 0 && popularity < p {
			return false
		}
		popularity = p
	}
	return true
}

func isEstatesInBoundingBox(estates []asset.Estate, boundingBox [2]point) bool {
	for _, estate := range estates {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil || !e.Equal(&estate) {
			return false
		}

		if !(boundingBox[0].Latitude <= e.Latitude && boundingBox[1].Latitude >= e.Latitude) {
			return false
		}

		if !(boundingBox[0].Longitude <= e.Longitude && boundingBox[1].Longitude >= e.Longitude) {
			return false
		}
	}

	return true
}
