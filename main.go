package main

import (
	"server/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.GetEnv()
	initializers.ConnectDB()
}

func main() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "Hello, world",
		})
	})
	router.Run()
}
