package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	Id          string `bson:"_id"`
	Description string
	User        User
}

type User struct {
	Id   string `bson:"_id"`
	Name string `json:"name"`
}

type Book struct {
	Id       string `bson:"_id"`
	Name     string
	Pages    int
	Language string
	ISBN     string
	Reviews  []Review
}

func (b Book) AddReview(description string, user User) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	//users := session.DB("bookstore").C("users")

	review := &Review{
		Id:          newID(),
		Description: description,
		User:        user,
	}
	reviews := session.DB("bookstore").C("reviews")
	err = reviews.Insert(review)
	if err != nil {
		//log.Println("Error inserting review", err)
		return errors.New("Error inserting review")
	}
	b.Reviews = append(b.Reviews, *review)
	err = UpdateBookById(b.Id, bson.M{"$set": bson.M{"reviews": b.Reviews}})
	if err != nil {
		//	log.Println("Updating error", err)
		return errors.New("Updating error")
	}
	return nil
}

// Format: {"id":"2","name":"joaozinho"}
// {"name":"alco"}
func NewUser(res http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, "Error decoding JSON", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")
	user.Id = newID()
	err = c.Insert(user)
	if err != nil {
		msg := fmt.Sprintf("Error inserting user: %s", user.Name)
		//http.Error(res, msg, http.StatusInternalServerError)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode("{'user':'saved'}")
}

func GetUserByName(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if name == "" {
		http.Error(res, "No name given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")

	var user User
	err = c.Find(bson.M{"name": name}).One(&user)
	if err != nil {
		msg := fmt.Sprintf("User %s not found", name)
		http.Error(res, msg, http.StatusNotFound)
		//	http.NotFound(res, req)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(user)
}

func GetUserById(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "No id given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")

	var user User
	err = c.FindId(id).One(&user)
	if err != nil {
		msg := fmt.Sprintf("User with id %s not found", id)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(user)
}

func DeleteUserByName(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if name == "" {
		http.Error(res, "No name has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")
	if err = c.Remove(bson.M{"name": name}); err != nil {
		msg := fmt.Sprintf("User with name %s not found", name)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.WriteHeader(http.StatusNoContent)
}

func DeleteUserById(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
	}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")
	if err = c.RemoveId(id); err != nil {
		msg := fmt.Sprintf("User with id %s not found", id)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.WriteHeader(http.StatusNoContent)
}

func UpdateUserById(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
	}

	var params map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(res, "Problem decoding JSON", http.StatusUnprocessableEntity)
	}

	//var user User
	//err := json.NewDecoder(req.Body).Decode(&user)
	//if err != nil {
	//http.Error(res, "Problem decoding JSON", http.StatusUnprocessableEntity)
	//}

	//update := bson.M{"$set": bson.M{"name": "modified", "pages": 29}}
	session, err := getSession()
	if err != nil {
		http.Error(res, "Error opening session", http.StatusInternalServerError)
	}
	defer session.Close()
	c := session.DB("bookstore").C("users")

	update := bson.M{"$set": bson.M{params}}
	if err = c.UpdateId(id, params); err != nil {
		msg := fmt.Sprintf("User with id %s not found", id)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.Write([]byte("user updated"))
}

func GetReviewById(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "No id has been given", http.StatusBadRequest)
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
		msg := fmt.Sprintf("Review with id %s has not been found", id)
		http.Error(res, msg, http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(review)
}

func DeleteReviewById(id string) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("reviews")
	err = c.RemoveId(id)
	return err
}

func UpdateReviewById(id string, params bson.M) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("reviews")
	err = c.UpdateId(id, params)
	return err
}

func NewBook(name string, pages int, language string, ISBN string) (*Book, error) {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	book := &Book{
		Id:       newID(),
		Name:     name,
		Pages:    pages,
		Language: language,
		ISBN:     ISBN,
	}
	err = c.Insert(book)
	if err != nil {
		log.Println("Error inserting book: ", book)
	}
	return book, err
}

func DeleteBookByName(name string) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.Remove(bson.M{"name": name})
	return err
}

func DeleteBookById(id string) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.RemoveId(id)
	return err
}

func GetBookByName(name string) (Book, error) {
	var book Book
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.Find(bson.M{"name": name}).One(&book)
	return book, err
}

func GetBookById(id string) (Book, error) {
	var book Book
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.FindId(id).One(&book)
	return book, err
}

func UpdateBookById(id string, params bson.M) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.UpdateId(id, params)
	return err
}

func AddReview() error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	users := session.DB("bookstore").C("users")
	//books := session.DB("bookstore").C("books")
	reviews := session.DB("bookstore").C("reviews")
	user := &User{Id: newID(), Name: "Renato"}
	err = users.Insert(user)
	if err != nil {
		log.Println("Error inserting user: ", user)
	}
	review := &Review{Id: newID(), Description: "foobar", User: *user}
	err = reviews.Insert(review)
	if err != nil {
		log.Println("Error inserting review", err)
	}
	book, _ := GetBookByName("english for dummies")
	fmt.Println(book)
	//book.Reviews = append(book.Reviews, *review)
	err = UpdateBookById(book.Id, bson.M{"$set": bson.M{"reviews": review}})
	if err != nil {
		log.Println("Updating error", err)
	}
	return nil
}

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
	http.HandleFunc("/teste", NewUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
