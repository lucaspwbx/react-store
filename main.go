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

func DeleteBook(name string) error {
	session, err := getSession()
	if err != nil {
		log.Fatalln("Error opening session")
	}
	defer session.Close()
	c := session.DB("bookstore").C("books")
	err = c.Remove(bson.M{"name": name})
	return err
}

func GetBook(name string) (Book, error) {
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

type Review struct {
	Id          string `bson:"_id"`
	Description string
	User        User
	Book        Book
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
	book, err := GetBook("english for dummies")
	if err != nil {
		fmt.Println("Nao achou")
		return
	}
	fmt.Println(book.Name)
	//fmt.Println("Inseted")
	//fmt.Println(book)
	err = DeleteBook("english for dummies")
	if err != nil {
		fmt.Println("Not deleted")
	}
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
