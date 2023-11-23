package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"example/data-access/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./recordings.db")

	log.Println("Creating database...")
	db, err := sql.Open("sqlite3", "./recordings.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("ping error")

		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	albums := database.NewAlbums(db)

	found, err := albums.AlbumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", found)

	albID, err := albums.Add(database.Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)

	found, err = albums.AlbumsByArtist("Betty Carter")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", found)

	alb, err := albums.AlbumByID(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)
}
