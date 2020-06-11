package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	main "github.com/isucon/isucon10-qualify/webapp/go"
	"github.com/stretchr/testify/assert"
)

func TestGetEstateDetail(t *testing.T) {
	client := new(http.Client)
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET]/api/estate/:id, id=10, to get info of estate", func(t *testing.T) {
		expectedID := 10
		path := "/api/estate/" + strconv.Itoa(expectedID)
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		var preEstate main.EstateSchema
		var postEstate main.EstateSchema

		q := `select * from estate where id = ?`
		db.Get(&preEstate, q, expectedID)
		resp, err := client.Do(req)
		if err != nil {
			panic("")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualEstate main.Estate
		_ = json.Unmarshal(body, &actualEstate)

		db.Get(&postEstate, q, expectedID)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, *preEstate.ToEstate(), actualEstate)
		assert.EqualValues(t, postEstate.ViewCount-preEstate.ViewCount, 1)
	})
}

func TestBuyEstate(t *testing.T) {
	client := new(http.Client)
	t.Run("[POST] /api/estate/req_doc/:id, id=10, to post buy", func(t *testing.T) {
		expectedID := 10
		path := "/api/estate/req_doc/" + strconv.Itoa(expectedID)
		url := getURL()
		req, _ := http.NewRequest("POST", url+path, nil)

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestResponseEstateRange(t *testing.T) {
	client := new(http.Client)
	t.Run("[GET] /api/estate/range, to get estate range", func(t *testing.T) {
		path := "/api/estate/range"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSearchEstates(t *testing.T) {

	client := new(http.Client)
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET] /api/estate/search, to get estate list restricted by props", func(t *testing.T) {
		path := "/api/estate/search"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)
		params := req.URL.Query()
		params.Add("doorWidthRangeId", "1")
		widthMin := 80
		widthMax := 110
		params.Add("doorHeightRangeId", "1")
		heightMin := 80
		heightMax := 110
		params.Add("rentRangeId", "1")
		rentMin := 50000
		rentMax := 100000
		params.Add("perPage", "20")
		params.Add("page", "0")
		perPage := 20
		startPos := perPage * (0)

		var estates []main.EstateSchema

		q := `select * from estate where door_width< ? and door_width >= ? and door_height < ? and door_height >= ? and rent < ? and rent >= ? order by view_count desc limit ? offset ?`
		db.Select(&estates, q, widthMax, widthMin, heightMax, heightMin, rentMax, rentMin, perPage, startPos)

		req.URL.RawQuery = params.Encode()
		fmt.Printf("url: %s", req.URL.String())
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		var actualEstates main.EstateSearchResponse
		_ = json.Unmarshal(body, &actualEstates)
		var ae []main.Estate
		for _, estatePointer := range actualEstates.Estates {
			ae = append(ae, *estatePointer)
		}

		defer resp.Body.Close()
		var re []main.Estate
		for _, estate := range estates {
			re = append(re, *estate.ToEstate())
		}
		assert.Equal(t, re, ae)
	})
}
