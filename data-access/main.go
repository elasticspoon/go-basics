package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	Title  string
	Artist string
	ID     int
	Price  float32
}

func main() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.DBName = "recordings"
	cfg.Addr = "0.0.0.0:3306"

	// var dp *sql.DB
	// var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to DB")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("ping error")

		log.Fatal(pingErr)
	}

	fmt.Println("Successfully connected!")

	albums, err := albumsByArtist("John Coltrane", db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

func albumByID(id int, db *sql.DB) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID(%d): no such album", id)
		}
		return alb, fmt.Errorf("albumsByID %d: %v", id, err)
	}

	return alb, nil
}

func albumsByArtist(artist string, db *sql.DB) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?;", artist)
	if err != nil {
		// add specific error if album with ID not found
		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
		}
		albums = append(albums, alb)
	}

	return albums, nil
}

func addAlbum(alb Album, db *sql.DB) (int64, error) {
	res, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", alb)
	}

	return id, nil
}
