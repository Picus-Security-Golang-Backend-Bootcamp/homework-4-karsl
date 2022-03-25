package book

import (
	"strconv"
)

// lineToBook parses single line to Book
func lineToBook(line []string) (Book, error) {
	name := line[0]
	stockCode := line[1]
	authorId, err := strconv.Atoi(line[2])
	if err != nil {
		return Book{}, err
	}
	isbn, err := strconv.Atoi(line[3])
	if err != nil {
		return Book{}, err
	}
	numberOfPages, err := strconv.Atoi(line[4])
	if err != nil {
		return Book{}, err
	}
	price, err := strconv.ParseFloat(line[5], 64)
	if err != nil {
		return Book{}, err
	}
	quantity, err := strconv.Atoi(line[6])
	if err != nil {
		return Book{}, err
	}

	data := Book{
		Name:          name,
		StockCode:     stockCode,
		ISBN:          isbn,
		NumberOfPages: numberOfPages,
		Price:         price,
		Quantity:      quantity,
		AuthorID:      authorId,
	}

	return data, nil
}

// linesToBook parses lines to Books
func linesToBook(lines [][]string) ([]Book, error) {
	var result []Book
	for _, line := range lines[1:] {
		data, err := lineToBook(line)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}
