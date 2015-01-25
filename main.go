package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	Id          string `bson:"_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Book struct {
	Id       string   `bson:"_id"`
	Name     string   `json:"name"`
	Pages    int      `json:"pages"`
	Language string   `json:"languages"`
	ISBN     string   `json:"isbn"`
	Reviews  []Review `json:"reviews"`
}

//func (b Book) AddReview(description string, user User) error {
//session, err := getSession()
//if err != nil {
//log.Fatalln("Error opening session")
//}
//defer session.Close()
//users := session.DB("bookstore").C("users")

//review := &Review{
//Id:          newID(),
//Description: description,
//	User:        user,
//}
//reviews := session.DB("bookstore").C("reviews")
//err = reviews.Insert(review)
//if err != nil {
//return errors.New("Error inserting review")
//}
//b.Reviews = append(b.Reviews, *review)
//err = UpdateBookById(b.Id, bson.M{"$set": bson.M{"reviews": b.Reviews}})
//if err != nil {
//return errors.New("Updating error")
//}
//return nil
//}

//now using mux
func GetReviewByIdHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("reviews")

	var review Review
	err = c.FindId(id).One(&review)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(review)
}

//now using mux
func DeleteReviewByIdHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("reviews")

	err = c.RemoveId(id)
	if err != nil {
		http.Error(res, err.Error(), 422)
	}
	res.WriteHeader(http.StatusNoContent)
}

//now using mux
func UpdateReviewById(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id given", http.StatusBadRequest)
	}

	// TODO -> remove after second handling is working
	//var review Review
	//err := json.NewDecoder(req.Body).Decode(&review)
	//if err != nil {
	//http.Error(res, "Error decoding JSON", http.StatusBadRequest)
	//}

	//var params map[string]interface{}
	//params["description"] = review.Description
	//params["name"] = review.Name

	var params map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(res, "Error decoding JSON", http.StatusBadRequest)
	}

	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("reviews")
	if err = c.UpdateId(id, bson.M{"$set": params}); err != nil {
		http.Error(res, err.Error(), 422)
	}
	res.WriteHeader(http.StatusOK)
}

//improved
func NewBookHandler(res http.ResponseWriter, req *http.Request) {
	var book Book
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		http.Error(res, "Error decoding json", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	c := session.DB("bookstore").C("books")
	book.Id = newID()
	if err = c.Insert(book); err != nil {
		http.Error(res, err.Error(), 422)
	}
	location := fmt.Sprintf("/books/%s", book.Id)
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Location", location)
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(book)
}

//mux improved
func DeleteBookByIdHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")

	err = c.RemoveId(id)
	if err != nil {
		http.Error(res, err.Error(), 422)
	}
	res.WriteHeader(http.StatusNoContent)
}

func GetBookByIdHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")

	var book Book
	err = c.FindId(id).One(&book)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(book)
}

func GetBookByNameHandler(res http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	if name == "" {
		http.Error(res, "No name has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")

	var book Book
	err = c.Find(bson.M{"name": name}).One(&book)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(book)
}

func DeleteBookByNameHandler(res http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	if name == "" {
		http.Error(res, "No name has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")

	err = c.Remove(bson.M{"name": name})
	if err != nil {
		http.Error(res, err.Error(), 422)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusNoContent)
}

func UpdateBookByIdHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
	}
	var params map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(res, "Problems decoding JSON", 422)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	//c := session.DB("bookstore").C("books")

	//jupdate := bson.M{"$set": bson.M{params}}
	//if err = c.UpdateId(id, update); err != nil {
	//http.Error(res, err.Error(), 422)
	//}
	//res.WriteHeader(http.StatusOK)
}

//func AddReview() error {
//session, err := getSession()
//if err != nil {
//log.Fatalln("Error opening session")
//}
//defer session.Close()
//users := session.DB("bookstore").C("users")
//reviews := session.DB("bookstore").C("reviews")
//user := &User{Id: newID(), Name: "Renato"}
//err = users.Insert(user)
//if err != nil {
//log.Println("Error inserting user: ", user)
//}
//review := &Review{Id: newID(), Description: "foobar", User: *user}
//err = reviews.Insert(review)
//if err != nil {
//log.Println("Error inserting review", err)
//}
//book, _ := GetBookByName("english for dummies")
//fmt.Println(book)
//err = UpdateBookById(book.Id, bson.M{"$set": bson.M{"reviews": review}})
//if err != nil {
//log.Println("Updating error", err)
//}
//return nil
//}

func main() {
	//_, err := NewBook("english for dummies", 10, "english", "0xyzueue")
	//if err != nil {
	//fmt.Println("Error inserting new book")
	//return
	//}
	//book, err := GetBookByName("english for dummies")
	//if err != nil {
	//fmt.Println("Nao achou")
	//return
	//}
	//fmt.Println(book.Name)
	//id := book.Id
	//update := bson.M{"$set": bson.M{"name": "modified", "pages": 29}}
	//err = UpdateBookById(id, update)
	//if err != nil {
	//fmt.Println("Error updating: ", err)
	//return
	//}
	//fmt.Println("Updated")
	//fmt.Println("Inseted")
	//fmt.Println(book)
	//err = DeleteBookByName("english for dummies")
	//if err != nil {
	//fmt.Println("Not deleted")
	//}
	//AddReview()
	//t := User{Id: "1", Name: "teste"}
	//json.NewEncoder(os.Stdout).Encode(t)
	//http.HandleFunc("/teste", NewUser)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	r := mux.NewRouter()
	r.HandleFunc("/reviews/{id}", GetReviewByIdHandler).Methods("GET")
	r.HandleFunc("/reviews/{id}", DeleteReviewByIdHandler).Methods("DELETE")
	r.HandleFunc("/books", NewBookHandler).Methods("POST")
	r.HandleFunc("/books/{id}", DeleteBookByIdHandler).Methods("DELETE")
	r.HandleFunc("/books/{id}", GetBookByIdHandler).Methods("GET")
	r.HandleFunc("/books/{id}", UpdateBookByIdHandler).Methods("PUT")
	r.HandleFunc("/books/{name}", GetBookByNameHandler).Methods("GET")
	r.HandleFunc("/books/{name}", GetBookByNameHandler).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		return nil, err
	}
	return session, nil
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
