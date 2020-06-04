package asset

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	chairs  []Chair
	estates []Estate
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(dataDir string) {
	f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	chair := &Chair{}

	for scanner.Scan() {
		err := json.Unmarshal([]byte(scanner.Text()), chair)
		if err != nil {
			log.Fatal(err)
		}
		chairs = append(chairs, *chair)
	}
	f.Close()

	f, err = os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
	if err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(f)
	estate := &Estate{}

	for scanner.Scan() {
		err := json.Unmarshal([]byte(scanner.Text()), estate)
		if err != nil {
			log.Fatal(err)
		}
		estates = append(estates, *estate)
	}
	f.Close()
}
