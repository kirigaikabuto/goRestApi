package bookstore

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func NewBookStore(filename string) (BookStore,error){
	file,err := os.Open(filename)
	if err!=nil{
		return nil,err
	}
	buffer := bufio.NewReader(file)

	data,err := ioutil.ReadAll(buffer)
	if err!=nil{
		return nil,errors.New("Ваш json пуст")
	}
	var books []*Book
	if err:=json.Unmarshal(data,&books);err!=nil{
		return nil,err
	}
	defer file.Close()
	return &bookStoreClass{books},nil
}
type bookStoreClass struct {
	books[]*Book
}
func (bsc *bookStoreClass) SaveBooks(filename string) error{
	file,err:=os.OpenFile(filename,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0604)
	if err!=nil{
		return err
	}
	data,err:=json.Marshal(bsc.books)
	if err!=nil{
		return err
	}
	_,err = file.Write(data)
	if err!=nil{
		return err
	}
	return nil
}
func (bsc *bookStoreClass)  CreateBook(book *Book) (*Book,error){
	bsc.books = append(bsc.books,book)
	return book,nil
}
func (bsc *bookStoreClass) GetBook(id int64) (*Book,error) {
	for _,v := range bsc.books{
		if v.ID == id{
			return v,nil
		}
	}
	return nil,errors.New("Not found")
}
func (bsc *bookStoreClass) ListBooks() ([]*Book,error){
	books:=bsc.books
	return books,nil
}
func (bsc *bookStoreClass) DeleteBook(id int64) error{
	for i,v := range bsc.books{
		if v.ID == id{
			bsc.books = append(bsc.books[:i],bsc.books[i+1:]...)
			return nil
		}
	}
	return errors.New("Not deleted")
}
func (bsc *bookStoreClass) UpdateBook(id int64,book *Book) (*Book,error){
	for i,v := range bsc.books{
		if v.ID == id{
			bsc.books[i]=book
			return bsc.books[i],nil
		}
	}
	return nil,errors.New("Not updated")
}