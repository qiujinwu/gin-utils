package utils

import (
	"github.com/gin-gonic/gin"
)

type IRoutesWarapper struct {
	router   gin.IRoutes
	base_url string
}

type RuleFunc func(string,string)
type RuleFuncs []RuleFunc

func NewIRoutesWarapper(router gin.IRoutes) *IRoutesWarapper {
	base_url := ""
	switch v := router.(type) {
	case *gin.RouterGroup:{
		base_url = v.BasePath()
	}
	}

	return &IRoutesWarapper{
		router:   router,
		base_url: base_url,
	}
}

func (irouter *IRoutesWarapper) callRuleHandles(method string,url string,regex string,rule_handles RuleFuncs){
	if rule_handles != nil{
		var newurl string
		if url[0] == '/'{
			newurl = method + irouter.base_url + url
			if regex != ""{
				regex = method + irouter.base_url + regex
			}
		}else{
			newurl = method + irouter.base_url + "/" + url
			if regex != "" {
				regex = method + irouter.base_url + "/" + regex
			}
		}
		for _, value := range rule_handles {
			value(newurl,regex)
		}
	}
}

func (irouter *IRoutesWarapper) GET(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	return irouter.GET2(url,"",rule_handles,handlers...)
}

func (irouter *IRoutesWarapper) GET2(url string,regex string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("GET",url,regex,rule_handles,)
	irouter.router.GET(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) POST(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("POST",url,"",rule_handles)
	irouter.router.POST(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) DELETE(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("DELETE",url,"",rule_handles)
	irouter.router.DELETE(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) PATCH(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("PATCH",url,"",rule_handles)
	irouter.router.PATCH(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) PUT(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("PUT",url,"",rule_handles)
	irouter.router.PUT(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) OPTIONS(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("OPTIONS",url,"",rule_handles)
	irouter.router.OPTIONS(url, handlers...)
	return irouter
}

func (irouter *IRoutesWarapper) HEAD(url string,
	rule_handles RuleFuncs,handlers ...gin.HandlerFunc) *IRoutesWarapper {
	irouter.callRuleHandles("HEAD",url,"",rule_handles)
	irouter.router.HEAD(url, handlers...)
	return irouter
}