package authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"strings"
	"github.com/qiujinwu/gin-utils/sessions"
	"github.com/qiujinwu/gin-utils/utils"
)

type Config struct{
	// true 黑名单，false 白名单
	Blacklist bool
	// 名单列表
	Items utils.Set

	store sessions.Store
}

/**
 不在黑名单或者在白名单的被忽略
*/
func (config *Config)need_ignore(url string)bool{
	return config.Blacklist && !config.Items.Contains(url) ||
		!config.Blacklist && config.Items.Contains(url)
}

func newConfig(blacklist bool) *Config{
	return &Config{
		Blacklist: blacklist,
		Items:utils.NewHashSet(),
	}
}


//--------------------------------------
var _config *Config = nil
var _handle gin.HandlerFunc = nil

func AddUrl(url string){
	if _config == nil{
		log.Fatal("add url before new filter")
		return
	}
	_config.Items.Add(url)
}

func Login(c *gin.Context){
	if _config == nil{
		log.Fatal("add url before new filter")
		return
	}
	session_inst, _ := _config.store.Get(c,"session")
	session_inst.Set("count", 1)
	session_inst.Save()
}

func Logout(c* gin.Context){
	if _config == nil{
		log.Fatal("add url before new filter")
		return
	}

	_config.store.Delete(c,"session")
}


func NewFilter(blacklist bool,store sessions.Store) gin.HandlerFunc {
	if _handle != nil{
		log.Fatal("bind filter more than once")
		return nil
	}
	if _config == nil{
		_config = newConfig(blacklist)
		_config.store = store
	}
	_handle = func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)
		if _config.need_ignore(method + c.Request.RequestURI){
			c.Next()
			return
		}

		session_inst, _ := store.Get(c,"session")
		v := session_inst.Get("count")
		if v != nil {
			c.Next()
		}else{
			c.JSON(http.StatusForbidden, gin.H{
				"message": "StatusForbidden",
			})
			c.Abort()
			return
		}

	}
	return _handle
}
