package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	router.POST("/", func(c *gin.Context) {
		text := c.PostForm("text")
		log.Printf("text: %#v", text)
		c.String(http.StatusOK, "Hello World")
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Run("0.0.0.0:8080")
}
