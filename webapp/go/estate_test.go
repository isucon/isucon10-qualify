package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	t.Run("[GET] /api/estate/search/condition, to get estate range", func(t *testing.T) {
		path := "/api/estate/search/condition"
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

func TestSearchEstatenazotte(t *testing.T) {
	if testing.Short() {
		t.Skip("\"TestSearchEstatenazotte\" skipping")
	}
	db, err := MySQLConnectionData.ConnectDB()
	if err != nil {
		t.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	t.Run("[POST] /api/estate/nazotte, to get estate from nazotte API", func(t *testing.T) {
		coordinatesData := &main.Coordinates{
			Coordinates: []main.Coordinate{
				{35.78746239087127, 139.64962005615237},
				{35.78746239087127, 139.65511322021487},
				{35.78857638547713, 139.71485137939456},
				{35.78885488168885, 139.71519470214847},
				{35.79832317220566, 139.71553802490237},
				{35.839803254941856, 139.71622467041018},
				{35.84175145058553, 139.71622467041018},
				{35.84342129448237, 139.7007751464844},
				{35.84453450421662, 139.6901321411133},
				{35.848430615242236, 139.64447021484378},
				{35.83200999390394, 139.64447021484378},
				{35.826164545769274, 139.64481353759768},
				{35.8228240963784, 139.6451568603516},
				{35.78746239087127, 139.64962005615237},
			},
		}

		b, err := json.Marshal(coordinatesData)
		if err != nil {
			t.Fatalf("create request body error: %v", err)
		}

		path := "/api/estate/nazotte"
		url := getURL()
		resp, err := http.Post(url+path, "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("nazotte search POST request failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("nazotte search invalid response status code: want 200 but got %v", resp.StatusCode)
		}

		respbody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("nazotte search response body read failed: %v", err)
		}

		var estateResponse main.EstateSearchResponse
		err = json.Unmarshal(respbody, &estateResponse)
		if err != nil {
			t.Fatalf("response Unmarshal failed: %v", err)
		}

		var areas []string
		for _, coordinate := range coordinatesData.Coordinates {
			areas = append(areas, fmt.Sprintf("%v %v", coordinate.Latitude, coordinate.Longitude))
		}

		// Latitudeを昇順でソート
		sort.Slice(coordinatesData.Coordinates, func(i, j int) bool {
			return coordinatesData.Coordinates[i].Latitude < coordinatesData.Coordinates[j].Latitude
		})
		minLatitude := coordinatesData.Coordinates[0].Latitude
		maxLatitude := coordinatesData.Coordinates[len(coordinatesData.Coordinates)-1].Latitude

		// Longitudeを昇順でソート
		sort.Slice(coordinatesData.Coordinates, func(i, j int) bool {
			return coordinatesData.Coordinates[i].Longitude < coordinatesData.Coordinates[j].Longitude
		})
		minLongitude := coordinatesData.Coordinates[0].Longitude
		maxLongitude := coordinatesData.Coordinates[len(coordinatesData.Coordinates)-1].Longitude

		q := `SELECT * FROM estate WHERE latitude < ? AND latitude > ? AND longitude < ? AND longitude > ? AND ST_Contains(ST_PolygonFromText(?), POINT(latitude, longitude))`
		var estatesFromDB []main.EstateSchema
		err = db.Select(&estatesFromDB, q, maxLatitude, minLatitude, maxLongitude, minLongitude, fmt.Sprintf("POLYGON((%s))", strings.Join(areas, ",")))
		if err != nil {
			t.Fatalf("select estate error: %v", err)
		}

		if len(estateResponse.Estates) != len(estatesFromDB) {
			t.Fatalf("Not correct response length: want estates %v, but got %v", len(estatesFromDB), len(estateResponse.Estates))
		}

	})

}
