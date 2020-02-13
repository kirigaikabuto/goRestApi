package bookstore

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Endpoints interface {
	GetBook(idParam string) func(w http.ResponseWriter, r *http.Request)
	CreateBook() func(w http.ResponseWriter, r *http.Request)
	GetBooks() func(w http.ResponseWriter, r *http.Request)
	DeleteBook(idParam string) func(w http.ResponseWriter, r *http.Request)
	UpdateBook(idParam string) func(w http.ResponseWriter, r *http.Request)
	CreateBookDb() func(w http.ResponseWriter, r *http.Request)
	GetBooksDb() func(w http.ResponseWriter, r *http.Request)
	UpdateBookDb(idParam string) func(w http.ResponseWriter, r *http.Request)
	GetBookDb(idParam string) func(w http.ResponseWriter, r *http.Request)
	DeleteBookDb(idParam string) func(w http.ResponseWriter, r *http.Request)
}
type endpointsFactory struct {
	bookstore BookStore
	ps *postgreStore
}

func NewEndpointsFactory(bookstore BookStore, ps *postgreStore) Endpoints {
	return &endpointsFactory{bookstore: bookstore, ps: ps}
}
func SendResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
func (ef *endpointsFactory) GetBookDb(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)
		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !err {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, error := ef.ps.GetBook(newint)
		if error != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+error.Error())
			return
		}
		newdata, newerr := json.Marshal(book)
		if newerr != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+newerr.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(newdata))
	}
}

func (ef *endpointsFactory) CreateBookDb() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		book := &Book{}
		if err := json.Unmarshal(data, book); err != nil {
			SendResponse(w, http.StatusBadRequest, "Error"+err.Error())
			return
		}
		result, err := ef.ps.CreateBook(book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
		}
		response, err := json.Marshal(result)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusCreated, string(response))
	}
}
func (ef *endpointsFactory) GetBooksDb() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := ef.ps.GetBooks()
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		newbooks, error := json.Marshal(books)
		if error != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(newbooks))
	}
}
func (ef *endpointsFactory) UpdateBookDb(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)
		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, err := ef.ps.GetBook(newint)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		if err := json.Unmarshal(data, &book); err != nil {
			SendResponse(w, http.StatusBadRequest, "Error"+err.Error())
			return
		}
		updated_book, err := ef.ps.UpdateBook(book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		result, err := json.Marshal(updated_book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(result))
	}
}
func (ef *endpointsFactory) DeleteBookDb(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)

		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, err := ef.ps.GetBook(newint)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		err = ef.ps.DeleteBook(book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, "Element was deleted")
	}
}
func (ef *endpointsFactory) GetBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)
		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !err {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, error := ef.bookstore.GetBook(newint)
		if error != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+error.Error())
			return
		}
		newdata, newerr := json.Marshal(book)
		if newerr != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+newerr.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(newdata))
	}
}

func (ef *endpointsFactory) CreateBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		book := &Book{}
		if err := json.Unmarshal(data, book); err != nil {
			SendResponse(w, http.StatusBadRequest, "Error"+err.Error())
			return
		}
		result, err := ef.bookstore.CreateBook(book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
		}
		response, err := json.Marshal(result)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusCreated, string(response))
	}
}

func (ef *endpointsFactory) GetBooks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := ef.bookstore.ListBooks()
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		newbooks, error := json.Marshal(books)
		if error != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(newbooks))
	}
}
func (ef *endpointsFactory) DeleteBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)
		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		err := ef.bookstore.DeleteBook(newint)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, "Element was deleted")
	}
}
func (ef *endpointsFactory) UpdateBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		newint, newerr := strconv.ParseInt(id, 10, 0)
		if newerr != nil {
			SendResponse(w, http.StatusBadRequest, "Аргумент должен быть числом")
			return
		}
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, err := ef.bookstore.GetBook(newint)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		if err := json.Unmarshal(data, &book); err != nil {
			SendResponse(w, http.StatusBadRequest, "Error"+err.Error())
			return
		}
		updated_book, err := ef.bookstore.UpdateBook(newint, book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		result, err := json.Marshal(updated_book)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, string(result))
	}
}
