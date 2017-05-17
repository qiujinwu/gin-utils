package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qiujinwu/gin-utils/sessions"
	"gopkg.in/redis.v5"
)

func main() {
	r := gin.Default()
	// cookie session，所有的内容都存在客户端的cookie中
	// var store = sessions.NewCookieStore([]byte("something-very-secret"))

	// 存本地文件
	// var store = sessions.NewFilesystemStore("/tmp/11", []byte("something-very-secret"))

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "", // no password set
		DB: 2,  // use default DB
	})
	var store,_ = sessions.NewRediStore(client,[]byte("something-very-secret"))


	r.GET("/", func(c *gin.Context) {
		session, _ := store.Get(c,"session-name")
		session2, _ := store.Get(c,"session-name2")
		var count,count2 int
		v := session.Get("count")
		v2 := session2.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}
		session.Set("count", count)

		if v2 == nil {
			count2 = 0
		} else {
			count2 = v2.(int)
			count2 += 1
		}
		session2.Set("count", count2)

		// 保存session
		// session.Save()
		// session2.Save()
		sessions.Save(c)
		c.JSON(200, gin.H{"count": count,"count2": count2})
	})
	r.Run(":4000") // listen and serve on 0.0.0.0:8080
}
