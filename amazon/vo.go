package amazon

type Config struct {
	AssociateTag   string `json:"associateTag"`
	AWSAccessKeyId string `json:"awsAccessKeyId"`
	AWSSecretKey   string `json:"awsSecretKey"`
}

type BuyUrl struct {
	Isbn  string
	Tp    int
	Url   string
	Price string
}
