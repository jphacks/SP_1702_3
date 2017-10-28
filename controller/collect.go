package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtangai/Otomo_backend/model"
	"github.com/youtangai/Otomo_backend/storage"
)

//Collect is collection indoor info function
func Collect(c *gin.Context) {
	indoorInfo := new(model.IndoorInfoJSON)
	c.BindJSON(indoorInfo)
	InsertIndoorInfo(*indoorInfo)
	c.String(http.StatusOK, "success")
}

//InsertIndoorInfo is saving function
func InsertIndoorInfo(info model.IndoorInfoJSON) {
	record := new(model.IndoorInfo)
	record.Temperature = info.Temperature
	record.Humidity = info.Humidity
	record.Illumination = info.Illumination
	db := storage.GetDBContext()
	db.Create(&record)
}
