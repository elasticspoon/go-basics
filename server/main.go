package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"example/data-access/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var albums *database.Albums

func getAlbums(w http.ResponseWriter, _ *http.Request) {
	albs, err := albums.All()
	catch(err)

	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/base.html", "templates/index.html")
	catch(err)
	err = t.Execute(w, albs)
	catch(err)
}

func getAlbumsApi(w http.ResponseWriter, _ *http.Request) {
	albs, err := albums.All()
	if err != nil {
		log.Printf("Error getting albums: %v", err)
	}
	bytes, err := json.Marshal(albs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func postAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newAlbum database.Album
	err := json.NewDecoder(r.Body).Decode(&newAlbum)
	if err != nil {
		log.Printf("Error decoding album: %v", err)
	}

	albums.Add(newAlbum)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	albumID := chi.URLParam(r, "albumID")

	id, err := strconv.Atoi(albumID)
	if err != nil {
		log.Printf("Error converting albumID to int: %v", err)
	}

	albums.Delete(id)
}

func getAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	albumID := chi.URLParam(r, "albumID")

	id, err := strconv.Atoi(albumID)
	catch(err)

	album, err := albums.AlbumByID(id)
	catch(err)

	t, err := template.ParseFiles("templates/base.html", "templates/album.html")
	catch(err)

	err = t.Execute(w, album)
	catch(err)
}

func serverAlive(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		catch(err)
	}
}

func getAlbumApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	albumID := chi.URLParam(r, "albumID")

	id, err := strconv.Atoi(albumID)
	if err != nil {
		log.Printf("Error converting albumID to int: %v", err)
	}

	album, err := albums.AlbumByID(id)
	if err != nil {
		log.Printf("Error getting album by ID: %v", err)
		return
	}

	bytes, err := json.Marshal(album)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(bytes)
}

func showIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

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

	albums = database.NewAlbums(db)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/livereload.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./js/livereload.js")
	})
	r.Get("/", showIndex)
	r.Get("/albums", getAlbums)
	r.Get("/alive", serverAlive)
	// r.Get("/albums", getAlbumsApi)
	r.Get("/albums/{albumID}", getAlbum)
	r.Delete("/albums/{albumID}", deleteAlbum)
	r.Post("/albums", postAlbums)

	err = http.ListenAndServe(":8080", r)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
