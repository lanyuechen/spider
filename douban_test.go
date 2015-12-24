package douban

import (
	"testing"
)

func TestSplitSlice(t *testing.T) {
	book, err := GetBookByIsbn("9787505715660")
	if err != nil {
		t.Log(err)
	}
	if book.Isbn != "9787505715660" {
		t.Log(book)
	}

	book, err = GetBookById("1003078")
	if err != nil {
		t.Log(err)
	}
	if book.Id != "1003078" {
		t.Log(book)
	}

	books, err := SearchBook("小王子", "", 0, 2)
	if err != nil {
		t.Log(err)
	}
	if books.Start != 0 || books.Count != 2 || len(books.Books) != books.Count {
		t.Log(books)
	}
}
