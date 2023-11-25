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

func postAlbumsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newAlbum database.Album
	err := json.NewDecoder(r.Body).Decode(&newAlbum)
	if err != nil {
		log.Printf("Error decoding album: %v", err)
		http.Redirect(w, r, "/albums/new", http.StatusSeeOther)
	}

	id, err := albums.Add(newAlbum)
	if err != nil {
		log.Printf("Error adding album: %v", err)
		http.Redirect(w, r, "/albums/new", http.StatusSeeOther)
	}

	http.Redirect(w, r, fmt.Sprintf("/albums/%d", id), http.StatusSeeOther)
}

func postAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newAlbum database.Album
	newAlbum.Title = template.HTMLEscapeString(r.FormValue("title"))
	newAlbum.Artist = template.HTMLEscapeString(r.FormValue("artist"))

	price := r.FormValue("price")
	val, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Printf("Error decoding album: %v", err)
		http.Redirect(w, r, "/albums/new", http.StatusSeeOther)
	}

	newAlbum.Price = float32(val)

	id, err := albums.Add(newAlbum)
	if err != nil {
		log.Printf("Error adding album: %v", err)
		http.Redirect(w, r, "/albums/new", http.StatusSeeOther)
	}

	http.Redirect(w, r, fmt.Sprintf("/albums/%d", id), http.StatusSeeOther)
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
	if err != nil {
		http.Redirect(w, r, "/albums", http.StatusSeeOther)
		catch(err)
	}

	album, err := albums.AlbumByID(id)
	if err != nil {
		http.Redirect(w, r, "/albums", http.StatusSeeOther)
		catch(err)
	}

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

func showIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Hello World!"))
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
	r.Get("/alive", serverAlive)
	r.Route("/albums", func(r chi.Router) {
		r.Get("/", getAlbums)
		r.Get("/new", newAlbum)
		r.Get("/{albumID}/edit", editAlbum)
		r.Post("/", postAlbum)
		r.Get("/{albumID}", getAlbum)
		r.Delete("/{albumID}", deleteAlbum)
	})
	// r.Get("/albums", getAlbumsApi)

	err = http.ListenAndServe(":8080", r)
	catch(err)
}

func editAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.Header().Set("Content-Type", "text/html")
	albumID := chi.URLParam(r, "albumID")

	id, err := strconv.Atoi(albumID)
	if err != nil {
		http.Redirect(w, r, "/albums", http.StatusSeeOther)
		catch(err)
	}

	album, err := albums.AlbumByID(id)
	if err != nil {
		http.Redirect(w, r, "/albums", http.StatusSeeOther)
		catch(err)
	}

	t, err := template.ParseFiles("templates/base.html", "templates/albums/new.html", "templates/albums/form.html")
	catch(err)
	err = t.Execute(w, album)
	catch(err)
}

func newAlbum(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/base.html", "templates/albums/new.html", "templates/albums/form.html")
	catch(err)
	err = t.Execute(w, nil)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
