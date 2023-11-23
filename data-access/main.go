package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var dp *sql.DB

type Album struct {
  ID     int
  Title  string
  Artist string
  Price  float32
}

func main() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.DBName = "recordings"
	cfg.Addr = "0.0.0.0:3306"

	var err error
	dp, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to DB")
		log.Fatal(err)
	}

	pingErr := dp.Ping()
	if pingErr != nil {
		fmt.Println("ping error")

		log.Fatal(pingErr)
	}

	fmt.Println("Successfully connected!")

	rows, err := dp.Query("SELECT ID FROM album;")
	if err != nil {
		log.Fatal(err)
	}

  for rows.Next() {
    id := 0
    rows.Scan()
}
