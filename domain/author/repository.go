package author

import (
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-karsl/infrastructure/data"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

type Repository interface {
	Migration() error
	InsertSampleData() error
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Migration() error {
	err := r.db.AutoMigrate(&Author{})
	if err != nil {
		return err
	}

	return nil
}

// InsertSampleData reads data from author.csv and writes them to table author
func (r *GormRepository) InsertSampleData() error {
	lines, err := data.GetCellsFromCSV("author.csv")
	if err != nil {
		return err
	}

	authors, err := linesToAuthors(lines)
	if err != nil {
		return err
	}

	for _, c := range authors {
		r.db.FirstOrCreate(&c, Author{
			Model: gorm.Model{ID: c.ID},
		})
	}

	return nil
}
