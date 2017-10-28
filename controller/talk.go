package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/gin-gonic/gin"
)

const (
	//DOCOMOURL is docomo's chatting api endpoint
	DOCOMOURL = "https://api.apigw.smt.docomo.ne.jp/dialogue/v1/dialogue"
	//DOCOMOAPIKEY is docomo's chatting api key
	DOCOMOAPIKEY = "3163666e31357858574e786733386432482e4b6a73664c7656514d315567716d4337527a42656b44774130"
)

//TalkContent is sending from iPhone or RaspberyyPi
type TalkContent struct {
	Text string `json:"text"`
}

//DocomoJSON is docomo response json
type DocomoJSON struct {
	Response string `json:"utt"`
	Read     string `json:"yomi"`
	Mode     string `json:"dialog"`
	Da       string `json:"da"`
	Context  string `json:"context"`
}

//Talk is function of talking lapping talking api
func Talk(c *gin.Context) {
	talkContent := new(TalkContent)
	c.BindJSON(talkContent)
	c.JSON(http.StatusOK, talkContent)
}

//Chatting return chatting response
func Chatting(text string) (DocomoJSON, error) {
	resp, err := http.Get(DOCOMOURL + "?APIKEY=" + DOCOMOAPIKEY)
	if err != nil {
		log.Fatal(err)
		return DocomoJSON{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return DocomoJSON{}, err
	}

	docomo := new(DocomoJSON)
	if err := json.Unmarshal(bytes, &docomo); err != nil {
		log.Fatal(err)
		return DocomoJSON{}, err
	}

	return *docomo, nil
}
