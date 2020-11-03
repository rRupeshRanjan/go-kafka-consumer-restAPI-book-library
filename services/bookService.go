package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-kafka-consumer-restAPI-book-library/domain"
	"go-kafka-consumer-restAPI-book-library/repository"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type BookRepositoryInterface interface {
	getBookById(id string) ([]domain.Book, error)
	getAllBooks() ([]domain.Book, error)
	createBook(book domain.Book) (int64, error)
	updateBook(book domain.Book) (int64, error)
}

func (b BookRepository) getBookById(id string) ([]domain.Book, error) {
	return repository.GetBookById(id)
}

func (b BookRepository) getAllBooks() ([]domain.Book, error) {
	return repository.GetAllBooks()
}

func (b BookRepository) createBook(book domain.Book) (int64, error) {
	return repository.CreateBook(book)
}

func (b BookRepository) updateBook(book domain.Book) (int64, error) {
	return repository.UpdateBook(book)
}

type BookRepository struct{}

var booksRepository BookRepositoryInterface
var log, _ = zap.NewProduction()

func ProcessMessage(msg []byte) {
	var book domain.Book
	err := json.Unmarshal(msg, &book)
	if err != nil {
		log.Error("Error while decoding kafka message to book object: " + err.Error())
	} else if isValid(book) {
		id := strconv.FormatInt(book.ISBN, 10)
		bookRow, _ := repository.GetBookById(id)
		if len(bookRow) == 1 {
			updateBook(bookRow[0])
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
	isbn, err := booksRepository.createBook(book)
	if err == nil {
		book.ISBN = isbn
		log.Info("Successfully inserted kafka record to database: " + getString(book))
	}
}

func updateBook(book domain.Book) {
	isbn, err := booksRepository.updateBook(book)
	if err == nil {
		book.ISBN = isbn
		log.Info("Successfully updated book: " + getString(book))
	}
}

func GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	books, err := booksRepository.getBookById(id)

	if err == nil {
		if len(books) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_, _ = fmt.Fprint(w, getString(books[0]))
		}
	} else {
		log.Error("Error while getting book with id " + id + " with error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func GetAllBooksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []domain.Book
	books, err := booksRepository.getAllBooks()

	if err == nil {
		_, _ = fmt.Fprint(w, getString(books))
	} else {
		log.Error("Error while getting all books from db: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getString(input interface{}) string {
	deserializedObject, deserializationErr := json.Marshal(input)
	stringObject := ""

	if deserializationErr != nil {
		log.Error("Error while deserializing data: " + deserializationErr.Error())
	} else {
		stringObject = string(deserializedObject)
	}

	return stringObject
}
