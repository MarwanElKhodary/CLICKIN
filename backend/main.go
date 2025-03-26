package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
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
	fmt.Printf("Original count: %v\n", cnt)

	fmt.Printf("Incrementing count by 1!\n")
	incrementCount(rand.IntN(100), 1)

	cnt2, cntErr2 := Count()
	if cntErr2 != nil {
		log.Fatal(cntErr2)
	}
	fmt.Printf("New count: %v\n", cnt2)
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
