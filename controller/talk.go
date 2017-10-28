package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//TalkContent is sending from iPhone or RaspberyyPi
type TalkContent struct {
	Text string `json:"text"`
}

//Talk is function of talking lapping talking api
func Talk(c *gin.Context) {
	talkContent := new(TalkContent)
	c.BindJSON(talkContent)
	c.JSON(http.StatusOK, talkContent)
}
