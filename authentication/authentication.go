package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/qiujinwu/gin-utils/sessions"
	"log"
	"net/http"
	"encoding/gob"
	"github.com/qiujinwu/gin-utils/utils"
)

const (
	CurrentUserKey = "current_user"
)



//--------------------------------------
var _config *utils.Config = nil
var _store sessions.Store = nil
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
	session_inst, _ := _store.Get(c, "session")
	session_inst.Set("user", user)
	session_inst.Save()
}

func Logout(c *gin.Context) {
	if _config == nil {
		log.Fatal("add url before new filter")
		return
	}

	_store.Delete(c, "session")
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
		_config = utils.New(options.Blacklist)
	}

	_store = options.SessionStore
	gob.Register(options.User)
	_handle = func(c *gin.Context) {
		if _config.AllowAccess(c) {
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
