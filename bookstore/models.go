package bookstore

type BookStore interface {
	SaveBooks(filename string) error
	CreateBook(book *Book) (*Book,error)
	GetBook(id int64) (*Book,error)
	ListBooks()([]*Book,error)
	DeleteBook(id int64) error
	UpdateBook(id int64,book *Book) (*Book,error)
}
type Book struct {
	ID int64 `json:"id" pg:"id,pk"`
	Title string `json:"title,omitempty" pg:"title"`
	Author string `json:"author,omitempty"  pg:"author"`
	Description string `json:"description,omitempty"  pg:"description"`
	NumberOfPages int `json:"number_of_pages"  pg:"numberofpages"`
}
