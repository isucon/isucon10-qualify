package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	main "github.com/isucon/isucon10-qualify/webapp/go"
	"github.com/stretchr/testify/assert"
)

func TestRecommendEstate(t *testing.T) {
	client := new(http.Client)
	db, err := main.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET]/api/recommended_estate , to get recommended_estate", func(t *testing.T) {
		path := "/api/recommended_estate"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		var estates []main.EstateSchema

		q := `select * from estate order by view_count desc limit 20`
		db.Select(&estates, q)

		var expectedEstates []main.Estate
		for _, estate := range estates {
			expectedEstates = append(expectedEstates, *estate.ToEstate())
		}

		resp, err := client.Do(req)
		if err != nil {
			panic("")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualEstates main.EstateSearchResponse
		_ = json.Unmarshal(body, &actualEstates)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedEstates, actualEstates.Estates)
	})
}

func TestRecommendEstateWithChair(t *testing.T) {
	client := new(http.Client)
	db, err := main.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET]/api/recommended_estate/:id , to get recommended_estate", func(t *testing.T) {
		expectedID := 10
		path := "/api/recommended_estate/" + strconv.Itoa(expectedID)
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		var estates []main.EstateSchema
		var chair main.ChairSchema
		chairSQL := `select * from chair where id = ?`
		db.Get(&chair, chairSQL, expectedID)
		chairLength := []int{int(chair.Height), int(chair.Width), int(chair.Depth)}
		sort.Ints(chairLength)

		t.Logf("%v", chairLength)
		q := `select * from estate where (door_width >= ? and door_height >= ?) or (door_width >= ? and door_height>=?) order by view_count desc limit 20`
		db.Select(&estates, q, chairLength[0], chairLength[1], chairLength[1], chairLength[0])
		var expectedEstates []main.Estate
		for _, estate := range estates {
			expectedEstates = append(expectedEstates, *estate.ToEstate())
		}

		resp, err := client.Do(req)
		if err != nil {
			panic("")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualEstates main.EstateSearchResponse
		_ = json.Unmarshal(body, &actualEstates)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedEstates, actualEstates.Estates)
	})
}
func TestRecommendChair(t *testing.T) {
	client := new(http.Client)
	db, err := main.ConnectDB()
	if err != nil {
		fmt.Printf("DB connection failed :%v", err)
	}
	defer db.Close()

	t.Run("[GET]/api/recommended_chair, to get recommended_chair", func(t *testing.T) {
		path := "/api/recommended_chair"
		url := getURL()
		req, _ := http.NewRequest("GET", url+path, nil)

		var chairs []main.ChairSchema

		q := `select * from chair where stock >=1 order by view_count desc limit 20`
		db.Select(&chairs, q)

		var expectedChairs []main.Chair
		for _, chair := range chairs {
			expectedChairs = append(expectedChairs, *(chair.ToChair()))
		}

		resp, err := client.Do(req)
		if err != nil {
			panic("")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var actualChairs main.ChairSearchResponce
		_ = json.Unmarshal(body, &actualChairs)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedChairs, actualChairs.Chairs)
	})
}
