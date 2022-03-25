package book

import (
	"errors"
	"fmt"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/domain/author"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrQuantityTooLow              = errors.New("there is not enough items in the stock")
	ErrNameTooShort                = errors.New("name too short")
	ErrInvalidStockCode            = errors.New("invalid stock code")
	ErrInvalidISBN                 = errors.New("ISBN must be 11 digits")
	ErrNumberOfPagesMustBePositive = errors.New("number of pages must be positive")
	ErrPriceMustBePositive         = errors.New("price must be positive")
	ErrQuantityMustBePositive      = errors.New("quantity must be positive")
	ErrAuthorIDMustBePositive      = errors.New("author id must be positive")
)

type Book struct {
	gorm.Model
	Name          string  `gorm:"not null" json:"name"`
	StockCode     string  `gorm:"not null;unique" json:"stockCode"`
	ISBN          int     `gorm:"not null" json:"ISBN"`
	NumberOfPages int     `gorm:"not null" json:"numberOfPages"`
	Price         float64 `gorm:"not null" json:"price"`
	Quantity      int     `gorm:"not null;default:0" json:"quantity"`
	AuthorID      int     `gorm:"not null" json:"authorID"`
	Author        author.Author
}

func (book Book) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("Deleting book: ", book.Name)
	return nil
}

func (book Book) String() string {
	return fmt.Sprintf("{Id: %d, Name: %s, StockCode: %s, ISBN: %d, NumberOfPages: %d, Price: %f, Quantity: %d, Author: %s}",
		book.ID, book.Name, book.StockCode, book.ISBN, book.NumberOfPages, book.Price, book.Quantity, book.Author)
}

// buy decreases stock count. A Book can't be bought if there is not enough stock or deleted already.
func (book *Book) buy(quantityToBuy int) error {
	if book.Quantity < quantityToBuy {
		return ErrQuantityTooLow
	}

	book.Quantity -= quantityToBuy
	return nil
}

func (book Book) validate() error {
	if len(book.Name) <= 0 {
		return ErrNameTooShort
	} else if !strings.HasPrefix(book.StockCode, "SKU-") {
		return ErrInvalidStockCode
	} else if len(string(rune(book.ISBN))) != 11 {
		return ErrInvalidISBN
	} else if book.NumberOfPages <= 0 {
		return ErrNumberOfPagesMustBePositive
	} else if book.Price <= 0 {
		return ErrPriceMustBePositive
	} else if book.Quantity <= 0 {
		return ErrQuantityMustBePositive
	} else if book.AuthorID <= 0 {
		return ErrAuthorIDMustBePositive
	}

	return nil
}
