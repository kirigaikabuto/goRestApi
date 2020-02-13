package main

import (
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
//abstract factory
var (
	PATH string = ""
	PORT string = ""
)

var flags []cli.Flag = []cli.Flag{
	//&cli.StringFlag{
	//	Name:        "data",
	//	Usage:       "Load JSON FILE",
	//	Aliases:     []string{"d"},
	//	Destination: &PATH,
	//},
	//&cli.StringFlag{
	//	Name:        "port",
	//	Usage:       "SET PORT TO RUN",
	//	Aliases:     []string{"p"},
	//	Destination: &PORT,
	//},
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
	config:=bookstore.PostgreConfig{
				User:     "postgres",
				Password: "passanya",
				Port:     "5432",
				Host:     "127.0.0.1",
				Database: "apibooks",
			}
	ps,err := bookstore.NewPostgreBookStore(config)
	if err!=nil{
		log.Fatal(err)
	}
	endpoints:=bookstore.NewEndpointsFactory(bookStore,ps)
	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(endpoints.GetBook( "id"))
	router.Methods("GET").Path("/").HandlerFunc(endpoints.GetBooks())
	router.Methods("POST").Path("/").HandlerFunc(endpoints.CreateBook())
	router.Methods("DELETE").Path("/{id}").HandlerFunc(endpoints.DeleteBook( "id"))
	router.Methods("PUT").Path("/{id}").HandlerFunc(endpoints.UpdateBook( "id"))

	router.Methods("GET").Path("/books/").HandlerFunc(endpoints.GetBooksDb())
	router.Methods("GET").Path("/books/{id}").HandlerFunc(endpoints.GetBookDb("id"))
	router.Methods("PUT").Path("/books/{id}").HandlerFunc(endpoints.UpdateBookDb("id"))
	router.Methods("POST").Path("/books/").HandlerFunc(endpoints.CreateBookDb())
	router.Methods("DELETE").Path("/books/{id}").HandlerFunc(endpoints.DeleteBookDb("id"))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		http.ListenAndServe(PORT, router)
	}()
	log.Println("Server is running on port" + PORT)

	<-done

	log.Print("Server Stopped")
	ExitWithSave(bookStore)
	return nil
}



func ExitWithSave(store bookstore.BookStore) {
	err := store.SaveBooks(PATH)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Data is saved")
}