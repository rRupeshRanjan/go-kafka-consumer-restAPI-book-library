package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-kafka-consumer-restAPI-book-library/domain"
	"go.uber.org/zap"
)

var database *sql.DB
var log, _ = zap.NewProduction()

const (
	getQuery                = "SELECT * FROM books WHERE ISBN=?"
	getAllQuery             = "SELECT ISBN, name, author FROM books"
	createQuery             = "INSERT INTO books (name, author) VALUES (?, ?)"
	updateQuery             = "UPDATE books SET name=?, author=? where ISBN=?"
	initializeDatabaseQuery = `CREATE TABLE IF NOT EXISTS books (
		ISBN INTEGER PRIMARY KEY,
		name TEXT,
		author TEXT);`
)

func InitBooksDb() error {
	database, _ = sql.Open("sqlite3", "books.sql")
	_, err := database.Exec(initializeDatabaseQuery)
	return err
}

func GetAllBooks() ([]domain.Book, error) {
	rows, err := database.Query(getAllQuery)

	var books []domain.Book
	if err == nil {
		var book domain.Book
		for rows.Next() {
			_ = rows.Scan(&book.ISBN, &book.Name, &book.Author)
			books = append(books, book)
		}
	}

	return books, err
}

func GetBookById(id string) ([]domain.Book, error) {
	rows, err := database.Query(getQuery, id)

	var books []domain.Book
	if err == nil {
		var book domain.Book
		for rows.Next() {
			_ = rows.Scan(&book.ISBN, &book.Name, &book.Author)
			books = append(books, book)
		}
	}

	return books, err
}

func CreateBook(book domain.Book) (int64, error) {
	statement, _ := database.Prepare(createQuery)
	result, createRecordErr := statement.Exec(book.Name, book.Author)
	if createRecordErr != nil {
		log.Error("Error while creating record into books table: " + createRecordErr.Error())
	}
	return result.LastInsertId()
}

func UpdateBook(book domain.Book) (int64, error) {
	statement, _ := database.Prepare(updateQuery)
	result, updateRecordErr := statement.Exec(book.Name, book.Author, book.ISBN)
	if updateRecordErr != nil {
		log.Error("Error while updating record in books table: " + updateRecordErr.Error())
	}
	return result.LastInsertId()
}
