package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

type Albums struct {
	albums []album
}

func NewAlbums() *Albums {
	return &Albums{
		albums: []album{
			{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: "56.99"},
			{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: "17.99"},
			{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: "39.99"},
		},
	}
}

var albums = NewAlbums()

func (a *Albums) getAlbums() *[]album {
	return &a.albums
}

func (a *Albums) getAlbum(albumID string) *album {
	for _, album := range a.albums {
		if album.ID == albumID {
			return &album
		}
	}
	return nil
}

func (a *Albums) deleteAlbum(albumID string) {
	for i, album := range a.albums {
		if album.ID == albumID {
			a.albums = append(a.albums[:i], a.albums[i+1:]...)
		}
	}
}

func (a *Albums) createAlbum(newAlbum album) {
	a.albums = append(a.albums, newAlbum)
}

func getAlbums(w http.ResponseWriter, r *http.Request) {
	albs := albums.getAlbums()
	bytes, err := json.Marshal(albs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func postAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newAlbum album
	err := json.NewDecoder(r.Body).Decode(&newAlbum)
	if err != nil {
		log.Printf("Error decoding album: %v", err)
	}

	albums.createAlbum(newAlbum)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	albumID := chi.URLParam(r, "albumID")
	albums.deleteAlbum(albumID)
}

func getAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	albumID := chi.URLParam(r, "albumID")

	album := albums.getAlbum(albumID)
	if album == nil {
		return
	}

	bytes, err := json.Marshal(album)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(bytes)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/albums", getAlbums)
	r.Get("/albums/{albumID}", getAlbum)
	r.Delete("/albums/{albumID}", deleteAlbum)
	r.Post("/albums", postAlbums)

	if err := http.ListenAndServe(":8080", r); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
