package dangdang

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/quexer/utee"
	"net/url"
	"regexp"
	"strings"
)

var (
	Cfg *Config
)

/**
 * 根据isbn获取图书购买链接
 * tp＝0：电子书，tp＝1：实体书
 */
func GetBuyUrlByIsbn(isbn string, tp int) (*BuyUrl, error) {
	medium := "22" //电子书
	if tp == 1 {
		medium = "01" //实体书
	}

	url := fmt.Sprintf("http://search.dangdang.com/?medium=%s&key4=%s", medium, isbn)

	//获取文档（编码为gbk，需要先转码）
	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	html := mahonia.NewDecoder("gb18030").ConvertString(string(b))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	buyUrl := &BuyUrl{
		Isbn: isbn,
		Tp:   tp,
	}

	sel := doc.Find(".search_left .line1").First()

	price := sel.Find(".search_now_price").First().Text()
	match, _ := regexp.Compile(`[0-9.]+`)
	buyUrl.Price = match.FindString(price)

	buyLink, ok := sel.Find("a[name=itemlist-picture]").First().Attr("href")
	if !ok {
		return nil, fmt.Errorf("no buyLink")
	}
	buyUrl.Url = cps(buyLink)

	return buyUrl, nil
}

/**
 * 返利链接
 */
func cps(link string) string {
	return fmt.Sprintf(
		"http://union.dangdang.com/transfer.php?from=%s&ad_type=10&sys_id=1&backurl=%s",
		Cfg.CpsId,
		url.QueryEscape(link),
	)
}
