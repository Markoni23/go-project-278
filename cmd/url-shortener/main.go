package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}

func main() {
	router := setupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
