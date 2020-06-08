package asset

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	chairMap  sync.Map
	estateMap sync.Map
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(dataDir string) {
	f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var chair Chair
		err := json.Unmarshal([]byte(scanner.Text()), &chair)
		if err != nil {
			log.Fatal(err)
		}
		chairMap.Store(chair.ID, &chair)
	}
	f.Close()

	f, err = os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		var estate Estate
		err := json.Unmarshal([]byte(scanner.Text()), &estate)
		if err != nil {
			log.Fatal(err)
		}
		estateMap.Store(estate.ID, &estate)
	}
	f.Close()
}

func GetChairFromID(ID int64) *Chair {
	chair, ok := chairMap.Load(ID)
	if !ok {
		return nil
	}
	return chair.(*Chair)
}

func GetEstateFromID(ID int64) *Estate {
	estate, ok := estateMap.Load(ID)
	if !ok {
		return nil
	}
	return estate.(*Estate)
}
