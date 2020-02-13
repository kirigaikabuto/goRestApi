package bookstore

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type PostgreConfig struct {
	User string
	Password string
	Port string
	Host string
	Database string
}

type postgreStore struct {
	db * pg.DB
}
func NewPostgreBookStore(config PostgreConfig) (*postgreStore, error) {
	db:=pg.Connect(&pg.Options{
		Addr:config.Host+":"+config.Port,
		User:config.User,
		Password:config.Password,
		Database:config.Database,
	})
	err:=createSchema(db)
	if err!=nil{
		return nil, err
	}
	return &postgreStore{db:db},nil
}
func createSchema(db *pg.DB) error{
	for _, model := range []interface{}{(*Book)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:true,
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (ps *postgreStore) CreateBook(book *Book)  (*Book,error) {
	return book,ps.db.Insert(book)
}
func (ps *postgreStore) GetBooks() ([]*Book,error){
	var books[]*Book
	err := ps.db.Model(&books).Select()
	return books,err
}
func (ps *postgreStore) GetBook(id int64) (*Book,error){
	book_find:=&Book{ID:id}
	err:= ps.db.Select(book_find)
	if err!=nil{
		return nil,errors.New("No book by id")
	}
	return book_find,nil
}
func (ps *postgreStore) UpdateBook(book *Book) (*Book,error){
	err:=ps.db.Update(book)
	if err!=nil{
		return nil,errors.New("Some error when updated")
	}
	return book,nil
}
func (ps *postgreStore) DeleteBook(book *Book) error{

	err:=ps.db.Delete(book)
	if err!=nil{
		return errors.New("Book not deleted")
	}
	return nil
}