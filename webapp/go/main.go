package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const LIMIT = 20
const NAZOTTE_LIMIT = 50

var db *sqlx.DB
var MySQLConnectionData *MySQLConnectionEnv
var chairSearchCondition ChairSearchCondition
var estateSearchCondition EstateSearchCondition

type InitializeResponse struct {
	Language string `json:"language"`
}

type Chair struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Thumbnail   string `db:"thumbnail" json:"thumbnail"`
	Price       int64  `db:"price" json:"price"`
	Height      int64  `db:"height" json:"height"`
	Width       int64  `db:"width" json:"width"`
	Depth       int64  `db:"depth" json:"depth"`
	Color       string `db:"color" json:"color"`
	Features    string `db:"features" json:"features"`
	Kind        string `db:"kind" json:"kind"`
	ViewCount   int64  `db:"view_count" json:"-"`
	Stock       int64  `db:"stock" json:"-"`
}

type ChairSearchResponse struct {
	Count  int64   `json:"count"`
	Chairs []Chair `json:"chairs"`
}

type ChairRecommendResponse struct {
	Chairs []Chair `json:"chairs"`
}

//Estate 物件
type Estate struct {
	ID          int64   `db:"id" json:"id"`
	Thumbnail   string  `db:"thumbnail" json:"thumbnail"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Latitude    float64 `db:"latitude" json:"latitude"`
	Longitude   float64 `db:"longitude" json:"longitude"`
	Address     string  `db:"address" json:"address"`
	Rent        int64   `db:"rent" json:"rent"`
	DoorHeight  int64   `db:"door_height" json:"doorHeight"`
	DoorWidth   int64   `db:"door_width" json:"doorWidth"`
	Features    string  `db:"features" json:"features"`
	ViewCount   int64   `db:"view_count" json:"-"`
}

//EstateSearchResponse estate/searchへのレスポンスの形式
type EstateSearchResponse struct {
	Count   int64    `json:"count"`
	Estates []Estate `json:"estates"`
}

type EstateRecommendResponse struct {
	Estates []Estate `json:"estates"`
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

type RangeCondition struct {
	Prefix string   `json:"prefix"`
	Suffix string   `json:"suffix"`
	Ranges []*Range `json:"ranges"`
}

type ListCondition struct {
	List []string `json:"list"`
}

type EstateSearchCondition struct {
	DoorWidth  RangeCondition `json:"doorWidth"`
	DoorHeight RangeCondition `json:"doorHeight"`
	Rent       RangeCondition `json:"rent"`
	Feature    ListCondition  `json:"feature"`
}

type ChairSearchCondition struct {
	Width   RangeCondition `json:"width"`
	Height  RangeCondition `json:"height"`
	Depth   RangeCondition `json:"depth"`
	Price   RangeCondition `json:"price"`
	Color   ListCondition  `json:"color"`
	Feature ListCondition  `json:"feature"`
	Kind    ListCondition  `json:"kind"`
}

type BoundingBox struct {
	// TopLeftCorner 緯度経度が共に最小値になるような点の情報を持っている
	TopLeftCorner Coordinate
	// BottomRightCorner 緯度経度が共に最大値になるような点の情報を持っている
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

func init() {
	jsonText, err := ioutil.ReadFile("../fixture/chair_condition.json")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(jsonText, &chairSearchCondition)

	jsonText, err = ioutil.ReadFile("../fixture/estate_condition.json")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(jsonText, &estateSearchCondition)
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
	e.POST("/initialize", initialize)

	// Chair Handler
	e.GET("/api/chair/:id", getChairDetail)
	e.GET("/api/chair/search", searchChairs)
	e.GET("/api/chair/search/condition", getChairSearchCondition)
	e.POST("/api/chair/buy/:id", buyChair)

	// Estate Handler
	e.GET("/api/estate/:id", getEstateDetail)
	e.GET("/api/estate/search", searchEstates)
	e.POST("/api/estate/req_doc/:id", postEstateRequestDocument)
	e.POST("/api/estate/nazotte", searchEstateNazotte)
	e.GET("/api/estate/search/condition", getEstateSearchCondition)

	// Recommended Handler
	e.GET("/api/recommended_estate", searchRecommendEstate)
	e.GET("/api/recommended_estate/:id", searchRecommendEstateWithChair)
	e.GET("/api/recommended_chair", searchRecommendChair)

	MySQLConnectionData = NewMySQLConnectionEnv()

	var err error
	db, err = MySQLConnectionData.ConnectDB()
	if err != nil {
		e.Logger.Fatalf("DB connection failed : %v", err)
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

	return c.JSON(http.StatusOK, InitializeResponse{
		Language: "go",
	})
}

func getChairDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Errorf("Request parameter \"id\" parse error : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	chair := Chair{}
	sqlstr := "SELECT * FROM chair WHERE id = ?"
	err = db.Get(&chair, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("requested id's chair not found : %v", id)
			return c.NoContent(http.StatusNotFound)
		}
		c.Echo().Logger.Errorf("Failed to get the chair from id : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	} else if chair.Stock <= 0 {
		c.Echo().Logger.Infof("requested id's chair is sold out : %v", id)
		return c.NoContent(http.StatusNotFound)
	}

	tx, err := db.Begin()
	if err != nil {
		c.Echo().Logger.Errorf("failed to create transaction : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE chair SET view_count = ? WHERE id = ?", chair.ViewCount+1, id)
	if err != nil {
		c.Echo().Logger.Errorf("view_count of chair update failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		c.Echo().Logger.Errorf("transaction commit error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, chair)
}

func searchChairs(c echo.Context) error {
	var chairHeight, chairWidth, chairDepth, chairPrice *Range
	var err error

	searchQueryArray := make([]string, 0)
	queryParams := make([]interface{}, 0)

	if c.QueryParam("priceRangeId") != "" {
		chairPrice, err = getRange(chairSearchCondition.Price, c.QueryParam("priceRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("priceRangeID invalid, %v : %v", c.QueryParam("priceRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

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
		chairHeight, err = getRange(chairSearchCondition.Height, c.QueryParam("heightRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("heightRangeIf invalid, %v : %v", c.QueryParam("heightRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

		if chairHeight.Min != -1 {
			searchQueryArray = append(searchQueryArray, "height >= ? ")
			queryParams = append(queryParams, chairHeight.Min)
		}
		if chairHeight.Max != -1 {
			searchQueryArray = append(searchQueryArray, "height < ? ")
			queryParams = append(queryParams, chairHeight.Max)
		}
	}

	if c.QueryParam("widthRangeId") != "" {
		chairWidth, err = getRange(chairSearchCondition.Width, c.QueryParam("widthRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("widthRangeID invalid, %v : %v", c.QueryParam("widthRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

		if chairWidth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "width >= ? ")
			queryParams = append(queryParams, chairWidth.Min)
		}
		if chairWidth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "width < ? ")
			queryParams = append(queryParams, chairWidth.Max)
		}
	}

	if c.QueryParam("depthRangeId") != "" {
		chairDepth, err = getRange(chairSearchCondition.Depth, c.QueryParam("depthRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("depthRangeId invalid, %v : %v", c.QueryParam("depthRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

		if chairDepth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "depth >= ? ")
			queryParams = append(queryParams, chairDepth.Min)
		}
		if chairDepth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "depth < ? ")
			queryParams = append(queryParams, chairDepth.Max)
		}
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
	}

	if len(searchQueryArray) == 0 {
		c.Echo().Logger.Infof("Search condition not found")
		return c.NoContent(http.StatusBadRequest)
	}

	searchQueryArray = append(searchQueryArray, "stock > 0")

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

	var chairs ChairSearchResponse
	chairs.Chairs = []Chair{}
	sqlstr := "SELECT * FROM chair WHERE "
	searchCondition := strings.Join(searchQueryArray, " AND ")

	limitOffset := " ORDER BY view_count DESC, id ASC LIMIT ? OFFSET ?"

	countsql := "SELECT COUNT(*) FROM chair WHERE "
	err = db.Get(&chairs.Count, countsql+searchCondition, queryParams...)
	if err != nil {
		c.Logger().Errorf("searchChairs DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	searchedchairs := []Chair{}
	queryParams = append(queryParams, perpage, page*perpage)
	err = db.Select(&searchedchairs, sqlstr+searchCondition+limitOffset, queryParams...)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, ChairSearchResponse{Count: 0, Chairs: []Chair{}})
		}
		c.Logger().Errorf("searchChairs DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, c := range searchedchairs {
		chairs.Chairs = append(chairs.Chairs, c)
	}

	return c.JSON(http.StatusOK, chairs)
}

func sendEmail(email string) {
	// Not implemented
}

func buyChair(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	email, ok := m["email"].(string)
	if !ok {
		c.Echo().Logger.Info("post request document failed : email not found in request body")
		return c.NoContent(http.StatusBadRequest)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	var chair Chair
	err = db.Get(&chair, "SELECT * FROM chair WHERE id = ? AND stock > 0", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("buyChair chair id \"%v\" not found", id)
			return c.NoContent(http.StatusNotFound)
		}
		c.Echo().Logger.Errorf("DB Execution Error: on getting a chair by id : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	tx, err := db.Begin()
	if err != nil {
		c.Echo().Logger.Errorf("failed to create transaction : %v", err)
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

	sendEmail(email)
	return c.NoContent(http.StatusOK)
}

func getChairSearchCondition(c echo.Context) error {
	return c.JSON(http.StatusOK, chairSearchCondition)
}

func searchRecommendChair(c echo.Context) error {
	var recommendChairs []Chair

	sqlstr := `SELECT * FROM chair WHERE stock > 0 ORDER BY view_count DESC, id ASC LIMIT ?`

	err := db.Select(&recommendChairs, sqlstr, LIMIT)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error("searchRecommendChair not found")
			return c.JSON(http.StatusOK, ChairRecommendResponse{[]Chair{}})
		}
		c.Logger().Errorf("searchRecommendChair DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var rc ChairRecommendResponse

	for _, chair := range recommendChairs {
		rc.Chairs = append(rc.Chairs, chair)
	}

	return c.JSON(http.StatusOK, rc)
}

func getEstateDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("Request parameter \"id\" parse error : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	var estate Estate
	err = db.Get(&estate, "SELECT * FROM estate WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Echo().Logger.Infof("getEstateDetail estate id %v not found", id)
			return c.NoContent(http.StatusNotFound)
		}
		c.Echo().Logger.Errorf("Database Execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	tx, err := db.Begin()
	if err != nil {
		c.Echo().Logger.Errorf("failed to create transaction : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()

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

	return c.JSON(http.StatusOK, estate)
}

func getRange(cond RangeCondition, rangeID string) (*Range, error) {
	RangeIndex, err := strconv.Atoi(rangeID)
	if err != nil {
		return nil, err
	}

	if RangeIndex < 0 && len(cond.Ranges) <= RangeIndex {
		return nil, fmt.Errorf("Unexpected Range ID")
	}

	return cond.Ranges[RangeIndex], nil
}

func searchEstates(c echo.Context) error {
	var doorHeight, doorWidth, estateRent *Range
	var err error

	searchQueryArray := make([]string, 0)
	var searchQueryParameter []interface{}

	if c.QueryParam("doorHeightRangeId") != "" {
		doorHeight, err = getRange(estateSearchCondition.DoorHeight, c.QueryParam("doorHeightRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("doorHeightRangeID invalid, %v : %v", c.QueryParam("doorHeightRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

		if doorHeight.Min != -1 {
			searchQueryArray = append(searchQueryArray, "door_height >= ? ")
			searchQueryParameter = append(searchQueryParameter, doorHeight.Min)
		}
		if doorHeight.Max != -1 {
			searchQueryArray = append(searchQueryArray, "door_height < ? ")
			searchQueryParameter = append(searchQueryParameter, doorHeight.Max)
		}
	}

	if c.QueryParam("doorWidthRangeId") != "" {
		doorWidth, err = getRange(estateSearchCondition.DoorWidth, c.QueryParam("doorWidthRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("doorWidthRangeID invalid, %v : %v", c.QueryParam("doorWidthRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

		if doorWidth.Min != -1 {
			searchQueryArray = append(searchQueryArray, "door_width >= ? ")
			searchQueryParameter = append(searchQueryParameter, doorWidth.Min)
		}
		if doorWidth.Max != -1 {
			searchQueryArray = append(searchQueryArray, "door_width < ? ")
			searchQueryParameter = append(searchQueryParameter, doorWidth.Max)
		}
	}

	if c.QueryParam("rentRangeId") != "" {
		estateRent, err = getRange(estateSearchCondition.Rent, c.QueryParam("rentRangeId"))
		if err != nil {
			c.Echo().Logger.Infof("rentRangeID invalid, %v : %v", c.QueryParam("rentRangeId"), err)
			return c.NoContent(http.StatusBadRequest)
		}

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
	}

	if len(searchQueryArray) == 0 {
		c.Echo().Logger.Infof("searchEstates search condition not found")
		return c.NoContent(http.StatusBadRequest)
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
	estates.Estates = []Estate{}
	sqlstr := "SELECT * FROM estate WHERE "
	searchQuery := strings.Join(searchQueryArray, " AND ")

	limitOffset := " ORDER BY view_count DESC, id ASC LIMIT ? OFFSET ?"

	countsql := "SELECT COUNT(*) FROM estate WHERE "
	err = db.Get(&estates.Count, countsql+searchQuery, searchQueryParameter...)
	if err != nil {
		c.Logger().Errorf("searchEstates DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	matchestates := []Estate{}
	searchQueryParameter = append(searchQueryParameter, perpage, page*perpage)
	err = db.Select(&matchestates, sqlstr+searchQuery+limitOffset, searchQueryParameter...)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, EstateSearchResponse{Count: 0, Estates: []Estate{}})
		}
		c.Logger().Errorf("searchEstates DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, e := range matchestates {
		estates.Estates = append(estates.Estates, e)
	}

	return c.JSON(http.StatusOK, estates)
}

func searchRecommendEstate(c echo.Context) error {
	recommendEstates := make([]Estate, 0, LIMIT)

	sqlstr := `SELECT * FROM estate ORDER BY view_count DESC, id ASC LIMIT ?`

	err := db.Select(&recommendEstates, sqlstr, LIMIT)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error("searchRecommendEstate not found")
			return c.JSON(http.StatusOK, EstateRecommendResponse{[]Estate{}})
		}
		c.Logger().Errorf("searchRecommendEstate DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var re EstateRecommendResponse

	for _, estate := range recommendEstates {
		re.Estates = append(re.Estates, estate)
	}

	return c.JSON(http.StatusOK, re)
}

func searchRecommendEstateWithChair(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Infof("Invalid format searchRecommendedEstateWithChair id : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	chair := Chair{}
	sqlstr := `SELECT * FROM chair WHERE id = ?`

	err = db.Get(&chair, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Infof("Requested chair id \"%v\" not found", id)
			return c.NoContent(http.StatusBadRequest)
		}
		c.Logger().Errorf("Database execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var recommendEstates []Estate
	w := chair.Width
	h := chair.Height
	d := chair.Depth
	sqlstr = `SELECT * FROM estate where (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) ORDER BY view_count DESC, id ASC LIMIT ?`
	err = db.Select(&recommendEstates, sqlstr, w, h, w, d, h, w, h, d, d, w, d, h, LIMIT)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, EstateRecommendResponse{[]Estate{}})
		}
		c.Logger().Errorf("Database execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var re EstateRecommendResponse

	for _, estate := range recommendEstates {
		re.Estates = append(re.Estates, estate)
	}

	return c.JSON(http.StatusOK, re)
}

func searchEstateNazotte(c echo.Context) error {
	coordinates := Coordinates{}
	err := c.Bind(&coordinates)
	if err != nil {
		c.Echo().Logger.Infof("post search estate nazotte failed : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if len(coordinates.Coordinates) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	b := coordinates.getBoundingBox()
	estatesInBoundingBox := []Estate{}

	sqlstr := `SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY view_count DESC, id ASC`

	err = db.Select(&estatesInBoundingBox, sqlstr, b.BottomRightCorner.Latitude, b.TopLeftCorner.Latitude, b.BottomRightCorner.Longitude, b.TopLeftCorner.Longitude)
	if err == sql.ErrNoRows {
		c.Echo().Logger.Infof("select * from estate where latitude ...", err)
		return c.JSON(http.StatusOK, EstateSearchResponse{Count: 0, Estates: []Estate{}})
	} else if err != nil {
		c.Echo().Logger.Errorf("database execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	estatesInPolygon := []Estate{}
	for _, estate := range estatesInBoundingBox {
		validatedEstate := Estate{}

		point := fmt.Sprintf("'POINT(%f %f)'", estate.Latitude, estate.Longitude)
		sqlstr := `SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))`
		sqlstr = fmt.Sprintf(sqlstr, coordinates.coordinatesToText(), point)

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

	var re EstateSearchResponse
	re.Estates = []Estate{}
	for i, estate := range estatesInPolygon {
		if i >= NAZOTTE_LIMIT {
			break
		}
		re.Estates = append(re.Estates, estate)
	}
	re.Count = int64(len(re.Estates))

	return c.JSON(http.StatusOK, re)
}

func postEstateRequestDocument(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	email, ok := m["email"].(string)
	if !ok {
		c.Echo().Logger.Info("post request document failed : email not found in request body")
		return c.NoContent(http.StatusBadRequest)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Infof("post request document failed : %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	sqlstr := `SELECT * FROM estate WHERE id = ?`

	estate := Estate{}
	err = db.Get(&estate, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		c.Logger().Errorf("postEstateRequestDocument DB execution error : %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	sendEmail(email)
	return c.NoContent(http.StatusOK)
}

func getEstateSearchCondition(c echo.Context) error {
	return c.JSON(http.StatusOK, estateSearchCondition)
}

func (cs Coordinates) getBoundingBox() BoundingBox {
	coordinates := cs.Coordinates
	boundingBox := BoundingBox{
		TopLeftCorner: Coordinate{
			Latitude: coordinates[0].Latitude, Longitude: coordinates[0].Longitude,
		},
		BottomRightCorner: Coordinate{
			Latitude: coordinates[0].Latitude, Longitude: coordinates[0].Longitude,
		},
	}
	for _, coordinate := range coordinates {
		if boundingBox.TopLeftCorner.Latitude > coordinate.Latitude {
			boundingBox.TopLeftCorner.Latitude = coordinate.Latitude
		}
		if boundingBox.TopLeftCorner.Longitude > coordinate.Longitude {
			boundingBox.TopLeftCorner.Longitude = coordinate.Longitude
		}

		if boundingBox.BottomRightCorner.Latitude < coordinate.Latitude {
			boundingBox.BottomRightCorner.Latitude = coordinate.Latitude
		}
		if boundingBox.BottomRightCorner.Longitude < coordinate.Longitude {
			boundingBox.BottomRightCorner.Longitude = coordinate.Longitude
		}
	}
	return boundingBox
}

func (cs Coordinates) coordinatesToText() string {
	PolygonArray := make([]string, 0, len(cs.Coordinates))
	for _, c := range cs.Coordinates {
		PolygonArray = append(PolygonArray, fmt.Sprintf("%f %f", c.Latitude, c.Longitude))
	}
	return fmt.Sprintf("'POLYGON((%s))'", strings.Join(PolygonArray, ","))
}
