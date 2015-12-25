package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/quexer/utee"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	ApiUrl = "http://webservices.amazon.cn/onca/xml"
)

var (
	Cfg *Config
)

type Asin struct {
	Ebook string `json:"ebook"`
	Sbook string `json:"Sbook"`
}

func getAsin(isbn string) (*Asin, error) {
	url := fmt.Sprintf("%s?%s&Signature=%s", ApiUrl, getStringToSign(isbn), sign(isbn))

	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	res := strings.Replace(string(b), "\n", "", -1)

	reg, _ := regexp.Compile(`<ASIN>([^>]*)</ASIN>`)
	matchAsin := reg.FindAllStringSubmatch(res, -1)

	if matchAsin == nil {
		return nil, fmt.Errorf("no match asin isbn: %v", isbn)
	}

	reg, _ = regexp.Compile(`<ProductGroup>([^>]*)</ProductGroup>`)
	matchType := reg.FindAllStringSubmatch(res, -1)

	if matchType == nil {
		return nil, fmt.Errorf("no match type: %v", isbn)
	}

	asin := &Asin{}

	if len(matchAsin) == len(matchType) {
		for i := 0; i < len(matchAsin); i++ {
			if matchType[i][1] == "eBooks" {
				asin.Ebook = matchAsin[i][1]
			} else {
				asin.Sbook = matchAsin[i][1]
			}
		}
	}

	return asin, nil
}

/**
 * 根据isbn获取图书购买链接，现在只支持实体书（tp＝1）
 */
func GetBuyUrlByIsbn(isbn string, tp int) (*BuyUrl, error) {
	url := fmt.Sprintf("%s?%s&Signature=%s", ApiUrl, getStringToSign(isbn), sign(isbn))

	b, err := utee.HttpGet(url)
	if err != nil {
		return nil, err
	}

	res := strings.Replace(string(b), "\n", "", -1)

	reg, _ := regexp.Compile(`<ASIN>([^>]*)</ASIN>`)
	matchAsin := reg.FindStringSubmatch(res)

	if matchAsin == nil {
		return nil, fmt.Errorf("no match asin isbn: %v", isbn)
	}

	reg, _ = regexp.Compile(`<FormattedPrice>[^>^\d]*(\d*\.\d*)</FormattedPrice>`)
	matchPrice := reg.FindStringSubmatch(res)

	if matchPrice == nil {
		return nil, fmt.Errorf("no match price isbn: %v", isbn)
	}

	buyUrl := &BuyUrl{
		Isbn:  isbn,
		Tp:    tp,
		Url:   rebateUrl(matchAsin[1]),
		Price: matchPrice[1],
	}

	return buyUrl, nil
}

func rebateUrl(asin string) string {
	url := fmt.Sprintf(
		"http://www.amazon.cn/gp/product/%v/ref=as_li_qf_br_asin_il_tl?ie=UTF8&camp=536&creative=3200&creativeASIN=%v&linkCode=as2&tag=%v",
		asin,
		asin,
		Cfg.AssociateTag,
	)
	return url
}

func asin2isbn(asin string) (string, error) {
	url := fmt.Sprintf("http://www.amazon.cn/dp/%s", asin)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}
	res := doc.Find("#detail_bullets_id ul").First().Text()

	reg, err := regexp.Compile(`978\d{10}`)
	if err != nil {
		return "", err
	}
	isbn := reg.FindString(res)

	return isbn, nil
}

func getStringToSign(isbn string) string {
	sortedPairs := []string{
		fmt.Sprintf("AWSAccessKeyId=%v&", Cfg.AWSAccessKeyId),
		fmt.Sprintf("AssociateTag=%v&", Cfg.AssociateTag),
		fmt.Sprintf("Condition=%v&", "New"),
		fmt.Sprintf("IdType=%v&", "ISBN"),
		fmt.Sprintf("ItemId=%v&", isbn),
		fmt.Sprintf("Operation=%v&", "ItemLookup"),
		fmt.Sprintf("ResponseGroup=%v&", "ItemAttributes"),
		fmt.Sprintf("SearchIndex=%v&", "Books"),
		fmt.Sprintf("Service=%v&", "AWSECommerceService"),
		fmt.Sprintf("Timestamp=%v&", timeStampFormat(utee.Tick())),
		fmt.Sprintf("Version=%v", "2013-08-01"),
	}

	str := ""
	for _, v := range sortedPairs {
		str += v
	}
	return str
}

func sign(isbn string) string {

	h := hmac.New(sha256.New, []byte(Cfg.AWSSecretKey))
	io.WriteString(h, "GET\n")
	io.WriteString(h, "webservices.amazon.cn\n")
	io.WriteString(h, "/onca/xml\n")
	io.WriteString(h, getStringToSign(isbn))

	e := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	res := e.EncodeToString(h.Sum(nil))

	res = strings.Replace(res, "+", "%2B", -1)
	res = strings.Replace(res, "=", "%3D", -1)

	return string(res)
}

func timeStampFormat(t int64) string {
	dt := time.Unix(t/1000, 0)
	f := fmt.Sprintf(
		"%d-%.2d-%.2dT%.2d:%.2d:%.2d.000Z",
		dt.Year(),
		dt.Month(),
		dt.Day(),
		dt.Hour(),
		dt.Minute(),
		dt.Second(),
	)
	return strings.Replace(f, ":", "%3A", -1)
}
