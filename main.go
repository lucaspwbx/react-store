package main

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Review struct {
	Id          string `bson:"_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Book struct {
	Id       string   `json:"id"`
	Title    string   `json:"title"`
	Pages    int      `json:"pages"`
	Language string   `json:"languages"`
	Reviews  []Review `json:"reviews"`
}

var db *sql.DB

func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Title, &book.Language, &book.Pages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		books = append(books, book)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	rows, err := db.Query(`SELECT * FROM books WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	var book Book
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&book.Id, &book.Title, &book.Language, &book.Pages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=lucasweiblen dbname=bookreviewer sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books", GetBooksHandler).Methods("GET")
	r.HandleFunc("/books/{id}", GetBookByIdHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func newID() string {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return ""
	}
	nanoUnix := time.Now().UnixNano()
	hash := sha512.New()
	hash.Write(buf[:])
	hash.Write([]byte(fmt.Sprintf("%d", nanoUnix)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
