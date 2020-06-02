package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	main "github.com/isucon/isucon10-qualify/webapp/go"
	"github.com/stretchr/testify/assert"
)

//ToDo: DBとmainの自動起動の仕組み
//Seedも外部設定できるとよき

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

var port = getEnv("API_PORT", "1323")
var host = getEnv("API_HOST", "localhost")
var url = fmt.Sprintf("http://%s:%s", host, port)

func TestGetChairDetail(t *testing.T) {
	client := new(http.Client)

	t.Run("[GET]/api/chair/:id, id=10, to get info of chair", func(t *testing.T) {
		expectedID := 10
		path := "/api/chair/" + strconv.Itoa(expectedID)
		req, _ := http.NewRequest("GET", url+path, nil)

		resp, err := client.Do(req)
		fmt.Printf("%v", resp)
		if err != nil {
			panic("")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualChair main.Chair
		_ = json.Unmarshal(body, &actualChair)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, int64(expectedID), actualChair.ID)
		//assert.Empty(t, actualChair.ViewCount)
		//assert.Empty(t, actualChair.Stock)
		//ToDo: DB内のViewCountの監視
	})
}

func TestBuyChair(t *testing.T) {
	client := new(http.Client)
	t.Run("[POST] /api/chair/buy/:id, id=10, to post buy", func(t *testing.T) {
		expectedID := 10
		path := "/api/chair/buy/" + strconv.Itoa(expectedID)
		req, _ := http.NewRequest("POST", url+path, nil)

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//ToDo: DB内のStockの監視
	})
}

func TestResponseChairRange(t *testing.T) {
	client := new(http.Client)
	t.Run("[GET] /api/chair/range, to get chair range", func(t *testing.T) {
		path := "/api/chair/range"
		req, _ := http.NewRequest("GET", url+path, nil)

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//ToDo: 構成要素の確認ぐらい?
	})
}

func TestSearchChairs(t *testing.T) {
	client := new(http.Client)
	t.Run("[GET] /api/chair/search, to get chair list restricted by props", func(t *testing.T) {
		path := "/api/chair/search"
		req, _ := http.NewRequest("GET", url+path, nil)
		params := req.URL.Query()
		params.Add("priceRangeId", "1")
		params.Add("hegihtRangeId", "1")
		params.Add("depthRangeId", "1")
		params.Add("color", "オレンジ")
		params.Add("page", "0")
		params.Add("perPage", "20")
		//params.Add("features", "駅から徒歩5分")
		req.URL.RawQuery = params.Encode()
		fmt.Printf("url: %s", req.URL.String())
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		var actualChairs main.ChairSearchResponce
		_ = json.Unmarshal(body, &actualChairs)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//ToDo: countでなんとかする?
		fmt.Printf("number: %d", len(actualChairs.Chairs))
	})
}
