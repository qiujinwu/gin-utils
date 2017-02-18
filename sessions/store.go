package sessions

import (
	"github.com/gin-gonic/gin"
)

type Store interface {
	// Get should return a cached session.
	Get(c *gin.Context, name string) (*SessionImp, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(c *gin.Context, name string) (*SessionImp, error)

	// Save should persist session to the underlying store implementation.
	Save(c *gin.Context, s* SessionImp) error

	// 删除cookie
	Delete(c *gin.Context, name string) error
}