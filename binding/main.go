package binding

import (
	"github.com/gin-gonic/gin"
)

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON binding
// 		"application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func Bind(c * gin.Context,obj interface{}) error {
	b := Default(c.Request.Method, c.ContentType())
	return BindWith(c, obj, b)
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
func BindJSON(c * gin.Context,obj interface{}) error {
	return BindWith(c, obj, JSON)
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func BindWith(c * gin.Context,obj interface{}, b Binding) error {
	if err := b.Bind(c.Request, obj); err != nil {
		return err
	}
	return nil
}
