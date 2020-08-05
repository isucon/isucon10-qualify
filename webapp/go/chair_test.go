package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	main "github.com/isucon/isucon10-qualify/webapp/go"
	"github.com/stretchr/testify/assert"
)

func TestGetChairDetail(t *testing.T) {
	client := new(http.Client)
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET]/api/chair/:id, id=10, to get info of chair", func(t *testing.T) {
		expectedID := 10
		path := "/api/chair/" + strconv.Itoa(expectedID)
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		var preChair main.ChairSchema
		var postChair main.ChairSchema

		q := `select * from chair where id = ?`
		db.Get(&preChair, q, expectedID)
		resp, err := client.Do(req)
		if err != nil {
			panic("")
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualChair main.Chair
		_ = json.Unmarshal(body, &actualChair)

		db.Get(&postChair, q, expectedID)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, *(preChair.ToChair()), actualChair)
		assert.EqualValues(t, postChair.ViewCount-preChair.ViewCount, 1)
	})
}

func TestBuyChair(t *testing.T) {
	client := new(http.Client)
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[POST] /api/chair/buy/:id, id=10, to post buy", func(t *testing.T) {
		expectedID := 10
		path := "/api/chair/buy/" + strconv.Itoa(expectedID)
		url := getURL()

		var preChair main.ChairSchema
		var postChair main.ChairSchema
		q := `select * from chair where id = ?`
		db.Get(&preChair, q, expectedID)

		req, _ := http.NewRequest("POST", url+path, nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		db.Get(&postChair, q, expectedID)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.EqualValues(t, 1, preChair.Stock-postChair.Stock, 1)
	})
}

func TestResponseChairRange(t *testing.T) {
	client := new(http.Client)
	t.Run("[GET] /api/chair/search/condition, to get chair range", func(t *testing.T) {
		path := "/api/chair/search/condition"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSearchChairs(t *testing.T) {
	client := new(http.Client)
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET] /api/chair/search, to get chair list restricted by props", func(t *testing.T) {
		path := "/api/chair/search"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)
		params := req.URL.Query()
		params.Add("priceRangeId", "1")
		priceMin := 3000
		priceMax := 6000
		params.Add("heightRangeId", "1")
		heightMin := 80
		heightMax := 110
		params.Add("widthRangeId", "1")
		widthMin := 80
		widthMax := 110
		params.Add("depthRangeId", "1")
		depthMin := 80
		depthMax := 110
		params.Add("color", "オレンジ")
		color := "オレンジ"
		params.Add("page", "0")
		params.Add("perPage", "20")
		perPage := 20
		startPos := perPage * (0)
		req.URL.RawQuery = params.Encode()
		fmt.Printf("url: %s", req.URL.String())
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		var chairs []main.ChairSchema

		q := `select * from chair where stock >= 1 and width< ? and width >= ? and height < ? and height >= ? and depth < ? and depth >= ? and price < ? and price >= ? and color = ? order by view_count desc limit ? offset ?`
		db.Select(&chairs, q, widthMax, widthMin, heightMax, heightMin, depthMax, depthMin, priceMax, priceMin, color, perPage, startPos)
		var actualChairs main.ChairSearchResponse
		_ = json.Unmarshal(body, &actualChairs)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rc []main.Chair
		for _, chair := range chairs {
			rc = append(rc, *(chair.ToChair()))
		}
		var ac []main.Chair
		for _, chairPointer := range actualChairs.Chairs {
			ac = append(ac, *chairPointer)
		}
		assert.Equal(t, rc, ac)
	})
}
