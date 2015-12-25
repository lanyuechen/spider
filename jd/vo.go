package jd

type Config struct {
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	UnionId    string `json:"unionId"`
	UnionAuth  string `json:"unionAuth"`
	UnionWebId string `json:"unionWebId"`
}

type BuyUrl struct {
	Isbn  string
	Tp    int
	Url   string
	Price string
}
