package api

import (
	"encoding/json"
	"errors"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/domain/author"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/domain/book"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/infrastructure/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	bookRepository   book.Repository
	authorRepository author.Repository
)

func init() {
	db := data.NewPostgresDB()

	authorRepository = author.New(db)
	err := authorRepository.Migration()
	if err != nil {
		panic(err.Error())
	}
	err = authorRepository.InsertSampleData()
	if err != nil {
		panic(err)
	}

	bookRepository = book.New(db)
	err = bookRepository.Migration()
	if err != nil {
		panic(err.Error())
	}
	err = bookRepository.InsertSampleData()
	if err != nil {
		panic(err)
	}
}

func getHandler() *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/books").Subrouter()
	s.HandleFunc("", createBook).Methods("POST")
	s.HandleFunc("/{id:[0-9]+}", deleteBook).Methods("DELETE")
	s.HandleFunc("", searchBook).Queries("term", "{term:[a-zA-Z]+}").Methods("GET")
	s.HandleFunc("", listBooks).Methods("GET")
	s.HandleFunc("/buy/{id:[0-9]+}", buyBook).Queries("amount", "{amount:[0-9]+}").Methods("POST")

	return r
}

func buyBook(writer http.ResponseWriter, request *http.Request) {
	// These are guaranteed to be valid numbers with regex
	amount, _ := strconv.Atoi(mux.Vars(request)["amount"])
	bookId, _ := strconv.Atoi(mux.Vars(request)["id"])

	err := bookRepository.Buy(bookId, amount)
	if err != nil {
		switch {
		case errors.Is(err, book.ErrBookNotFound):
			writer.WriteHeader(http.StatusNotFound)
		case errors.Is(err, book.ErrQuantityTooLow):
			writer.WriteHeader(http.StatusUnprocessableEntity)
		}
		writer.Write([]byte(err.Error()))
	}

	writer.WriteHeader(http.StatusOK)
}

func searchBook(writer http.ResponseWriter, request *http.Request) {
	term := mux.Vars(request)["term"]
	foundBooks := bookRepository.Search(term)

	switch numberOfBooksFound := len(foundBooks); {
	case numberOfBooksFound > 0:
		writer.WriteHeader(http.StatusOK)
	case numberOfBooksFound <= 0:
		writer.WriteHeader(http.StatusNotFound)
	}

	body, _ := json.Marshal(struct {
		Books []book.Book `json:"books"`
	}{
		Books: foundBooks,
	})

	writer.Write(body)
}

func listBooks(writer http.ResponseWriter, request *http.Request) {
	booksInTheBookShelf := bookRepository.List()

	writer.WriteHeader(http.StatusOK)

	body, _ := json.Marshal(struct {
		Books []book.Book `json:"books"`
	}{
		Books: booksInTheBookShelf,
	})

	writer.Write(body)
}

func deleteBook(writer http.ResponseWriter, request *http.Request) {
	// ID guaranteed to be number with regex.
	bookId, _ := strconv.Atoi(mux.Vars(request)["id"])

	err := bookRepository.Delete(bookId)
	if err != nil {
		if errors.Is(err, book.ErrBookNotFound) {
			writer.WriteHeader(http.StatusNotFound)
		}

		writer.WriteHeader(http.StatusUnprocessableEntity)
	}

	writer.WriteHeader(http.StatusOK)
}

func createBook(writer http.ResponseWriter, request *http.Request) {
	var b book.Book

	err := DecodeJSONBody(writer, request, &b)
	if err != nil {
		var mr *MalformedRequest
		if errors.As(err, &mr) {
			http.Error(writer, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = bookRepository.Create(b)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
	}

	writer.WriteHeader(http.StatusOK)
}
