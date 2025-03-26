package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	cnt, cntErr := Count()
	if cntErr != nil {
		log.Fatal(cntErr)
	}
	fmt.Printf("Count found: %v\n", cnt)
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
