package douban

type Books struct {
	Start int     `json:"start"`
	Count int     `json:"count"`
	Total int     `json:"total"`
	Books []*Book `json:"books"`
}

type Book struct {
	Id          string `json:"id"`
	Isbn        string `json:"isbn13"`
	Title       string `json:"title"`
	OriginTitle string `json:"origin_title"`
	AltTitle    string `json:"alt_title"`
	SubTitle    string `json:"subtitle"` //副标题
	Url         string `json:"url"`      //接口链接
	Alt         string `json:"alt"`      //豆瓣图书链接
	Image       string `json:"image"`    //中图
	Images      struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"images"`
	Author     []string `json:"author"`
	Translator []string `json:"translator"`
	Publisher  string   `json:"publisher"`
	Pubdate    string   `json:"pubdate"`
	Rating     struct {
		Max       int    `json:"max"`
		NumRaters int    `json:"numRaters"`
		Average   string `json:"average"`
		Min       int    `json:"min"`
	} `json:"rating"`
	Tags []struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"tags"`
	Binding string `json:"binding"` //封装，平装/精装等
	Price   string `json:"price"`
	Series  struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	} `json:"series"`
	Pages       string `json:"pages"`
	AuthorIntro string `json:"author_intro"`
	Summary     string `json:"summary"`
	CataLog     string `json:"catalog"`
	EbookUrl    string `json:"ebook_url"`
	EbookPrice  string `json:"ebook_price"`
}

type Reviews struct {
	Start   int       `json:"start"`
	Count   int       `json:"count"`
	Total   int       `json:"total"`
	Reviews []*Review `json:"reviews"`
}

type Review struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Alt    string `json:"alt"`
	Author *User  `json:"author"`
	Book   *Book  `json:"book"`
	Rating struct {
		Max   int    `json:"max"`
		Value string `json:"value"`
		Min   int    `json:"min"`
	} `json:"rating"`
	Votes     int    `json:"votes"`
	Useless   int    `json:"useless"`
	Comments  int    `json:"comments"`
	Summary   string `json:"summary"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

type User struct {
}
