package main

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Book struct {
	Id       string `bson:"_id"`
	Name     string
	Pages    int
	Language string
	ISBN     string
	Reviews  []Review
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

type Review struct {
	Id          string `bson:"_id"`
	Description string
	User        User
}

type User struct {
	Id   string `bson:"_id"`
	Name string
}

func main() {
	_, err := NewBook("english for dummies", 10, "english", "0xyzueue")
	if err != nil {
		fmt.Println("Error inserting new book")
		return
	}
	book, err := GetBookByName("english for dummies")
	if err != nil {
		fmt.Println("Nao achou")
		return
	}
	fmt.Println(book.Name)
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
	AddReview()
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
