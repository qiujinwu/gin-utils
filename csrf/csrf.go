package csrf

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"github.com/dchest/uniuri"
	"github.com/qiujinwu/gin-utils/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/qiujinwu/gin-utils/utils"
)

const (
	csrfSecret = "csrfSecret"
	csrfSalt   = "csrfSalt"
	csrfToken  = "csrfToken"
)

var (
	DefaultCookieKey  = "_crsf"
	default_store sessions.Store = nil
 	_config *utils.Config = nil
)

var defaultIgnoreMethods = []string{"GET", "HEAD", "OPTIONS"}

var defaultErrorFunc = func(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"code": http.StatusForbidden,
		"message": "CSRF token mismatch",
	})
}

var defaultTokenGetter = func(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}

	return ""
}

// Options stores configurations for a CSRF middleware.
type Options struct {
	Secret        string
	IgnoreMethods []string
	ErrorFunc     gin.HandlerFunc
	TokenGetter   func(c *gin.Context) string
	Blacklist bool
}

func tokenize(secret, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt+"-"+secret)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func inArray(arr []string, value string) bool {
	inarr := false

	for _, v := range arr {
		if v == value {
			inarr = true
			break
		}
	}

	return inarr
}

func AddUrl(url string,regex string) {
	if _config == nil {
		log.Fatal("add url before new filter")
		return
	}
	_config.Items[url] = regex
}

// Middleware validates CSRF token.
func Middleware(options Options) gin.HandlerFunc {
	ignoreMethods := options.IgnoreMethods
	errorFunc := options.ErrorFunc
	tokenGetter := options.TokenGetter

	if default_store == nil{
		default_store = sessions.NewCookieStore([]byte(options.Secret))
	}else{
		log.Fatal("bind filter more than once")
		return nil
	}

	if _config == nil {
		_config = utils.New(options.Blacklist)
	}

	if ignoreMethods == nil {
		ignoreMethods = defaultIgnoreMethods
	}

	if errorFunc == nil {
		errorFunc = defaultErrorFunc
	}

	if tokenGetter == nil {
		tokenGetter = defaultTokenGetter
	}

	return func(c *gin.Context) {
		if _config.AllowAccess(c) {
			return
		}

		session,_ := default_store.Get(c,DefaultCookieKey)
		c.Set(csrfSecret, options.Secret)

		if inArray(ignoreMethods, c.Request.Method) {
			c.Next()
			return
		}

		var salt string

		if s, ok := session.Get(csrfSalt).(string); !ok || len(s) == 0 {
			c.Next()
			return
		} else {
			salt = s
		}

		session.Delete(csrfSalt)

		token := tokenGetter(c)

		if tokenize(options.Secret, salt) != token {
			errorFunc(c)
			c.Abort()
			return
		}
	}
}

// GetToken returns a CSRF token.
func GetToken(c *gin.Context) string {
	session,_ := default_store.Get(c,DefaultCookieKey)
	secret := c.MustGet(csrfSecret).(string)

	if t, ok := c.Get(csrfToken); ok {
		return t.(string)
	}

	salt := uniuri.New()
	token := tokenize(secret, salt)
	session.Set(csrfSalt, salt)
	session.Save()
	c.Set(csrfToken, token)

	return token
}
