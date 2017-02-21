package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/qiujinwu/gin-utils/sessions"
	"log"
	"net/http"
	"strings"
	"regexp"
	"encoding/gob"
)

const (
	CurrentUserKey = "current_user"
)

type Config struct {
	// true 黑名单，false 白名单
	Blacklist bool
	// 名单列表
	Items map[string]string

	store sessions.Store
}

// regexp.MatchString(pat, src)
func (config *Config) contain(url string) bool {
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
func (config *Config) need_ignore(url string) bool {
	return config.Blacklist && !config.contain(url) ||
		!config.Blacklist && config.contain(url)
}

func newConfig(blacklist bool) *Config {
	return &Config{
		Blacklist: blacklist,
		Items:     make(map[string]string),
	}
}

//--------------------------------------
var _config *Config = nil
var _handle gin.HandlerFunc = nil

var defaultErrorFunc = func(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"code": http.StatusForbidden,
		"message": "authentication requried",
	})
}

// Options stores configurations for a CSRF middleware.
type Options struct {
	Blacklist bool
	ErrorFunc gin.HandlerFunc
	SessionStore sessions.Store
	User interface{}
}

func AddUrl(url string,regex string) {
	if _config == nil {
		log.Fatal("add url before new filter")
		return
	}
	_config.Items[url] = regex
}

func Login(c *gin.Context,user interface{}) {
	if _config == nil {
		log.Fatal("add url before new filter")
		return
	}
	session_inst, _ := _config.store.Get(c, "session")
	session_inst.Set("user", user)
	session_inst.Save()
}

func Logout(c *gin.Context) {
	if _config == nil {
		log.Fatal("add url before new filter")
		return
	}

	_config.store.Delete(c, "session")
}

func NewFilter(options Options) gin.HandlerFunc {
	if _handle != nil {
		log.Fatal("bind filter more than once")
		return nil
	}

	if options.ErrorFunc == nil{
		options.ErrorFunc = defaultErrorFunc
	}

	if options.SessionStore == nil{
		log.Fatal("store can NOT be empty")
		return nil
	}

	if _config == nil {
		_config = newConfig(options.Blacklist)
		_config.store = options.SessionStore
	}
	gob.Register(options.User)
	_handle = func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)
		if _config.need_ignore(method + c.Request.URL.Path) {
			c.Next()
			return
		}

		session_inst, _ := options.SessionStore.Get(c, "session")
		v := session_inst.Get("user")
		if v != nil {
			c.Set(CurrentUserKey,v)
			c.Next()
		} else {
			options.ErrorFunc(c)
			c.Abort()
			return
		}

	}
	return _handle
}
