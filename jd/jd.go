package jd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2/bson"
	"io"
	"strings"
	"time"
)

const (
	JD_API_HOST  = "http://gw.api.jd.com/routerjson"
	JD_API_UNION = "jingdong.service.promotion.getcode"
)

var (
	Cfg *Config
)

func GetBuyUrlByIsbn(isbn string, tp int) (*BuyUrl, error) {

	url := fmt.Sprintf("http://search.jd.com/bookadvsearch?isbn=%s&enc=utf-8", isbn)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	s := doc.Find("#plist ul:first-child li:first-child .p-price strong").First()

	skuid, ok := s.Attr("class")
	if !ok {
		return nil, fmt.Errorf("skuid err")
	}
	skuid = strings.Trim(skuid, "J_")
	price, ok := s.Attr("data-price")
	if !ok {
		return nil, fmt.Errorf("price err")
	}

	url, _ = Rebate(fmt.Sprintf("http://item.m.jd.com/product/%s.html", skuid))

	buyUrl := &BuyUrl{
		Isbn:  isbn,
		Tp:    tp,
		Price: price,
		Url:   url,
	}

	return buyUrl, nil
}

func Rebate(url string) (string, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	unionUrl := jdUnion(JD_API_UNION, timestamp, url, Cfg.UnionAuth)
	b, err := utee.HttpGet(unionUrl)
	if err != nil {
		return url, err
	}
	jdUnion := bson.M{}
	err = json.Unmarshal(b, &jdUnion)
	if err != nil {
		return url, err
	}
	if jdUnion["error_response"] != nil {
		jdErr := jdUnion["error_response"].(map[string]interface{})
		if jdErr["zh_desc"] != nil {
			return url, fmt.Errorf(jdErr["zh_desc"].(string))
		}
	}
	if jdUnion["jingdong_service_promotion_getcode_responce"] == nil {
		return url, fmt.Errorf("no union responce")
	}
	jdUnionResponce := jdUnion["jingdong_service_promotion_getcode_responce"].(map[string]interface{})
	if jdUnionResponce["queryjs_result"] == nil {
		return url, fmt.Errorf("no url result")
	}
	jdUnionResult := jdUnionResponce["queryjs_result"].(string)

	res := bson.M{}
	err = json.Unmarshal([]byte(jdUnionResult), &res)
	if err != nil {
		return url, err
	}
	if res["resultCode"] == nil {
		return url, fmt.Errorf("bad url")
	}
	if res["resultCode"].(string) != "0" {
		return url, fmt.Errorf("err code %v, msg %v", res["resultCode"], res["resultMessage"])
	}
	if res["url"] == nil {
		return url, fmt.Errorf("err url")
	}
	url = res["url"].(string)

	return url, nil
}

func jdUnion(method string, timestamp string, url string, access string) string {
	buy_param_json := fmt.Sprintf(
		`360buy_param_json={"channel":"WL","materialId":"%v","promotionType":7,"unionId":%d,"webId":"%v"}`,
		url,
		Cfg.UnionId,
		Cfg.UnionWebId,
	)
	return jdUnionUrl(method, timestamp, buy_param_json, access)
}

func jdUnionUrl(method string, timestamp string, buy_param_json string, access string) string {
	sign := jdUnionSign(method, timestamp, buy_param_json, access)
	timestamp = strings.Replace(timestamp, " ", "%20", -1)
	buy_param_json = strings.Replace(buy_param_json, `"`, "%22", -1)
	url := fmt.Sprintf(
		"%v?%v&access_token=%v&app_key=%v&method=%v&timestamp=%v&v=%v&sign=%v",
		JD_API_HOST,
		buy_param_json,
		access,
		Cfg.AppKey,
		method,
		timestamp,
		"2.0",
		sign,
	)
	return url
}

func jdUnionSign(method string, timestamp string, buy_param_json string, access string) string {
	buy_param_json = strings.Replace(buy_param_json, "=", "", -1)
	sign := fmt.Sprintf(
		"%v%vaccess_token%vapp_key%vmethod%vtimestamp%vv%v%v",
		Cfg.AppSecret,
		buy_param_json,
		access,
		Cfg.AppKey,
		method,
		timestamp,
		"2.0",
		Cfg.AppSecret,
	)
	h := md5.New()
	io.WriteString(h, sign)
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}
