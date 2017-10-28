package controller

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/youtangai/Otomo_backend/model"
	"github.com/youtangai/Otomo_backend/speak"
	"github.com/youtangai/Otomo_backend/storage"
)

//Collect is collection indoor info function
func Collect(c *gin.Context) {
	indoorInfo := new(model.IndoorInfoJSON)
	c.BindJSON(indoorInfo)
	err := RecomendMoving(*indoorInfo)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, "faild")
		return
	}
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

//RecomendMoving is if humidy very high recomend speak
func RecomendMoving(info model.IndoorInfoJSON) error {
	if info.Humidity >= 60 {
		err := speak.Speak("こんな蒸し蒸しした所居れないよ！早く引っ越そうよ！")
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
