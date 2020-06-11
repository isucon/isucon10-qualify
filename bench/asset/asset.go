package asset

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	chairMap  map[int64]*Chair
	estateMap map[int64]*Estate
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(dataDir string) {
	f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	chairMap = map[int64]*Chair{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var chair Chair
		err := json.Unmarshal([]byte(scanner.Text()), &chair)
		if err != nil {
			log.Fatal(err)
		}
		chairMap[chair.ID] = &chair
	}
	f.Close()

	f, err = os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	estateMap = map[int64]*Estate{}

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		var estate Estate
		err := json.Unmarshal([]byte(scanner.Text()), &estate)
		if err != nil {
			log.Fatal(err)
		}
		estateMap[estate.ID] = &estate
	}
	f.Close()
}

func ExistsChairInMap(id int64) bool {
	_, ok := chairMap[id]
	return ok
}

func GetChairFromID(id int64) *Chair {
	c, _ := chairMap[id]
	return c
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

func GetEstateFromID(id int64) *Estate {
	e, _ := estateMap[id]
	return e
}

func IncrementEstateViewCount(id int64) {
	if ExistsEstateInMap(id) {
		estateMap[id].IncrementViewCount()
	}
}
