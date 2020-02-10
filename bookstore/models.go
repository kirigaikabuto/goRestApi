package bookstore

type BookStore interface {
	SaveBooks(filename string) error
	CreateBook(book *Book) (*Book,error)
	GetBook(id string) (*Book,error)
	ListBooks()([]*Book,error)
	DeleteBook(id string) error
	UpdateBook(id string,book *Book) (*Book,error)
}
type Book struct {
	ID string `json:"id"`
	Title string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
	NumberOfPages int `json:"number_of_pages"`
}
