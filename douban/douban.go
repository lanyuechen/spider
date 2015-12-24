package douban

import (
	"encoding/json"
	"fmt"
	"github.com/quexer/utee"
)

var (
	Cfg *Config
)

/**
 * 通过id获取图书信息
 * id: 豆瓣图书id
 * GET  https://api.douban.com/v2/book/:id
 */
func GetBookById(id string) (*Book, error) {
	return getBook(fmt.Sprintf(`https://api.douban.com/v2/book/%s?apikey=%s`, id, Cfg.ApiKey))
}

/**
 * 通过isbn获取图书信息
 * name: isbn
 * GET  https://api.douban.com/v2/book/isbn/:name
 */
func GetBookByIsbn(isbn string) (*Book, error) {
	fmt.Println("apiKey", Cfg.ApiKey)
	return getBook(fmt.Sprintf(`https://api.douban.com/v2/book/isbn/%s?apikey=%s`, isbn, Cfg.ApiKey))
}

/**
 * 搜索图书
 * GET  https://api.douban.com/v2/book/search
 * q: 查询关键字,传空忽略
 * tag: 查询的标签,传空忽略
 * start: 起始位置,默认0
 * count: 取结果条数,默认20,最大100
 */
func SearchBook(q string, tag string, start int, count int) (*Books, error) {
	url := fmt.Sprintf(
		"https://api.douban.com/v2/book/search?q=%s&tag=%s&start=%d&count=%d&apikey=%s",
		q,
		tag,
		start,
		count,
		Cfg.ApiKey,
	)
	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var books *Books
	err = json.Unmarshal(b, &books)

	return books, err
}

/**
 * 根据id获取图书的书评
 * GET	https://api.douban.com/v2/book/:id/reviews
 * id: 豆瓣图书id
 * start: 起始位置,默认0
 * count: 取结果条数,默认20
 */
func GetBookReviewsById(id string, start int, count int) (*Reviews, error) {
	url := fmt.Sprintf(`https://api.douban.com/v2/book/%s/reviews?start=%d&count=%d&apikey=%s`, id, start, count, Cfg.ApiKey)
	return getBookReviews(url)
}

/**
 * 根据isbn获取图书的书评
 * GET	https://api.douban.com/v2/book/isbn/:name/reviews
 * name: isbn
 * start: 起始位置,默认0
 * count: 取结果条数,默认20
 */
func GetBookReviewsByIsbn(isbn string, start int, count int) (*Reviews, error) {
	url := fmt.Sprintf(`https://api.douban.com/v2/book/isbn/%s/reviews?start=%d&count=%d&apikey=%s`, isbn, start, count, Cfg.ApiKey)
	return getBookReviews(url)
}

/**
 * 获取图书信息
 * url: 连接地址
 */
func getBook(url string) (*Book, error) {
	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var book *Book
	err = json.Unmarshal(b, &book)

	return book, err
}

/**
 * 获取图书书评列表
 * url: 连接地址
 */
func getBookReviews(url string) (*Reviews, error) {
	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var reviews *Reviews
	err = json.Unmarshal(b, &reviews)

	return reviews, err
}
