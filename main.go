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
	Id          string `json:"id, omitempty"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Book struct {
	Id       string   `json:"id, omitempty"`
	Title    string   `json:"title"`
	Pages    int      `json:"pages"`
	Language string   `json:"language"`
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
			return
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
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func DeleteBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec(`DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func InsertBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := db.Exec(`INSERT INTO books (title, language, pages) VALUES ($1, $2, $3)`, book.Title, book.Language, book.Pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/books/blabla")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBookByIdHandler(w http.ResponseWriter, r *http.Request) {
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
			return
		}
	}
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec(`UPDATE books SET title = $1, language = $2, pages = $3 WHERE id = $4`, params["title"], params["language"], params["pages"], id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	r.HandleFunc("/books/{id}", DeleteBookByIdHandler).Methods("DELETE")
	r.HandleFunc("/books/{id}", UpdateBookByIdHandler).Methods("PUT")
	r.HandleFunc("/books", InsertBook).Methods("POST")
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
