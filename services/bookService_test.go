package services

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go-kafka-consumer-restAPI-book-library/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type booksRepositoryMock struct{}

var booksRepositoryGetMock func(id string) ([]domain.Book, error)
var booksRepositoryGetAllMock func() ([]domain.Book, error)
var booksRepositoryCreateMock func(book domain.Book) (int64, error)
var booksRepositoryUpdateMock func(book domain.Book) (int64, error)

func (b booksRepositoryMock) getBookById(id string) ([]domain.Book, error) {
	return booksRepositoryGetMock(id)
}

func (b booksRepositoryMock) getAllBooks() ([]domain.Book, error) {
	return booksRepositoryGetAllMock()
}

func (b booksRepositoryMock) createBook(book domain.Book) (int64, error) {
	return booksRepositoryCreateMock(book)
}

func (b booksRepositoryMock) updateBook(book domain.Book) (int64, error) {
	return booksRepositoryUpdateMock(book)
}

func TestIsValidFalseForInvalidData(t *testing.T) {
	var books = []domain.Book{
		{ISBN: 1, Name: "Book", Author: "Author"},
		{ISBN: 2, Name: "Book", Author: ""},
		{ISBN: 3, Name: "", Author: "Author"},
	}

	for _, book := range books {
		valid := isValid(book)
		assert.False(t, valid)
	}
}

func TestIsValidTrue(t *testing.T) {
	var book = domain.Book{Name: "Book", Author: "Author"}
	valid := isValid(book)
	assert.True(t, valid)
}

func TestProcessMessageCreateBookIfValid(t *testing.T) {
	bookStr := "{\"Name\":\"Book\",\"Author\":\"Author\"}"
	var bookSearchResult []domain.Book

	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return bookSearchResult, nil
	}
	booksRepositoryCreateMock = func(book domain.Book) (int64, error) {
		return int64(1), nil
	}

	isbn, err := ProcessMessage([]byte(bookStr))

	assert.Nil(t, err)
	assert.Equal(t, int64(1), isbn)
}

func TestProcessMessageUpdateBookIfValid(t *testing.T) {
	bookStr := "{\"Name\":\"Book\",\"Author\":\"Author\"}"
	bookSearchResult := []domain.Book{
		{
			ISBN:   1,
			Name:   "Book",
			Author: "Author",
		},
	}

	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return bookSearchResult, nil
	}
	booksRepositoryUpdateMock = func(book domain.Book) (int64, error) {
		return int64(1), nil
	}

	isbn, err := ProcessMessage([]byte(bookStr))

	assert.Nil(t, err)
	assert.Equal(t, int64(1), isbn)
}

func TestProcessMessageFailureIfInvalidRecord(t *testing.T) {
	bookStrs := []string{
		"{\"ISBN\": 19, \"Name\":\"Book\",\"Author\":\"Author\"}",
		"{\"Name\":\"Book\",\"Author\":}",
		"{\"Name\":\"\",\"Author\":\"\"}",
		"{\"Name\":\"Book\",\"Author\":\"\"}",
	}

	for _, bookStr := range bookStrs {
		isbn, err := ProcessMessage([]byte(bookStr))

		assert.NotNil(t, err)
		assert.Equal(t, int64(-1), isbn)
	}
}

func TestProcessMessageUpdateBookFailureForDatabaseErrors(t *testing.T) {
	bookStr := "{\"Name\":\"Book\",\"Author\":\"Author\"}"
	bookSearchResult := []domain.Book{
		{
			ISBN:   1,
			Name:   "Book",
			Author: "Author",
		},
	}

	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return bookSearchResult, nil
	}
	booksRepositoryUpdateMock = func(book domain.Book) (int64, error) {
		return int64(1), errors.New("failed updating database")
	}

	isbn, err := ProcessMessage([]byte(bookStr))

	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), isbn)
}

func TestProcessMessageCreateBookFailureForDatabaseErrors(t *testing.T) {
	bookStr := "{\"Name\":\"Book\",\"Author\":\"Author\"}"
	var bookSearchResult []domain.Book

	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return bookSearchResult, nil
	}
	booksRepositoryCreateMock = func(book domain.Book) (int64, error) {
		return int64(1), errors.New("failed updating database")
	}

	isbn, err := ProcessMessage([]byte(bookStr))

	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), isbn)
}

func TestGetAllBooksHandlerSuccess(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	books := []domain.Book{
		{ISBN: 1, Name: "Book1", Author: "Author1"},
		{ISBN: 2, Name: "Book2", Author: "Author2"},
		{ISBN: 3, Name: "Book3", Author: "Author3"},
	}

	booksRepositoryGetAllMock = func() ([]domain.Book, error) {
		return books, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/books", nil)

	GetAllBooksHandler(w, r)

	var book []domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	_ = assert.Equal(t, http.StatusOK, w.Code)

	for i, book := range books {
		i = i + 1
		index := strconv.FormatInt(int64(i), 10)
		_ = assert.Equal(t, int64(i), book.ISBN)
		_ = assert.Equal(t, "Book"+index, book.Name)
		_ = assert.Equal(t, "Author"+index, book.Author)
	}
}

func TestGetAllBooksHandlerFailure(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	booksRepositoryGetAllMock = func() ([]domain.Book, error) {
		return []domain.Book{}, errors.New("error fetching records")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/books", nil)

	GetAllBooksHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBookByIdHandlerSuccess(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{{ISBN: 1, Name: "Book", Author: "Author"}}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	GetBookByIdHandler(w, r)

	var book domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), book.ISBN)
	assert.Equal(t, "Book", book.Name)
	assert.Equal(t, "Author", book.Author)
}

func TestGetBookByIdHandlerFailureDatabaseError(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{}, errors.New("error fetching record")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	GetBookByIdHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBookByIdHandlerFailureNotFound(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	GetBookByIdHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
