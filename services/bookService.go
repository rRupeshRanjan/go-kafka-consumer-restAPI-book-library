package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kafka-consumer/domain"
	"kafka-consumer/repository"
	"net/http"
)

var log, _ = zap.NewProduction()

func ProcessMessage(msg []byte) {
	var book domain.Book
	err := json.Unmarshal(msg, &book)
	if err != nil {
		log.Error("Error while decoding kafka message to book object: " + err.Error())
	} else if isValid(book) {
		row, _ := repository.GetBookById(string(book.ISBN))
		if row.Next() {
			updateBook(book)
		} else {
			createBook(book)
		}
	} else {
		log.Error("Invalid record received (wont be inserted into db): " + string(msg))
	}
}

func isValid(book domain.Book) bool {
	return book.ISBN == 0 && len(book.Name) > 0 && len(book.Author) > 0
}

func createBook(book domain.Book) {
	isbn, err := repository.InsertBook(book)
	if err == nil {
		book.ISBN = isbn
		log.Info("Successfully inserted kafka record to database: " + getString(book))
	}
}

func updateBook(book domain.Book) {
	isbn, err := repository.UpdateBook(book)
	if err == nil {
		book.ISBN = isbn
		log.Info("Successfully updated book: " + getString(book))
	}
}

func GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	rows, err := repository.GetBookById(id)

	if err == nil {
		var book domain.Book
		var ISBN int64
		var name string
		var author string

		for rows.Next() {
			_ = rows.Scan(&ISBN, &name, &author)
			book = domain.Book{ISBN: ISBN, Name : name, Author: author}
		}
		_, _ = fmt.Fprint(w, getString(book))
	} else {
		log.Error("Error while getting book with id " + id + " with error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func GetAllBooksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []domain.Book
	rows, err := repository.GetAllBooks()

	if err == nil {
		var ISBN int64
		var name string
		var author string

		for rows.Next() {
			_ = rows.Scan(&ISBN, &name, &author)
			books = append(books, domain.Book{ISBN: ISBN, Name : name, Author: author})
		}
		_, _ = fmt.Fprint(w, getString(books))
	} else {
		log.Error("Error while getting all books from db: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getString(input interface{}) string {
	deserializedObject, deserializationErr := json.Marshal(input)

	if deserializationErr != nil {
		log.Error("Error while deserializing data: " + deserializationErr.Error())
		return ""
	} else {
		return string(deserializedObject)
	}
}
