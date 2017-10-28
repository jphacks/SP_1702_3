package main

import (
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/gin-gonic/gin"
	"github.com/youtangai/Otomo_backend/controller"
)

func main() {
	router := gin.Default()

	router.GET("hc", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm Healty!\n")
	})
	router.POST("/collect", controller.Collect)
	router.POST("/talk", controller.Talk)

	log.Fatal(router.Run(":9000"))
}
