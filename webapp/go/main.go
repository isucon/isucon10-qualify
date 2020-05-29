package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const SRID = 6668

var db *sqlx.DB

var estateRentRanges = []*Range{
	{
		ID:  0,
		Min: -1,
		Max: 50000,
	},
	{
		ID:  1,
		Min: 50000,
		Max: 100000,
	},
	{
		ID:  2,
		Min: 100000,
		Max: 150000,
	},
	{
		ID:  3,
		Min: 150000,
		Max: -1,
	},
}

var estateDoorHeightRanges = []*Range{
	{
		ID:  0,
		Min: -1,
		Max: 80,
	},
	{
		ID:  1,
		Min: 80,
		Max: 110,
	},
	{
		ID:  2,
		Min: 110,
		Max: 150,
	},
	{
		ID:  3,
		Min: 150,
		Max: -1,
	},
}

var estateDoorWidthRanges = []*Range{
	{
		ID:  0,
		Min: -1,
		Max: 80,
	},
	{
		ID:  1,
		Min: 80,
		Max: 110,
	},
	{
		ID:  2,
		Min: 110,
		Max: 150,
	},
	{
		ID:  3,
		Min: 150,
		Max: -1,
	},
}

//EstateSchema estate tableに格納されている物件データ
type EstateSchema struct {
	ID          int64   `db:"id"`
	Thumbnail   string  `db:"thumbnail"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
	Address     string  `db:"address"`
	Rent        int64   `db:"rent"`
	DoorHeight  int64   `db:"door_height"`
	DoorWidth   int64   `db:"door_width"`
	ViewCount   int64   `db:"view_count"`
	Features    string  `db:"features"`
}

func (es EstateSchema) ToEstate() Estate {
	return Estate{
		ID:          es.ID,
		Thumbnail:   es.Thumbnail,
		Name:        es.Name,
		Description: es.Description,
		Address:     es.Address,
		Latitude:    es.Latitude,
		Longitude:   es.Longitude,
		DoorHeight:  es.DoorHeight,
		DoorWidth:   es.DoorWidth,
		Rent:        es.Rent,
		Features:    es.Features,
	}
}

//EstateSearchResponse estate/searchへのレスポンスの形式
type EstateSearchResponse struct {
	Estates []Estate `json:"estates"`
}

//Estate 物件
type Estate struct {
	ID          int64   `json:"id"`
	Thumbnail   string  `json:"thumbnail"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DoorHeight  int64   `json:"doorHeight"`
	DoorWidth   int64   `json:"doorWidth"`
	Rent        int64   `json:"rent"`
	Features    string  `json:"features"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Coordinates struct {
	Coordinates []Coordinate `json:"coordinates"`
}

type Range struct {
	ID  int64 `json:"id"`
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type RangeResponse struct {
	Prefix string   `json:"prefix"`
	Suffix string   `json:"suffix"`
	Ranges []*Range `json:"ranges"`
}

type RangeResponseMap struct {
	DoorWidth  RangeResponse `json:"doorWidth"`
	DoorHeight RangeResponse `json:"doorHeight"`
	Rent       RangeResponse `json:"rent"`
}

type BoundingBox struct {
	TopLeftCorner     Coordinate
	BottomRightCorner Coordinate
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

//ConnectDB isuumoデータベースに接続する
func ConnectDB() (*sqlx.DB, error) {
	host := getEnv("MYSQL_HOST", "127.0.0.1")
	port := getEnv("MYSQL_PORT", "3306")
	user := getEnv("MYSQL_USER", "isucon")
	dbname := getEnv("MYSQL_DBNAME", "isuumo")
	password := getEnv("MYSQL_PASS", "isucon")
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, password, host, port, dbname)
	return sqlx.Open("mysql", dsn)
}

func main() {
	// Echo instance
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Estate Handler
	e.GET("/api/estate/:id", getEstateDetail)
	e.GET("/api/estate/search", searchEstates)
	e.POST("/api/estate/req_doc/:id", postEstateRequestDocument)
	e.POST("/api/estate/nazotte", searchEstateNazotte)
	e.GET("/api/estate/range", responseEstateRange)

	// Recommended Handler
	e.GET("/api/recommended_estate", searchRecommendEstate)

	var err error
	db, err = ConnectDB()
	if err != nil {
		e.Logger.Fatalf("DB connection faild : %v", err)
	}
	defer db.Close()

	// Start server
	serverPort := fmt.Sprintf(":%v", getEnv("SERVER_PORT", "1323"))
	e.Logger.Fatal(e.Start(serverPort))
}

func getEstateDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug("Request parameter \"id\" parse error :", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var estate EstateSchema
	err = db.Get(&estate, "SELECT * FROM estate WHERE id = ?", id)
	if err != nil {
		c.Echo().Logger.Debug("Database Execution error :", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		c.Echo().Logger.Debug("faild to create transaction : ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = tx.Exec("UPDATE estate SET view_count = ? WHERE id = ?", estate.ViewCount+1, id)
	if err != nil {
		c.Echo().Logger.Debug("view_count update failed :", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		c.Echo().Logger.Debug("transaction commit error : ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, estate.ToEstate())
}

func getRange(RangeID string, Ranges []*Range) (*Range, error) {
	specifyRange := &Range{}

	RangeIndex, err := strconv.Atoi(RangeID)
	if err != nil {
		return nil, err
	}
	if RangeIndex < len(Ranges) && RangeIndex > 0 {
		specifyRange = Ranges[RangeIndex]
	} else {
		err = fmt.Errorf("Unexpected Range ID")
	}

	return specifyRange, nil
}

func searchEstates(c echo.Context) error {
	var searchOption bool
	var doorHeight, doorWidth, estateRent *Range
	var err error

	var searchQueryArray []string
	var searchQueryParameter []interface{}

	if c.QueryParam("doorHeightRangeId") != "" {
		doorHeight, err = getRange(c.QueryParam("doorHeightRangeId"), estateDoorHeightRanges)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if doorHeight.Min != -1 {
			searchQueryArray = append(searchQueryArray, "door_height >= ? ")
			searchQueryParameter = append(searchQueryParameter, doorHeight.Min)
		}
		if doorHeight.Max != -1 {
			searchQueryArray = append(searchQueryArray, "door_height < ? ")
			searchQueryParameter = append(searchQueryParameter, doorHeight.Max)
		}

		searchOption = true
	}

	if c.QueryParam("doorWidthRangeId") != "" {
		doorWidth, err = getRange(c.QueryParam("doorWidthRangeId"), estateDoorWidthRanges)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if doorWidth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "door_width >= ? ")
			searchQueryParameter = append(searchQueryParameter, doorWidth.Min)
		}
		if doorWidth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "door_width < ? ")
			searchQueryParameter = append(searchQueryParameter, doorWidth.Max)
		}

		searchOption = true
	}

	if c.QueryParam("rentRangeId") != "" {
		estateRent, err = getRange(c.QueryParam("rentRangeId"), estateRentRanges)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		searchOption = true

		if estateRent.Min != -1 {
			searchQueryArray = append(searchQueryArray, "door_width >= ? ")
			searchQueryParameter = append(searchQueryParameter, estateRent.Min)
		}
		if estateRent.Max != -1 {
			searchQueryArray = append(searchQueryArray, "door_width < ? ")
			searchQueryParameter = append(searchQueryParameter, estateRent.Max)
		}

	}

	if c.QueryParam("features") != "" {
		for _, f := range strings.Split(c.QueryParam("features"), ",") {
			searchQueryArray = append(searchQueryArray, "features like concat('%', ?, '%')")
			searchQueryParameter = append(searchQueryParameter, f)
		}
		searchOption = true
	}

	if !searchOption {
		return c.String(http.StatusBadRequest, "search condition not found")
	}

	var estates EstateSearchResponse
	sqlstr := "select * from estate where "
	searchQuery := strings.Join(searchQueryArray, " and ")

	err = db.Select(&estates.Estates, sqlstr+searchQuery, searchQueryParameter...)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, estates)
}

func searchRecommendEstate(c echo.Context) error {
	limit := 20
	recommentEstates := make([]Estate, 0, limit)

	sqlstr := `select * from estate order by view_count desc limit ?`

	err := db.Select(&recommentEstates, sqlstr, limit)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, recommentEstates)
}

func searchEstateNazotte(c echo.Context) error {
	coordinates := Coordinates{}
	err := c.Bind(&coordinates)
	if err != nil {
		c.Echo().Logger.Debug("post search estate nazotte failed :", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = coordinates.coordinatesToPolygon()
	if err != nil {
		c.Echo().Logger.Debug("request coordinates are not WKT Polygon", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	b := coordinates.getBoundingBox()
	estatesInBoundingBox := []EstateSchema{}

	q := `SELECT * FROM estate WHERE latitude < ? AND latitude > ? AND longitude< ? AND longitude > ?`

	err = db.Select(&estatesInBoundingBox, q, b.TopLeftCorner.Latitude, b.BottomRightCorner.Latitude, b.TopLeftCorner.Longitude, b.BottomRightCorner.Longitude)
	if err == sql.ErrNoRows {
		c.Echo().Logger.Debug("select * from estate where latitude ...", err)
		return c.NoContent(http.StatusNoContent)
	} else if err != nil {
		c.Echo().Logger.Debug("database execution error : ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	estatesInPolygon := []EstateSchema{}
	for _, estate := range estatesInBoundingBox {
		validatedEstate := EstateSchema{}

		point := fmt.Sprintf("'POINT(%f %f)'", estate.Latitude, estate.Longitude)
		q := `SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s, %v))`
		q = fmt.Sprintf(q, coordinates.coordinatesToText(), point, SRID)

		err = db.Get(&validatedEstate, q, estate.ID)
		if err == sql.ErrNoRows {
			c.Echo().Logger.Debug("This estate is not in the polygon")
		} else if err != nil {
			c.Echo().Logger.Debug("db access is failed on executing validate if estate is in polygon", err)
			return c.NoContent(http.StatusInternalServerError)
		} else {
			estatesInPolygon = append(estatesInPolygon, validatedEstate)
		}
	}

	re := make([]Estate, 0, len(estatesInPolygon))
	for _, estate := range estatesInPolygon {
		re = append(re, estate.ToEstate())
	}

	return c.JSON(http.StatusOK, re)
}

func postEstateRequestDocument(c echo.Context) error {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug("post request document failed :", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func responseEstateRange(c echo.Context) error {
	ranges := RangeResponseMap{
		DoorHeight: RangeResponse{
			Prefix: "",
			Suffix: "cm",
			Ranges: estateDoorHeightRanges,
		},
		DoorWidth: RangeResponse{
			Prefix: "",
			Suffix: "cm",
			Ranges: estateDoorWidthRanges,
		},
		Rent: RangeResponse{
			Prefix: "",
			Suffix: "円",
			Ranges: estateRentRanges,
		},
	}

	return c.JSON(http.StatusOK, ranges)
}

//TODO: グラハムスキャンの実装
func (cs Coordinates) coordinatesToPolygon() error {
	// グラハムスキャンして、Polygonにして返す

	return nil
}

func (cs Coordinates) getBoundingBox() BoundingBox {
	coordinates := cs.Coordinates
	boundingBox := BoundingBox{
		TopLeftCorner: Coordinate{
			Latitude: coordinates[0].Latitude, Longitude: coordinates[0].Latitude,
		},
		BottomRightCorner: Coordinate{
			Latitude: coordinates[0].Longitude, Longitude: coordinates[0].Longitude,
		},
	}
	for _, coordinate := range coordinates {
		if boundingBox.TopLeftCorner.Latitude < coordinate.Latitude {
			boundingBox.TopLeftCorner.Latitude = coordinate.Latitude
		}
		if boundingBox.TopLeftCorner.Longitude < coordinate.Longitude {
			boundingBox.TopLeftCorner.Longitude = coordinate.Longitude
		}

		if boundingBox.BottomRightCorner.Latitude > coordinate.Latitude {
			boundingBox.BottomRightCorner.Latitude = coordinate.Latitude
		}
		if boundingBox.BottomRightCorner.Longitude > coordinate.Longitude {
			boundingBox.BottomRightCorner.Longitude = coordinate.Longitude
		}
	}
	return boundingBox
}

func (cs Coordinates) coordinatesToText() string {
	// return such as POLYGON((35 137,35 140,37 140, 37 137,35 137)),6668)
	// for _, c := range coordinates { fmt.Spritf("")	}
	PolygonArray := make([]string, 0, len(cs.Coordinates))
	for _, c := range cs.Coordinates {
		PolygonArray = append(PolygonArray, fmt.Sprintf("%f %f", c.Latitude, c.Longitude))
	}
	return fmt.Sprintf("'POLYGON((%s))', %d", strings.Join(PolygonArray, ","), SRID)
}
