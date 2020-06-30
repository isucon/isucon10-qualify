package scenario

import (
	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func checkSearchedEstateViewCount(e []asset.Estate) bool {
	var viewCount int64 = -1
	for i, estate := range e {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil {
			return false
		}
		vc := e.GetViewCount()
		if i > 0 && viewCount-vc < -3 {
			return false
		}
		viewCount = vc
	}
	return true
}

func checkSearchedChairViewCount(c []asset.Chair) bool {
	var viewCount int64 = -1
	for i, chair := range c {
		_chair, err := asset.GetChairFromID(chair.ID)
		if err != nil {
			return false
		}

		if _chair.GetStock() <= 0 {
			return false
		}

		vc := _chair.GetViewCount()

		if i > 0 && viewCount-vc < -3 {
			return false
		}
		viewCount = vc
	}
	return true
}
