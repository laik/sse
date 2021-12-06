package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type watcherEvent struct {
	Type   string      `json:"type"`
	Object interface{} `json:"object"`
	URL    string      `json:"url"`
	Status int         `json:"status"`
}

func main() {
	route := gin.New()

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	route.GET("/watch", watchfunc)
	route.GET("/watcher/watch", watchfunc)

	route.Run()
}

func watchfunc(c *gin.Context) {
	ticker := time.NewTicker(10 * time.Second)
	x := RandStringRunes(10)
	fmt.Printf("recv connect %s \n", x)
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			c.SSEvent("", watcherEvent{
				Type:   "ping",
				Object: x,
			})
		case <-c.Writer.CloseNotify():
			fmt.Printf("close %s\n", x)
			return false
		}

		return true
	})
}
