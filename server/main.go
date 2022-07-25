package main

import (
	"gin-example/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	Engine *gin.Engine
}

func main() {
	c := NewComsumer("localhost:6379", "", 0)
	go c.Consume(c.rdb.Options().Addr)

	RGin := Gin{
		Engine: gin.Default(),
	}
	RGin.Engine.GET(
		"/k",
		//?interval=1
		func(c *gin.Context) {
			interval := c.DefaultQuery("interval", "1")
			c.JSON(http.StatusOK, gin.H{
				"message":  "unsupported interval",
				"interval": interval,
			})
		})
	RGin.Engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (g *Gin) Get(
	post string,
	kline *model.KLine,
) {
	g.Engine.GET(post, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"kline": kline,
		})
	})
}

// response:
// first line： http/1.1 200 OK
// headers

// body
