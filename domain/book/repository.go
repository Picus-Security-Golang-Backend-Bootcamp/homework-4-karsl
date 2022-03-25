package book

import (
	"errors"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/infrastructure/data"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

var (
	ErrBookNotFound = errors.New("no such book found")
)

type GormRepository struct {
	db *gorm.DB
}

type Repository interface {
	Search(term string) []Book
	List() []Book
	Buy(bookId, quantity int) error
	Delete(bookId int) error
	Create(book Book) error
	Migration() error
	InsertSampleData() error
}

func New(db *gorm.DB) Repository {
	return GormRepository{db: db}
}

func (r GormRepository) Migration() error {
	err := r.db.AutoMigrate(&Book{})
	if err != nil {
		return err
	}

	return nil
}

// InsertSampleData reads data from book.csv and writes them to table book
func (r GormRepository) InsertSampleData() error {
	lines, err := data.GetCellsFromCSV("book.csv")
	if err != nil {
		return err
	}

	books, err := linesToBook(lines)
	if err != nil {
		return err
	}

	for _, c := range books {
		r.db.FirstOrCreate(&c, Book{Name: c.Name})
	}

	return nil
}

// Search returns all books that match with term in the bookshelf, by checking if it matches with given term by
// considering Name, Author.Name, StockCode fields of Book.
func (r GormRepository) Search(term string) []Book {
	var foundBooks []Book
	term = strings.ToLower(term)
	r.db.Preload(clause.Associations).Joins("JOIN Author on Author.id = Book.author_id").
		Where("lower(Book.name) LIKE ?", "%"+term+"%").
		Or("lower(Author.name) LIKE ?", "%"+term+"%").
		Or("lower(Book.stock_code) LIKE ?", "%"+term+"%").
		Find(&foundBooks)

	return foundBooks
}

// List returns all books in the bookshelf.
func (r GormRepository) List() []Book {
	var books []Book
	r.db.Find(&books)
	return books
}

// Buy find the book with given id and buys it.
func (r GormRepository) Buy(bookId, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}

	foundBook, err := r.findBookById(bookId)
	if err != nil {
		return ErrBookNotFound
	}

	err = foundBook.buy(quantity)
	if err != nil {
		return err
	}

	r.db.Save(foundBook)

	return nil
}

// Delete deletes the given book
func (r GormRepository) Delete(bookId int) error {
	result := r.db.Delete(&Book{}, bookId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}

		return result.Error
	}

	return nil
}

// Create book if it's valid
func (r GormRepository) Create(book Book) error {
	err := book.validate()
	if err != nil {
		return err
	}

	result := r.db.Create(&book)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// findBookById returns the first book with given id in books.
func (r GormRepository) findBookById(bookId int) (Book, error) {
	var book Book
	result := r.db.First(&book, bookId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Book{}, ErrBookNotFound
	}

	return book, nil
}
