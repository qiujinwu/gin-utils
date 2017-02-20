package flash

import (
	"github.com/gin-gonic/gin"
	"github.com/qiujinwu/gin-utils/sessions"
	"log"
	"encoding/gob"
)

var (
	DefaultCookieKey  = "_flash"
	default_store sessions.Store = nil
	// handle gin.HandlerFunc = nil
)

func AddFlash(c *gin.Context, vars ...string) {
	session,_ := default_store.Get(c,DefaultCookieKey)
	session.AddFlash(vars[0])
	session.Save()
}

func Flashes(c *gin.Context) []interface{} {
	session,_ := default_store.Get(c,DefaultCookieKey)
	flashes := session.Flashes()
	// session.Save()
	default_store.Delete(c,DefaultCookieKey)
	return flashes
}

func Init(keyPairs ...[]byte) bool{
	if default_store == nil{
		default_store = sessions.NewCookieStore(keyPairs...)
		gob.Register([]interface{}{})
		return true
	}else{
		log.Fatal("bind filter more than once")
		return false
	}
}


//func NewFilter(keyPairs ...[]byte) gin.HandlerFunc {
//	if default_store == nil{
//		default_store = sessions.NewCookieStore(keyPairs...)
//		gob.Register([]interface{}{})
//	}else{
//		log.Fatal("bind filter more than once")
//		return nil
//	}
//
//
//	handle = func(c *gin.Context) {
//		c.Next()
//		session,_ := default_store.Get(c,DefaultCookieKey)
//		session.Save()
//	}
//	return handle
//}
