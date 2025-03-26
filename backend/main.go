package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file")
	}

	//Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "simple_sre_db",
	}

	//Get a database handle.
	var dbErr error
	db, dbErr = sql.Open("mysql", cfg.FormatDSN())
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/count", getCountHandler)
	router.POST("/count", incrementCountHandler)
	router.Run("localhost:8080")
}

func incrementCountHandler(c *gin.Context) {
	id, err := incrementCount(rand.IntN(100), 1)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, id)
	}
}

func getCountHandler(c *gin.Context) {
	count, err := Count()

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, count)
	}

}

func Count() (int, error) {
	var cnt int
	row := db.QueryRow("SELECT SUM(count) as count FROM count_table")
	err := row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("Count: couldn't get count")
	}

	return cnt, nil
}

func incrementCount(slot int, count int) (int64, error) {
	result, err := db.Exec("INSERT INTO count_table (slot, count) VALUES (?, ?)", slot, count)

	if err != nil {
		return 0, fmt.Errorf("incrementCount: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("incrementCount: %v", err)
	}

	return id, nil
}
