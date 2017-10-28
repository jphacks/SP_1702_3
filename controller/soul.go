package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtangai/Otomo_backend/model"
	"github.com/youtangai/Otomo_backend/storage"
)

//Soul is soulcontroller
func Soul(c *gin.Context) {
	soul := new(model.SoulJSON)
	c.BindJSON(soul)
	ChangeSoul(*soul)
	c.JSON(http.StatusOK, "success")
}

//ChangeSoul is changes otomo status
func ChangeSoul(soul model.SoulJSON) {
	db := storage.GetDBContext()
	db.Model(&model.Soul{}).Where("user_id = ?", soul.UserID).Update("on_device", soul.OnDevice)
}
