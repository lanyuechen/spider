package dangdang

type Config struct {
	CpsId string `json:"cpsId"`
}

type BuyUrl struct {
	Isbn  string
	Tp    int
	Url   string
	Price string
}
