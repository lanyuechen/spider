package spider

import (
	"github.com/spider/douban"
	"github.com/spider/dangdang"
)

type Config struct {
	Douban *douban.Config `json:"douban"`
	Dangdang *dangdang.Config `json:"dangdang"`
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
	Cfg = &Config{
		Douban: douban.Cfg,
		Dangdang: dangdang.Cfg,
	}
}

func Conf(cfg *Config) error {
	douban.Cfg = cfg.Douban
	dangdang.Cfg = cfg.Dangdang
	Cfg = cfg
	return nil
}
