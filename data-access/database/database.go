package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Albums struct {
	db *sql.DB
}

func NewAlbums(db *sql.DB) *Albums {
	setupTable(db)
	return &Albums{
		db: db,
	}
}

func (a *Albums) Add(album Album) (int64, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(`
    INSERT INTO album (title, artist, price) VALUES (?, ?, ?);
    `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	price := int(album.Price * 100)
	res, err := stmt.Exec(album.Title, album.Artist, price)
	if err != nil {
		// tx.Rollback() I don't think this is needed
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", album)
	}
	return id, nil
}

func (a *Albums) AlbumsByArtist(artist string) ([]Album, error) {
	db := a.db
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?;", artist)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("albumsByArtist(%q): no such artist", artist)
		}
		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		var price int
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
		}
		alb.Price = float32(price) / 100

		albums = append(albums, alb)
	}

	return albums, nil
}

func (a *Albums) AlbumByID(id int) (Album, error) {
	db := a.db
	var alb Album

	stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?;")
	if err != nil {
		return alb, fmt.Errorf("getAlbumByID(%d): %v", id, err)
	}

	err = stmt.QueryRow(id).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("getAlbumByID(%d): no such album", id)
		}
		return alb, fmt.Errorf("getAlbumByID(%d): %v", id, err)
	}

	return alb, nil
}

type Album struct {
	Title  string
	Artist string
	Price  float32
	ID     int
}

func (alb *Album) String() string {
	return fmt.Sprintf("Title:\t%s\nArtist:\t%s\nPrice:\t%.2f\nID:\t%d\n", alb.Title, alb.Artist, alb.Price, alb.ID)
}

func insertAlbum(db *sql.DB, alb Album) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
    INSERT INTO album (title, artist, price) VALUES (?, ?, ?);
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// log.Printf("Inserting album: %v\n", alb)
	price := int(alb.Price * 100)
	res, err := stmt.Exec(alb.Title, alb.Artist, price)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", alb)
	}

	return id, nil
}

func setupTable(db *sql.DB) {
	sqlStmt := `
    CREATE TABLE album (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      artist TEXT NOT NULL,
      price INTEGER NOT NULL
    );
	 `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	albums := []Album{
		{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
		{Title: "Small Groups", Artist: "Art Blakey", Price: 21.99},
		{Title: "Newk's Time", Artist: "Sonny Rollins", Price: 34.99},
	}

	for _, alb := range albums {
		_, err := insertAlbum(db, alb)
		if err != nil {
			log.Fatal(err)
		}
	}
}
