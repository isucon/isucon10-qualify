package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var (
	db *sqlx.DB
)

//EstateSchema estate tableに格納されている物件データ
type EstateSchema struct {
	ID          int64   `db:"id"`
	Thumbnails  string  `db:"thumbnails"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Address     string  `db:"address"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
	DoorHeight  int64   `db:"door_height"`
	DoorWidth   int64   `db:"door_width"`
	Rent        int64   `db:"rent"`
	Features    string  `db:"features"`
	ViewCount   int64   `db:"view_count"`
}

func (es EstateSchema) ToEstate() *Estate {
	return &Estate{
		ID:          es.ID,
		Thumbnails:  es.Thumbnails,
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
	Thumbnails  string  `json:"thumbnails"`
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

	// Routes
	e.GET("/api/estate/:id", getEstateDetail)

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
