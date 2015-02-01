package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec(`DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func InsertBookHandler(w http.ResponseWriter, r *http.Request) {
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

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
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

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["book_id"]
	rows, err := db.Query(`SELECT * FROM reviews WHERE book_id = $1`, bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.Id, &review.Description, &review.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, review)
	}
	json.NewEncoder(w).Encode(reviews)
}

func GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["book_id"]
	id := mux.Vars(r)["id"]
	rows, err := db.Query(`SELECT * FROM reviews WHERE book_id = $1 AND id = $2`, bookId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()

	var review Review
	for rows.Next() {
		if err := rows.Scan(&review.Id, &review.Description, &review.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(review)
}

func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["book_id"]
	id := mux.Vars(r)["id"]
	_, err := db.Exec(`DELETE FROM reviews WHERE book_id = $1 AND id = $2`, bookId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func InsertReviewHandler(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["book_id"]
	var review Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := db.Exec(`INSERT INTO reviews (book_id, description, name) VALUES ($1, $2, $3)`, bookId, review.Description, review.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	location := fmt.Sprintf("/books/%d/reviews/foo", bookId)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["book_id"]
	id := mux.Vars(r)["id"]

	rows, err := db.Query(`SELECT * FROM reviews WHERE book_id = $1 AND id = $2`, bookId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	var review Review
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&review.Id, &review.Description, &review.Name)
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
	_, err = db.Exec(`UPDATE reviews SET description = $1, name = $2 WHERE book_id = $3 AND id = $4`, params["description"], params["name"], bookId, id)
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
	r.HandleFunc("/books/{id}", GetBookHandler).Methods("GET")
	r.HandleFunc("/books/{id}", DeleteBookHandler).Methods("DELETE")
	r.HandleFunc("/books/{id}", UpdateBookHandler).Methods("PUT")
	r.HandleFunc("/books", InsertBookHandler).Methods("POST")
	r.HandleFunc("/books/{book_id}/reviews", GetReviewsHandler).Methods("GET")
	r.HandleFunc("/books/{book_id}/reviews/{id}", GetReviewHandler).Methods("GET")
	r.HandleFunc("/books/{book_id}/reviews/{id}", DeleteReviewHandler).Methods("DELETE")
	r.HandleFunc("/books/{book_id}/reviews/{id}", InsertReviewHandler).Methods("POST")
	r.HandleFunc("/books/{book_id}/reviews/{id}", UpdateReviewHandler).Methods("PUT")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
