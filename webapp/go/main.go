package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var (
	db   *sqlx.DB
	SRID = 6668
)

var estateRentRanges = []*RangeInt{
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

var estateDoorHeightRanges = []*RangeInt{
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

var estateDoorWidthRanges = []*RangeInt{
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

type RangeFloat struct {
	Min float64
	Max float64
}

type RangeInt struct {
	ID  int64 `json:"id"`
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type RangeResponse struct {
	Prefix string      `json:"prefix"`
	Suffix string      `json:"suffix"`
	Ranges []*RangeInt `json:"ranges"`
}

type RangeResponseMap struct {
	DoorWidth  RangeResponse `json:"doorWidth"`
	DoorHeight RangeResponse `json:"doorHeight"`
	Rent       RangeResponse `json:"rent"`
}

type BoundingBox struct {
	LatitudeRange  RangeFloat
	LongitudeRange RangeFloat
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
	return sqlx.Connect("mysql", dsn)
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

	_, err = db.Exec("UPDATE estate SET view_count = ? WHERE id = ?", estate.ViewCount+1, id)
	if err != nil {
		c.Echo().Logger.Debug("view_count update failed :", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, estate.ToEstate())
}

func getDoorHeightRange(heightRangeID string) (specifyRange *RangeInt, err error) {
	specifyRange = nil

	heightRangeIndex, err := strconv.Atoi(heightRangeID)
	if err != nil {
		return
	}

	switch heightRangeIndex {
	case 0:
		specifyRange = estateDoorWidthRanges[heightRangeIndex]
	case 1:
		specifyRange = estateDoorWidthRanges[heightRangeIndex]
	case 2:
		specifyRange = estateDoorWidthRanges[heightRangeIndex]
	case 3:
		specifyRange = estateDoorWidthRanges[heightRangeIndex]
	default:
		err = fmt.Errorf("Unexpected DoorHeight Range")
	}

	return
}

func getDoorWidthRange(widthRangeID string) (specifyRange *RangeInt, err error) {
	specifyRange = nil

	widthRangeIndex, err := strconv.Atoi(widthRangeID)
	if err != nil {
		return
	}

	switch widthRangeIndex {
	case 0:
		specifyRange = estateDoorWidthRanges[widthRangeIndex]
	case 1:
		specifyRange = estateDoorWidthRanges[widthRangeIndex]
	case 2:
		specifyRange = estateDoorWidthRanges[widthRangeIndex]
	case 3:
		specifyRange = estateDoorWidthRanges[widthRangeIndex]
	default:
		err = fmt.Errorf("Unexpected DoorWidth Range")
	}

	return
}

func getRentRange(rent string) (specifyRange *RangeInt, err error) {
	specifyRange = nil

	rentRangeIndex, err := strconv.Atoi(rent)
	if err != nil {
		return
	}

	switch rentRangeIndex {
	case 0:
		specifyRange = estateRentRanges[rentRangeIndex]
	case 1:
		specifyRange = estateRentRanges[rentRangeIndex]
	case 2:
		specifyRange = estateRentRanges[rentRangeIndex]
	case 3:
		specifyRange = estateRentRanges[rentRangeIndex]
	default:
		err = fmt.Errorf("Unexpected Rent Range")
	}

	return

}

func searchEstates(c echo.Context) error {
	var optwidth, optheight, optfeature, optrent bool
	var doorHeight, doorWidth, estateRent *RangeInt
	estateFeatures := []string{}
	var err error

	if c.QueryParam("doorHeightRangeId") != "" {
		doorHeight, err = getDoorHeightRange(c.QueryParam("doorHeightRangeId"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		optheight = true
	}

	if c.QueryParam("doorWidthRangeId") != "" {
		doorWidth, err = getDoorWidthRange(c.QueryParam("doorWidthRangeId"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		optwidth = true
	}

	if c.QueryParam("features") != "" {
		estateFeatures = strings.Split(c.QueryParam("features"), ",")
		optfeature = true
	}

	if c.QueryParam("rentRangeId") != "" {
		estateRent, err = getRentRange(c.QueryParam("rentRangeId"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		optrent = true
	}

	if !(optheight || optwidth || optfeature || optrent) {
		return c.String(http.StatusBadRequest, "search condition not found")
	}

	var searchquery string
	if optheight {
		hmin := doorHeight.Min
		hmax := doorHeight.Max
		if hmin == -1 {
			searchquery += fmt.Sprintf("door_height < %v ", hmax)
		} else if hmax == -1 {
			searchquery += fmt.Sprintf("door_height >= %v ", hmin)
		} else {
			searchquery += fmt.Sprintf("door_height >= %v and door_height < %v", hmin, hmax)
		}
	}

	if optheight && optwidth {
		searchquery += "and "
	}

	if optwidth {
		wmin := doorWidth.Min
		wmax := doorWidth.Max
		if wmin == -1 {
			searchquery += fmt.Sprintf("door_width < %v ", wmax)
		} else if wmax == -1 {
			searchquery += fmt.Sprintf("door_width >= %v ", wmin)
		} else {
			searchquery += fmt.Sprintf("door_width >= %v and door_width < %v ", wmin, wmax)
		}
	}

	if (optheight || optwidth) && optrent {
		searchquery += "and "
	}

	if optrent {
		rmin := estateRent.Min
		rmax := estateRent.Max
		if rmin == -1 {
			searchquery += fmt.Sprintf("rent < %v ", rmax)
		} else if rmax == -1 {
			searchquery += fmt.Sprintf("rent >= %v ", rmin)
		} else {
			searchquery += fmt.Sprintf("rent >= %v and rent < %v", rmin, rmax)
		}
	}

	var estates EstateSearchResponse
	sqlstr := "select * from estate where "

	if optfeature {
		for _, f := range estateFeatures {
			var likequery string
			if optheight || optwidth || optrent {
				likequery = "and "
			}
			likequery += fmt.Sprintf("features like '%%%v%%'", f)

			rows, err := db.Queryx(sqlstr + searchquery + likequery)
			if err != nil {
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			for rows.Next() {
				e := &EstateSchema{}
				uniq := true
				err := rows.StructScan(e)
				if err != nil {
					c.Logger().Error(err)
					return c.String(http.StatusInternalServerError, err.Error())
				}
				for _, est := range estates.Estates {
					if reflect.DeepEqual(est, e) {
						uniq = false
					}
				}
				if uniq {
					estates.Estates = append(estates.Estates, e.ToEstate())
				}
			}
		}
	} else {
		rows, err := db.Queryx(sqlstr + searchquery)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		for rows.Next() {
			e := EstateSchema{}
			err := rows.StructScan(&e)
			if err != nil {
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			estates.Estates = append(estates.Estates, e.ToEstate())
		}
	}

	return c.JSON(http.StatusOK, estates)
}

func searchRecommendEstate(c echo.Context) error {
	limit := 20
	recommentEstates := make([]Estate, 0, limit)

	sqlstr := `select * from estate order by view_count desc limit ?`

	rows, err := db.Queryx(sqlstr, limit)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	for rows.Next() {
		e := EstateSchema{}
		if err := rows.StructScan(&e); err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		recommentEstates = append(recommentEstates, e.ToEstate())
	}

	return c.JSON(http.StatusOK, recommentEstates)
}

//TODO: グラハムスキャンの実装
func coordinatesToPolygon(coordinates Coordinates) error {
	// グラハムスキャンして、Polygonにして返す

	return nil
}

func getBoundingBox(coordinates []Coordinate) BoundingBox {
	//latitude := make([]float64)
	//ここをソートにしてみると n log(n)になりそう
	boundingBox := BoundingBox{
		LatitudeRange: RangeFloat{
			Min: coordinates[0].Latitude, Max: coordinates[0].Latitude,
		},
		LongitudeRange: RangeFloat{
			Min: coordinates[0].Longitude, Max: coordinates[0].Longitude,
		},
	}
	for _, coordinate := range coordinates {
		if boundingBox.LatitudeRange.Min > coordinate.Latitude {
			boundingBox.LatitudeRange.Min = coordinate.Latitude
		}
		if boundingBox.LatitudeRange.Max < coordinate.Latitude {
			boundingBox.LatitudeRange.Max = coordinate.Latitude
		}
		if boundingBox.LongitudeRange.Min > coordinate.Longitude {
			boundingBox.LongitudeRange.Min = coordinate.Longitude
		}
		if boundingBox.LongitudeRange.Max < coordinate.Longitude {
			boundingBox.LongitudeRange.Max = coordinate.Longitude
		}
	}
	return boundingBox
}

func coordinatesToText(coordinates Coordinates) string {
	// return such as POLYGON((35 137,35 140,37 140, 37 137,35 137)),6668)
	// for _, c := range coordinates { fmt.Spritf("")	}
	PolygonArray := make([]string, 0, len(coordinates.Coordinates))
	for _, coordinate := range coordinates.Coordinates {
		PolygonArray = append(PolygonArray, fmt.Sprintf("%f %f", coordinate.Latitude, coordinate.Longitude))
	}
	return fmt.Sprintf("'POLYGON((%s))', %d", strings.Join(PolygonArray, ","), SRID)
}

func searchEstateNazotte(c echo.Context) error {
	coordinates := Coordinates{}
	err := c.Bind(&coordinates)
	if err != nil {
		c.Echo().Logger.Debug("post search estate nazotte failed :", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = coordinatesToPolygon(coordinates)
	if err != nil {
		c.Echo().Logger.Debug("request coordinates are not WKT Polygon", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	boundingBox := getBoundingBox(coordinates.Coordinates)
	estates := []EstateSchema{}
	respEstates := []EstateSchema{}

	err = db.Select(&estates, "SELECT * FROM estate WHERE latitude < ? AND latitude > ? AND longitude< ? AND longitude> ?", boundingBox.LatitudeRange.Max, boundingBox.LatitudeRange.Min, boundingBox.LongitudeRange.Max, boundingBox.LongitudeRange.Min)
	if err == sql.ErrNoRows {
		c.Echo().Logger.Debug("select * from estate where latitude ...", err)
	}

	for _, estate := range estates {
		validatedEstate := EstateSchema{}
		point := "'" + fmt.Sprintf("POINT(%f %f)", estate.Latitude, estate.Longitude) + "'"
		polygonValidateSQL := fmt.Sprintf("SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s, %v))", coordinatesToText(coordinates), point, SRID)
		err = db.Get(&validatedEstate, polygonValidateSQL, estate.ID)
		if err == sql.ErrNoRows {
			//c.Echo().Logger.Debug("No Rows")
		} else if err != nil {
			c.Echo().Logger.Debug("db access is failed on executing validate estate isWithinPolygon", err)
			return c.NoContent(http.StatusInternalServerError)
		} else {
			respEstates = append(respEstates, validatedEstate)
		}
	}

	re := make([]Estate, 0, len(respEstates))
	for _, estate := range respEstates {
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
