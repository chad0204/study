package main

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {

	// json
	log.Debug().Msgf("web start success :%s", "localhost:8080/ping")


	c := cache.New(5*time.Minute, 10*time.Minute)
	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cache.DefaultExpiration)

	http()

}

//https://gin-gonic.com/zh-cn/docs/quickstart/
func http()  {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务

}
