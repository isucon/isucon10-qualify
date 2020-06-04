package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const SRID = 6668

const LIMIT = 20

var db *sqlx.DB
var MySQLConnectionData *MySQLConnectionEnv

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

var ChairPriceRanges = []*Range{
	{
		ID:  0,
		Min: -1,
		Max: 3000,
	},
	{
		ID:  1,
		Min: 3000,
		Max: 6000,
	},
	{
		ID:  2,
		Min: 6000,
		Max: 9000,
	},
	{
		ID:  3,
		Min: 9000,
		Max: 12000,
	},
	{
		ID:  4,
		Min: 12000,
		Max: 15000,
	},
	{
		ID:  5,
		Min: 15000,
		Max: -1,
	},
}
var ChairHeightRanges = []*Range{
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

var ChairWidthRanges = []*Range{
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

var ChairDepthRanges = []*Range{
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

type ChairSchema struct {
	ID          int64  `db:"id"`
	Thumbnail   string `db:"thumbnail"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Price       int64  `db:"price"`
	Height      int64  `db:"height"`
	Width       int64  `db:"width"`
	Depth       int64  `db:"depth"`
	ViewCount   int64  `db:"view_count"`
	Stock       int64  `db:"stock"`
	Color       string `db:"color"`
	Features    string `db:"features"`
	Kind        string `db:"kind"`
}

type Chair struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Price       int64  `json:"price"`
	Height      int64  `json:"height"`
	Width       int64  `json:"width"`
	Depth       int64  `json:"depth"`
	Color       string `json:"color"`
	Features    string `json:"features"`
	Kind        string `json:"kind"`
}

type ChairSearchResponce struct {
	Chairs []Chair `json:"chairs"`
}

func (cs *ChairSchema) ToChair() *Chair {
	return &Chair{
		ID:          cs.ID,
		Name:        cs.Name,
		Description: cs.Description,
		Thumbnail:   cs.Thumbnail,
		Price:       cs.Price,
		Height:      cs.Height,
		Width:       cs.Width,
		Depth:       cs.Depth,
		Color:       cs.Color,
		Features:    cs.Features,
		Kind:        cs.Kind,
	}
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

func (es *EstateSchema) ToEstate() *Estate {
	return &Estate{
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

type RangeIntResponce struct {
	ID    int64 `json:"id"`
	Range Range `json:"range`
}

type RangeResponse struct {
	Prefix string   `json:"prefix"`
	Suffix string   `json:"suffix"`
	Ranges []*Range `json:"ranges"`
}

type RangeResponseEstateMap struct {
	DoorWidth  RangeResponse `json:"doorWidth"`
	DoorHeight RangeResponse `json:"doorHeight"`
	Rent       RangeResponse `json:"rent"`
}

type RangeResponseChairMap struct {
	Width  RangeResponse `json:"width"`
	Height RangeResponse `json:"height"`
	Depth  RangeResponse `json:"depth"`
	Price  RangeResponse `json:"price"`
}

type BoundingBox struct {
	TopLeftCorner     Coordinate
	BottomRightCorner Coordinate
}

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func NewMySQLConnectionEnv() *MySQLConnectionEnv {
	return &MySQLConnectionEnv{
		Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "isuumo"),
		Password: getEnv("MYSQL_PASS", "isucon"),
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

//ConnectDB isuumoデータベースに接続する
func (mc *MySQLConnectionEnv) ConnectDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
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

	// Initialize
	e.GET("/initialize", initialize)

	// Chair Handler
	e.GET("/api/chair/:id", getChairDetail)
	e.GET("/api/chair/search", searchChairs)
	e.POST("/api/chair/buy/:id", buyChair)
	e.GET("/api/chair/range", responseChairRange)

	// Estate Handler
	e.GET("/api/estate/:id", getEstateDetail)
	e.GET("/api/estate/search", searchEstates)
	e.POST("/api/estate/req_doc/:id", postEstateRequestDocument)
	e.POST("/api/estate/nazotte", searchEstateNazotte)
	e.GET("/api/estate/range", responseEstateRange)

	// Recommended Handler
	e.GET("/api/recommended_estate", searchRecommendEstate)
	e.GET("/api/recommended_estate/:id", searchRecommendEstateWithChair)
	e.GET("/api/recommendes_chair", searchRecommendChair)

	MySQLConnectionData = NewMySQLConnectionEnv()

	var err error
	db, err = MySQLConnectionData.ConnectDB()
	if err != nil {
		e.Logger.Fatalf("DB connection faild : %v", err)
	}
	defer db.Close()

	// Start server
	serverPort := fmt.Sprintf(":%v", getEnv("SERVER_PORT", "1323"))
	e.Logger.Fatal(e.Start(serverPort))
}

func initialize(c echo.Context) error {
	fpathprefix := filepath.Join("..", "mysql", "db")
	paths := []string{
		filepath.Join(fpathprefix, "0_Schema.sql"),
		filepath.Join(fpathprefix, "1_DummyEstateData.sql"),
		filepath.Join(fpathprefix, "2_DummyChairData.sql"),
	}

	for _, p := range paths {
		sqlFile, _ := filepath.Abs(p)
		cmdstr := fmt.Sprintf("mysql -h %v -u %v -p%v %v < %v",
			MySQLConnectionData.Host,
			MySQLConnectionData.User,
			MySQLConnectionData.Password,
			MySQLConnectionData.DBName,
			sqlFile,
		)
		if err := exec.Command("bash", "-c", cmdstr).Run(); err != nil {
			c.Logger().Errorf("Initialize script error : %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.NoContent(http.StatusOK)
}

func getChairDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Errorf("Request parameter \"id\" parse error : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	chair := ChairSchema{}
	sqlstr := "SELECT * FROM chair WHERE id = ?"
	err = db.Get(&chair, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("requested id's chair not found : %v", id)
			return c.JSON(http.StatusNoContent, chair.ToChair())
		}
		c.Echo().Logger.Errorf("Faild to get the chair from id : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	} else if chair.Stock <= 0 {
		return c.NoContent(http.StatusNotFound)
	}

	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		c.Echo().Logger.Errorf("faild to create transaction : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = tx.Exec("UPDATE chair SET view_count = ? WHERE id = ?", chair.ViewCount+1, id)
	if err != nil {
		c.Echo().Logger.Errorf("view_count update failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		c.Echo().Logger.Errorf("transaction commit error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, chair.ToChair())

}

func searchChairs(c echo.Context) error {
	var searchOption bool
	var chairHeight, chairWidth, chairDepth, chairPrice *Range
	var err error

	var searchQueryArray []string
	queryParams := make([]interface{}, 0)

	if c.QueryParam("priceRangeId") != "" {
		chairPrice, err = getRange(c.QueryParam("priceRangeId"), ChairPriceRanges)
		if err != nil {
			c.Echo().Logger.Infof("priceRangeID invalid, %v : %v", c.QueryParam("priceRangeId"), err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		searchOption = true

		if chairPrice.Min != -1 {
			searchQueryArray = append(searchQueryArray, "price >= ? ")
			queryParams = append(queryParams, chairPrice.Min)
		}
		if chairPrice.Max != -1 {
			searchQueryArray = append(searchQueryArray, "price < ? ")
			queryParams = append(queryParams, chairPrice.Max)
		}
	}

	if c.QueryParam("heightRangeId") != "" {
		chairHeight, err = getRange(c.QueryParam("heightRangeId"), ChairHeightRanges)
		if err != nil {
			c.Echo().Logger.Infof("heightRangeIf invalid, %v : %v", c.QueryParam("heightRangeId"), err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		if chairHeight.Min != -1 {
			searchQueryArray = append(searchQueryArray, "height >= ? ")
			queryParams = append(queryParams, chairHeight.Min)
		}
		if chairHeight.Max != -1 {
			searchQueryArray = append(searchQueryArray, "height < ? ")
			queryParams = append(queryParams, chairHeight.Max)
		}

		searchOption = true
	}

	if c.QueryParam("widthRangeId") != "" {
		chairWidth, err = getRange(c.QueryParam("widthRangeId"), ChairWidthRanges)
		if err != nil {
			c.Echo().Logger.Infof("widthRangeID invalid, %v : %v", c.QueryParam("widthRangeId"), err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		if chairWidth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "width >= ? ")
			queryParams = append(queryParams, chairWidth.Min)
		}
		if chairWidth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "width < ? ")
			queryParams = append(queryParams, chairWidth.Max)
		}

		searchOption = true
	}

	if c.QueryParam("depthRangeId") != "" {
		chairDepth, err = getRange(c.QueryParam("depthRangeId"), ChairDepthRanges)
		if err != nil {
			c.Echo().Logger.Infof("depthRangeId invalid, %v : %v", c.QueryParam("depthRangeId"), err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		if chairDepth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "depth>= ? ")
			queryParams = append(queryParams, chairDepth.Min)
		}
		if chairDepth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "depth < ? ")
			queryParams = append(queryParams, chairDepth.Max)
		}

		searchOption = true
	}

	if c.QueryParam("color") != "" {
		searchQueryArray = append(searchQueryArray, "color = ?")
		queryParams = append(queryParams, c.QueryParam("color"))
	}

	if c.QueryParam("features") != "" {
		for _, f := range strings.Split(c.QueryParam("features"), ",") {
			searchQueryArray = append(searchQueryArray, "features LIKE CONCAT('%', ?, '%')")
			queryParams = append(queryParams, f)
		}
		searchOption = true
	}

	if !searchOption {
		c.Echo().Logger.Infof("Search condition not found")
		return c.String(http.StatusBadRequest, "search condition not found")
	} else {
		searchQueryArray = append(searchQueryArray, "stock > 0")
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		c.Logger().Infof("Invalid format page parameter : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	perpage, err := strconv.Atoi(c.QueryParam("perPage"))
	if err != nil {
		c.Logger().Infof("Invalid format perPage parameter : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	var chairs ChairSearchResponce
	sqlstr := "SELECT * FROM chair WHERE "
	searchCondition := strings.Join(searchQueryArray, " AND ")

	limitOffset := " ORDER BY view_count DESC LIMIT ? OFFSET ?"
	queryParams = append(queryParams, perpage, page*perpage)

	searchedchairs := []ChairSchema{}
	err = db.Select(&searchedchairs, sqlstr+searchCondition+limitOffset, queryParams...)
	if err != nil {
		c.Logger().Errorf("searchChairs DB execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	for _, c := range searchedchairs {
		chairs.Chairs = append(chairs.Chairs, *c.ToChair())
	}

	return c.JSON(http.StatusOK, chairs)
}

func buyChair(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var chair ChairSchema
	err = db.Get(&chair, "SELECT * FROM chair WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("buyChair chair id \"%v\" not found", id)
			return c.NoContent(http.StatusNoContent)
		}
		c.Echo().Logger.Errorf("DB Execution Error: on getting a chair by id : %v", err)
		return c.NoContent(http.StatusNotFound)
	}

	tx, err := db.Begin()
	if err != nil {
		c.Echo().Logger.Errorf("faild to create transaction : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE chair SET stock = ? WHERE id = ?", chair.Stock-1, id)
	if err != nil {
		c.Echo().Logger.Errorf("view_count update failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		c.Echo().Logger.Errorf("transaction commit error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func responseChairRange(c echo.Context) error {
	ranges := RangeResponseChairMap{
		Height: RangeResponse{
			Prefix: "",
			Suffix: "cm",
			Ranges: ChairHeightRanges,
		},
		Width: RangeResponse{
			Prefix: "",
			Suffix: "cm",
			Ranges: ChairWidthRanges,
		},
		Depth: RangeResponse{
			Prefix: "",
			Suffix: "cm",
			Ranges: ChairDepthRanges,
		},
		Price: RangeResponse{
			Prefix: "",
			Suffix: "円",
			Ranges: ChairPriceRanges,
		},
	}
	return c.JSON(http.StatusOK, ranges)
}

func searchRecommendChair(c echo.Context) error {
	var recommendChairs []ChairSchema

	sqlstr := `SELECT * FROM chair WHERE stock >= 1 ORDER BY view_count DESC LIMIT ?`

	err := db.Select(&recommendChairs, sqlstr, LIMIT)
	if err != nil {
		c.Logger().Errorf("searchRecommendChair DB execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var rc ChairSearchResponce

	for _, chair := range recommendChairs {
		rc.Chairs = append(rc.Chairs, *(chair.ToChair()))
	}

	return c.JSON(http.StatusOK, rc)
}

func getEstateDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("Request parameter \"id\" parse error : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	var estate EstateSchema
	err = db.Get(&estate, "SELECT * FROM estate WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("getEstateDetail estate id %v not found", id)
			return c.JSON(http.StatusNoContent, estate.ToEstate)
		}
		c.Echo().Logger.Errorf("Database Execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		c.Echo().Logger.Errorf("faild to create transaction : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = tx.Exec("UPDATE estate SET view_count = ? WHERE id = ?", estate.ViewCount+1, id)
	if err != nil {
		c.Echo().Logger.Errorf("view_count update failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		c.Echo().Logger.Errorf("transaction commit error : %v", err)
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
			c.Echo().Logger.Infof("doorHeightRangeID invalid, %v : %v", c.QueryParam("doorHeightRangeId"), err)
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
			c.Echo().Logger.Infof("doorWidthRangeID invalid, %v : %v", c.QueryParam("doorWidthRangeId"), err)
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
			c.Echo().Logger.Infof("rentRangeID invalid, %v : %v", c.QueryParam("rentRangeId"), err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		searchOption = true

		if estateRent.Min != -1 {
			searchQueryArray = append(searchQueryArray, "rent >= ? ")
			searchQueryParameter = append(searchQueryParameter, estateRent.Min)
		}
		if estateRent.Max != -1 {
			searchQueryArray = append(searchQueryArray, "rent < ? ")
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
		c.Echo().Logger.Infof("searchEstates search condition not found")
		return c.String(http.StatusBadRequest, "search condition not found")
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		c.Logger().Infof("Invalid format page parameter : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	perpage, err := strconv.Atoi(c.QueryParam("perPage"))
	if err != nil {
		c.Logger().Infof("Invalid format perPage parameter : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	var estates EstateSearchResponse
	sqlstr := "SELECT * FROM estate WHERE "
	searchQuery := strings.Join(searchQueryArray, " AND ")

	limitOffset := " ORDER BY view_count DESC LIMIT ? OFFSET ?"
	searchQueryParameter = append(searchQueryParameter, perpage, page*perpage)

	matchestates := []EstateSchema{}
	err = db.Select(&matchestates, sqlstr+searchQuery+limitOffset, searchQueryParameter...)
	if err != nil {
		c.Logger().Errorf("searchEstates DB execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	for _, e := range matchestates {
		estates.Estates = append(estates.Estates, *e.ToEstate())
	}

	return c.JSON(http.StatusOK, estates)
}

func searchRecommendEstate(c echo.Context) error {
	limit := 20
	recommentEstates := make([]EstateSchema, 0, limit)

	sqlstr := `SELECT * FROM estate ORDER BY view_count DESC LIMIT ?`

	err := db.Select(&recommentEstates, sqlstr, limit)
	if err != nil {
		c.Logger().Errorf("searchRecommendEstate DB execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var re EstateSearchResponse

	for _, estate := range recommentEstates {
		re.Estates = append(re.Estates, *estate.ToEstate())
	}

	return c.JSON(http.StatusOK, re)
}

func searchRecommendEstateWithChair(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Logger().Infof("Invalid format searchRecommendedEstateWithChair id : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	chair := ChairSchema{}
	sqlstr := `SELECT * FROM chair WHERE id = ?`

	err = db.Get(&chair, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Infof("Reqested chair id \"%v\" nof found", id)
			return c.String(http.StatusBadRequest, "Chair Not Found Invalid Chair ID")
		}
		c.Logger().Errorf("Database execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	limit := 20
	var recommendEstates []EstateSchema
	w := chair.Width
	h := chair.Height
	d := chair.Depth
	sqlstr = `SELECT * FROM estate where (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) order by view_count desc limit ?`
	err = db.Select(&recommendEstates, sqlstr, w, h, w, d, h, w, h, d, d, w, d, h, limit)
	if err != nil {
		c.Logger().Errorf("Database execution error : %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var re EstateSearchResponse

	for _, estate := range recommendEstates {
		re.Estates = append(re.Estates, *estate.ToEstate())
	}

	return c.JSON(http.StatusOK, re)
}

func searchEstateNazotte(c echo.Context) error {
	coordinates := Coordinates{}
	err := c.Bind(&coordinates)
	if err != nil {
		c.Echo().Logger.Infof("post search estate nazotte failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = coordinates.coordinatesToPolygon()
	if err != nil {
		c.Echo().Logger.Errorf("request coordinates are not WKT Polygon : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	b := coordinates.getBoundingBox()
	estatesInBoundingBox := []EstateSchema{}

	sqlstr := `SELECT * FROM estate WHERE latitude < ? AND latitude > ? AND longitude < ? AND longitude > ?`

	err = db.Select(&estatesInBoundingBox, sqlstr, b.TopLeftCorner.Latitude, b.BottomRightCorner.Latitude, b.BottomRightCorner.Longitude, b.TopLeftCorner.Longitude)
	if err == sql.ErrNoRows {
		c.Echo().Logger.Infof("select * from estate where latitude ...", err)
		return c.NoContent(http.StatusNoContent)
	} else if err != nil {
		c.Echo().Logger.Errorf("database execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	estatesInPolygon := []EstateSchema{}
	for _, estate := range estatesInBoundingBox {
		validatedEstate := EstateSchema{}

		point := fmt.Sprintf("'POINT(%f %f)'", estate.Latitude, estate.Longitude)
		sqlstr := `SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s, %v))`
		sqlstr = fmt.Sprintf(sqlstr, coordinates.coordinatesToText(), point, SRID)

		err = db.Get(&validatedEstate, sqlstr, estate.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				c.Echo().Logger.Errorf("db access is failed on executing validate if estate is in polygon : %v", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		} else {
			estatesInPolygon = append(estatesInPolygon, validatedEstate)
		}
	}

	re := make([]Estate, 0, len(estatesInPolygon))
	for _, estate := range estatesInPolygon {
		re = append(re, *estate.ToEstate())
	}

	return c.JSON(http.StatusOK, re)
}

func postEstateRequestDocument(c echo.Context) error {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func responseEstateRange(c echo.Context) error {
	ranges := RangeResponseEstateMap{
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
		if boundingBox.TopLeftCorner.Longitude > coordinate.Longitude {
			boundingBox.TopLeftCorner.Longitude = coordinate.Longitude
		}

		if boundingBox.BottomRightCorner.Latitude > coordinate.Latitude {
			boundingBox.BottomRightCorner.Latitude = coordinate.Latitude
		}
		if boundingBox.BottomRightCorner.Longitude < coordinate.Longitude {
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
