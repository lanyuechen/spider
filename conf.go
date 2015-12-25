package spider

import (
	"github.com/spider/amazon"
	"github.com/spider/dangdang"
	"github.com/spider/douban"
	"github.com/spider/jd"
)

type Config struct {
	Douban   *douban.Config   `json:"douban"`
	Dangdang *dangdang.Config `json:"dangdang"`
	Amazon   *amazon.Config   `json:"dangdang"`
	Jd       *jd.Config       `json:"dangdang"`
}

var (
	Cfg *Config
)

func init() {
	douban.Cfg = &douban.Config{
		ApiKey: "",
	}
	dangdang.Cfg = &dangdang.Config{
		CpsId: "",
	}
	amazon.Cfg = &amazon.Config{
		AssociateTag:   "",
		AWSAccessKeyId: "",
		AWSSecretKey:   "",
	}
	jd.Cfg = &jd.Config{
		AppKey:     "",
		AppSecret:  "",
		UnionId:    "",
		UnionAuth:  "",
		UnionWebId: "",
	}
	Cfg = &Config{
		Douban:   douban.Cfg,
		Dangdang: dangdang.Cfg,
	}
}

func Conf(cfg *Config) error {
	douban.Cfg = cfg.Douban
	dangdang.Cfg = cfg.Dangdang
	amazon.Cfg = cfg.Amazon
	jd.Cfg = cfg.Jd
	Cfg = cfg
	return nil
}
