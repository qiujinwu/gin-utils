package utils

import (
	"regexp"
	"github.com/gin-gonic/gin"
	"strings"
)

type Config struct {
	// true 黑名单，false 白名单
	Blacklist bool
	// 名单列表
	Items map[string]string
}

// regexp.MatchString(pat, src)
func (config *Config) Contain(url string) bool {
	for k, v := range config.Items {
		if v != ""{
			if match,err := regexp.MatchString(v, url);match && err == nil{
				return true
			}
		}else{
			if k == url{
				return true
			}
		}
	}
	return false
}

/**
不在黑名单或者在白名单的被忽略
*/
func (config *Config) NeedIgnore(url string) bool {
	return config.Blacklist && !config.Contain(url) ||
		!config.Blacklist && config.Contain(url)
}

func (config *Config) AllowAccess(c * gin.Context) bool{
	method := strings.ToUpper(c.Request.Method)
	return config.NeedIgnore(method + c.Request.URL.Path)
}

func New(blacklist bool) *Config {
	return &Config{
		Blacklist: blacklist,
		Items:     make(map[string]string),
	}
}
