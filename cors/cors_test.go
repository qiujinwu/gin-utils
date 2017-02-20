package cors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(config Config) *gin.Engine {
	router := gin.New()
	router.Use(New(config))
	router.GET("/", func(c *gin.Context) {
		c.String(200, "get")
	})
	router.POST("/", func(c *gin.Context) {
		c.String(200, "post")
	})
	router.PATCH("/", func(c *gin.Context) {
		c.String(200, "patch")
	})
	return router
}

func performRequest(r http.Handler, method, origin string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, "/", nil)
	if len(origin) > 0 {
		req.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestConfigAddAllow(t *testing.T) {
	config := Config{}
	config.AddAllowMethods("POST")
	config.AddAllowMethods("GET", "PUT")
	config.AddExposeHeaders()

	config.AddAllowHeaders("Some", " cool")
	config.AddAllowHeaders("header")
	config.AddExposeHeaders()

	config.AddExposeHeaders()
	config.AddExposeHeaders("exposed", "header")
	config.AddExposeHeaders("hey")

	assert.Equal(t, config.AllowMethods, []string{"POST", "GET", "PUT"})
	assert.Equal(t, config.AllowHeaders, []string{"Some", " cool", "header"})
	assert.Equal(t, config.ExposeHeaders, []string{"exposed", "header", "hey"})

}

func TestBadConfig(t *testing.T) {
	assert.Panics(t, func() { New(Config{}) })
	assert.Panics(t, func() {
		New(Config{
			AllowAllOrigins: true,
			AllowOrigins:    []string{"http://google.com"},
		})
	})
	assert.Panics(t, func() {
		New(Config{
			AllowAllOrigins: true,
			AllowOriginFunc: func(origin string) bool { return false },
		})
	})
	assert.Panics(t, func() {
		New(Config{
			AllowOrigins: []string{"google.com"},
		})
	})
}

func TestNormalize(t *testing.T) {
	values := normalize([]string{
		"http-Access ", "Post", "POST", " poSt  ",
		"HTTP-Access", "",
	})
	assert.Equal(t, values, []string{"http-access", "post", ""})

	values = normalize(nil)
	assert.Nil(t, values)

	values = normalize([]string{})
	assert.Equal(t, values, []string{})
}

func TestConvert(t *testing.T) {
	methods := []string{"Get", "GET", "get"}
	headers := []string{"X-CSRF-TOKEN", "X-CSRF-Token", "x-csrf-token"}

	assert.Equal(t, []string{"GET", "GET", "GET"}, convert(methods, strings.ToUpper))
	assert.Equal(t, []string{"X-Csrf-Token", "X-Csrf-Token", "X-Csrf-Token"}, convert(headers, http.CanonicalHeaderKey))
}

func TestGenerateNormalHeaders_AllowAllOrigins(t *testing.T) {
	header := generateNormalHeaders(Config{
		AllowAllOrigins: false,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Origin"), "")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 1)

	header = generateNormalHeaders(Config{
		AllowAllOrigins: true,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Origin"), "*")
	assert.Equal(t, header.Get("Vary"), "")
	assert.Len(t, header, 1)
}

func TestGenerateNormalHeaders_AllowCredentials(t *testing.T) {
	header := generateNormalHeaders(Config{
		AllowCredentials: true,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Credentials"), "true")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestGenerateNormalHeaders_ExposedHeaders(t *testing.T) {
	header := generateNormalHeaders(Config{
		ExposeHeaders: []string{"X-user", "xPassword"},
	})
	assert.Equal(t, header.Get("Access-Control-Expose-Headers"), "X-User,Xpassword")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestGeneratePreflightHeaders(t *testing.T) {
	header := generatePreflightHeaders(Config{
		AllowAllOrigins: false,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Origin"), "")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 1)

	header = generateNormalHeaders(Config{
		AllowAllOrigins: true,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Origin"), "*")
	assert.Equal(t, header.Get("Vary"), "")
	assert.Len(t, header, 1)
}

func TestGeneratePreflightHeaders_AllowCredentials(t *testing.T) {
	header := generatePreflightHeaders(Config{
		AllowCredentials: true,
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Credentials"), "true")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestGeneratePreflightHeaders_AllowedMethods(t *testing.T) {
	header := generatePreflightHeaders(Config{
		AllowMethods: []string{"GET ", "post", "PUT", " put  "},
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Methods"), "GET,POST,PUT")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestGeneratePreflightHeaders_AllowedHeaders(t *testing.T) {
	header := generatePreflightHeaders(Config{
		AllowHeaders: []string{"X-user", "Content-Type"},
	})
	assert.Equal(t, header.Get("Access-Control-Allow-Headers"), "X-User,Content-Type")
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestGeneratePreflightHeaders_MaxAge(t *testing.T) {
	header := generatePreflightHeaders(Config{
		MaxAge: 12 * time.Hour,
	})
	assert.Equal(t, header.Get("Access-Control-Max-Age"), "43200") // 12*60*60
	assert.Equal(t, header.Get("Vary"), "Origin")
	assert.Len(t, header, 2)
}

func TestValidateOrigin(t *testing.T) {
	cors := newCors(Config{
		AllowAllOrigins: true,
	})
	assert.True(t, cors.validateOrigin("http://google.com"))
	assert.True(t, cors.validateOrigin("https://google.com"))
	assert.True(t, cors.validateOrigin("example.com"))

	cors = newCors(Config{
		AllowOrigins: []string{"https://google.com", "https://github.com"},
		AllowOriginFunc: func(origin string) bool {
			return (origin == "http://news.ycombinator.com")
		},
	})
	assert.False(t, cors.validateOrigin("http://google.com"))
	assert.True(t, cors.validateOrigin("https://google.com"))
	assert.True(t, cors.validateOrigin("https://github.com"))
	assert.True(t, cors.validateOrigin("http://news.ycombinator.com"))
	assert.False(t, cors.validateOrigin("http://example.com"))
	assert.False(t, cors.validateOrigin("google.com"))
}

func TestPassesAllowedOrigins(t *testing.T) {
	router := newTestRouter(Config{
		AllowOrigins:     []string{"http://google.com"},
		AllowMethods:     []string{" GeT ", "get", "post", "PUT  ", "Head", "POST"},
		AllowHeaders:     []string{"Content-type", "timeStamp "},
		ExposeHeaders:    []string{"Data", "x-User"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://github.com"
		},
	})

	// no CORS request, origin == ""
	w := performRequest(router, "GET", "")
	assert.Equal(t, w.Body.String(), "get")
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, w.Header().Get("Access-Control-Expose-Headers"))

	// allowed CORS request
	w = performRequest(router, "GET", "http://google.com")
	assert.Equal(t, w.Body.String(), "get")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "http://google.com")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "")
	assert.Equal(t, w.Header().Get("Access-Control-Expose-Headers"), "Data,X-User")

	// deny CORS request
	w = performRequest(router, "GET", "https://google.com")
	assert.Equal(t, w.Code, 403)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, w.Header().Get("Access-Control-Expose-Headers"))

	// allowed CORS prefligh request
	w = performRequest(router, "OPTIONS", "http://github.com")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "http://github.com")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "GET,POST,PUT,HEAD")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type,Timestamp")
	assert.Equal(t, w.Header().Get("Access-Control-Max-Age"), "43200")

	// deny CORS prefligh request
	w = performRequest(router, "OPTIONS", "http://example.com")
	assert.Equal(t, w.Code, 403)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.Empty(t, w.Header().Get("Access-Control-Max-Age"))
}

func TestPassesAllowedAllOrigins(t *testing.T) {
	router := newTestRouter(Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{" Patch ", "get", "post", "POST"},
		AllowHeaders:     []string{"Content-type", "  testheader "},
		ExposeHeaders:    []string{"Data2", "x-User2"},
		AllowCredentials: false,
		MaxAge:           10 * time.Hour,
	})

	// no CORS request, origin == ""
	w := performRequest(router, "GET", "")
	assert.Equal(t, w.Body.String(), "get")
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, w.Header().Get("Access-Control-Expose-Headers"))

	// allowed CORS request
	w = performRequest(router, "POST", "example.com")
	assert.Equal(t, w.Body.String(), "post")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	assert.Equal(t, w.Header().Get("Access-Control-Expose-Headers"), "Data2,X-User2")
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))

	// allowed CORS prefligh request
	w = performRequest(router, "OPTIONS", "https://facebook.com")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "PATCH,GET,POST")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type,Testheader")
	assert.Equal(t, w.Header().Get("Access-Control-Max-Age"), "36000")
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
}
