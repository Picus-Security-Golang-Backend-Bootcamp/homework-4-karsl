package author

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Author struct {
	gorm.Model
	Name      string    `gorm:"not null"`
	BirthDate time.Time `gorm:"not null"`
}

func (author Author) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("Deleting author: ", author.Name)
	return nil
}

func (author Author) String() string {
	return fmt.Sprintf("{ID: %d, Name: %s, BirthDate: %s}", author.ID, author.Name, author.BirthDate.Format("02/01/2006"))
}
