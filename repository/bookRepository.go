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
	insertQuery             = "INSERT INTO books (name, author) VALUES (?, ?)"
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

func GetAllBooks() (*sql.Rows, error) {
	return database.Query(getAllQuery)
}

func GetBookById(id string) (*sql.Rows, error) {
	return database.Query(getQuery, id)
}

func InsertBook(book domain.Book) (int64, error) {
	statement, _ := database.Prepare(insertQuery)
	result, insertRecordErr := statement.Exec(book.Name, book.Author)
	if insertRecordErr != nil {
		log.Error("Error while inserting record into books table: " + insertRecordErr.Error())
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
