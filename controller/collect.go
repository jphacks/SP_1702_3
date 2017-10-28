package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//IndoorInfo is users indoor information
type IndoorInfo struct {
	Temperature  float32 `json:"temperature"`
	Humidity     int     `json:"humidity"`
	Illumination int     `json:"illumination"`
}

//Collect is collection indoor info function
func Collect(c *gin.Context) {
	indoorInfo := new(IndoorInfo)
	c.BindJSON(indoorInfo)
	c.JSON(http.StatusOK, indoorInfo)
}
