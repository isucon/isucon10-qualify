package asset

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	chairMap  map[int64]*Chair
	chairIDs  []int64
	estateMap map[int64]*Estate
	estateIDs []int64
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(dataDir string) {
	f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	chairMap = map[int64]*Chair{}
	chairIDs = make([]int64, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var chair Chair
		err := json.Unmarshal([]byte(scanner.Text()), &chair)
		if err != nil {
			log.Fatal(err)
		}
		chairMap[chair.ID] = &chair
		chairIDs = append(chairIDs, chair.ID)
	}
	f.Close()

	f, err = os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	estateMap = map[int64]*Estate{}
	estateIDs = make([]int64, 0)

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		var estate Estate
		err := json.Unmarshal([]byte(scanner.Text()), &estate)
		if err != nil {
			log.Fatal(err)
		}
		estateMap[estate.ID] = &estate
		estateIDs = append(estateIDs, estate.ID)
	}
	f.Close()
}

func ExistsChairInMap(id int64) bool {
	_, ok := chairMap[id]
	return ok
}

func GetChairIDs() []int64 {
	return chairIDs
}

func GetChairFromID(id int64) (*Chair, error) {
	var c *Chair
	if ExistsChairInMap(id) {
		c, _ = chairMap[id]
		return c, nil
	}

	return nil, errors.New("requested chair not found")
}

func IncrementChairViewCount(id int64) {
	if ExistsChairInMap(id) {
		chairMap[id].IncrementViewCount()
	}
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

func GetEstateIDs() []int64 {
	return estateIDs
}

func GetEstateFromID(id int64) (*Estate, error) {
	var e *Estate
	if ExistsEstateInMap(id) {
		e, _ = estateMap[id]
		return e, nil
	}
	return nil, errors.New("requested estate not found")
}

func IncrementEstateViewCount(id int64) {
	if ExistsEstateInMap(id) {
		estateMap[id].IncrementViewCount()
	}
}
