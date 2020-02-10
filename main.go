package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Lesson06.02.2020/bookstore"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

var (
	PATH string = ""
	PORT string = ""
)

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:        "data",
		Usage:       "Load JSON FILE",
		Aliases:     []string{"d"},
		Destination: &PATH,
	},
	&cli.StringFlag{
		Name:        "port",
		Usage:       "SET PORT TO RUN",
		Aliases:     []string{"p"},
		Destination: &PORT,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "test RESTFULL api"
	app.Usage = "simple CRUD"
	app.Version = "1.0.0"
	app.Flags = flags
	app.Action = runRestApi
	fmt.Println(app.Run(os.Args))
}
func check_file() error {
	file, _ := ioutil.ReadFile(PATH)
	if len(file) == 0 {
		return errors.New("база данных пуста или не существует")
	}
	return nil
}
func runRestApi(*cli.Context) error {
	if PATH == "" {
		PATH = "books.json"
	}
	if PORT == "" {
		PORT = ":8080"
	}
	if errorfile := check_file(); errorfile != nil {
		log.Fatal(errorfile)
	}
	bookStore, err := bookstore.NewBookStore(PATH)
	if err != nil {
		log.Fatal(err.Error())
	}
	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(GetBook(bookStore, "id"))
	router.Methods("POST").Path("/").HandlerFunc(CreateBook(bookStore))
	router.Methods("GET").Path("/").HandlerFunc(GetBooks(bookStore))
	router.Methods("DELETE").Path("/{id}").HandlerFunc(DeleteBook(bookStore, "id"))
	router.Methods("PUT").Path("/{id}").HandlerFunc(UpdateBook(bookStore, "id"))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(PORT, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server is running on port" + PORT)

	<-done

	log.Print("Server Stopped")
	ExitWithSave(bookStore)
	return nil
}
func SendResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
func GetBook(store bookstore.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := vars[idParam]
		if !err {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, error := store.GetBook(id)
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

func CreateBook(store bookstore.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		book := &bookstore.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			SendResponse(w, http.StatusBadRequest, "Error"+err.Error())
			return
		}
		result, err := store.CreateBook(book)
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
func ExitWithSave(store bookstore.BookStore) {
	err := store.SaveBooks(PATH)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Data is saved")
}
func GetBooks(store bookstore.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := store.ListBooks()
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
func DeleteBook(store bookstore.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		err := store.DeleteBook(id)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		SendResponse(w, http.StatusOK, "Element was deleted")
	}
}
func UpdateBook(store bookstore.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, error := vars[idParam]
		if !error {
			SendResponse(w, http.StatusBadRequest, "не был передан аргумент")
			return
		}
		book, err := store.GetBook(id)
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
		updated_book, err := store.UpdateBook(id, book)
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
