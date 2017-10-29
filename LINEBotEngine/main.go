package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	dproxy "github.com/koron/go-dproxy"
)

type LINEEvents struct {
	Events []LINEMessaging `json:"events"`
}

type LINEMessaging struct {
	ReplyToken string      `json:"replyToken"`
	Type       string      `json:"type"`
	TimeStamp  int         `json:"timestamp"`
	Source     SourceType  `json:"source"`
	Message    MessageType `json:"message"`
}

type SourceType struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type MessageType struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type ReplyMessageType struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LINEReply struct {
	ReplyToken    string             `json:"replyToken"`
	ReplyMessages []ReplyMessageType `json:"messages"`
}

//DocomoJSON is docomo response json
type DocomoJSON struct {
	Response string `json:"utt"`
	Read     string `json:"yomi"`
	Mode     string `json:"mode"`
	Da       string `json:"da"`
	Context  string `json:"context"`
}

//ShopInfo is
type ShopInfo struct {
	ShopName string `json:"shop_name"`
	Address  string `json:"address"`
	Tel      string `json:"tel"`
}

type PushMessage struct {
	To            string             `json:"to"`
	ReplyMessages []ReplyMessageType `json:"messages"`
}

type Auth struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	Limit     string `json:"expires_in"`
}

const (
	//URL is url
	URL = "https://api.apigw.smt.docomo.ne.jp/dialogue/v1/dialogue"
	//API is api
	API = "3163666e31357858574e786733386432482e4b6a73664c7656514d315567716d4337527a42656b44774130"
	//TOKEN is token
	TOKEN = "Jz4Jj6G6QDjWolOw/I0o260NoC85VAFCgSHNVr8JXxfRZAITyWnU+gkcS7IQynF5MlLkYe6Oi429kMvd5L+gl1YMJTgEiPwHuKz328LFVrRe+pokMdCTytdaVjW7KgMo415e0QSyuSrt3isWNf/NUwdB04t89/1O/w1cDnyilFU="
	//SHOPURL is url
	SHOPURL = "https://glacial-anchorage-80532.herokuapp.com/food/food/"

	LIFFULEAPIURL  = "https://api.homes.co.jp/v1/realestate_article/search"
	LIFFULEAUTHURL = "https://auth.homes.co.jp/token"
	USER           = "OTk4MmE5ZTBjOGYyMTJiM2I4ZmExZjFlZTk0MmRk"
	PASS           = "ZDVhMWE5N2Q2YzNmOGFlZjRlZGU3MzlkODVlZTY3"
)

var (
	source = ""
)

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "faild!")
	}

	line := new(LINEEvents)
	if err := json.Unmarshal(jsonBytes, &line); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "faild!")
	}

	text := line.Events[0].Message.Text
	log.Println("text is", text)
	replyText := ""
	if strings.Contains(text, "/shop") {
		trimed := strings.TrimPrefix(text, "/shop ")
		log.Println("trimed is", trimed)
		locate := strings.Split(trimed, ",")
		log.Println("locate is", locate)
		shopInfo := Shopping(locate[0], locate[1])
		replyText = shopInfo
		log.Println("replyText is", replyText)
	} else if strings.Contains(text, "へや") {
		replyText = Lifful()
	} else {
		docomo, err := Chatting(text)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "faild!")
		}
		replyText = docomo.Read
	}

	reply := new(LINEReply)
	reply.ReplyToken = line.Events[0].ReplyToken

	source = line.Events[0].Source.UserID
	log.Println(source)

	message := ReplyMessageType{Type: "text", Text: replyText}
	messages := []ReplyMessageType{message}
	reply.ReplyMessages = messages
	jsonBytes, err = json.Marshal(reply)

	req, err := http.NewRequest(
		"POST",
		"https://api.line.me/v2/bot/message/reply",
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
	log.Println("/reply accessed")
}

//Chatting return chatting response
func Chatting(text string) (DocomoJSON, error) {
	jsonStr := `{"utt":"` + text + `"}`
	url := URL + "?APIKEY=" + API

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		log.Fatal(err)
		return DocomoJSON{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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

func Shopping(lat, lon string) string {
	resp, err := http.Get(SHOPURL + "?lat=" + lat + "&lon=" + lon)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	shopInfo := new(ShopInfo)
	err = json.Unmarshal(byteJSON, &shopInfo)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	str := "１番近いお店は " + shopInfo.ShopName + "\n住所は " + shopInfo.Address + "\n電話番号は " + shopInfo.Tel + "\nだよ。"
	return str
}

func Lifful() string {
	tokenValues := url.Values{}
	tokenValues.Add("grant_type", "client_credentials")

	tokenReq, err := http.NewRequest(
		"POST",
		LIFFULEAUTHURL,
		strings.NewReader(tokenValues.Encode()),
	)
	if err != nil {
		log.Fatal(err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.SetBasicAuth(USER, PASS)

	tokenClient := &http.Client{}
	tokenResp, err := tokenClient.Do(tokenReq)
	if err != nil {
		log.Fatal(err)
	}
	defer tokenResp.Body.Close()

	byteTokenJSON, err := ioutil.ReadAll(tokenResp.Body)
	auth := new(Auth)
	json.Unmarshal(byteTokenJSON, &auth)

	apiValues := url.Values{}
	apiValues.Add("hits", "1")
	//apiValues.Add("pref_id", "12")
	apiValues.Add("city_id", "12245")

	apiReq, err := http.NewRequest(
		"GET",
		LIFFULEAPIURL,
		strings.NewReader(apiValues.Encode()),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Content-Type 設定
	apiReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	apiReq.Header.Set("Authorization", "Bearer "+auth.Token)

	apiClient := &http.Client{}
	apiResp, err := apiClient.Do(apiReq)
	if err != nil {
		log.Fatal(err)
	}
	defer apiResp.Body.Close()

	byteJSON, err := ioutil.ReadAll(apiResp.Body)
	var result interface{}
	json.Unmarshal(byteJSON, &result)
	v := dproxy.New(result)
	name, _ := v.M("row_set").A(0).M("realestate_article_name").String()
	addr, _ := v.M("row_set").A(0).M("full_address").String()
	money, _ := v.M("row_set").A(0).M("money_room_text").String()
	return fmt.Sprintf("%s に\n%s の\n%s って\n名前のお家があるよ！！", addr, money, name)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "62043"
	}

	http.HandleFunc("/hc", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I'm Healty\n")
		log.Println("/hc accessed")
	})
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		push := new(PushMessage)
		message := ReplyMessageType{Type: "text", Text: "こんなジメジメしたところはいやだよー，はやく引っ越そうよ！"}
		messages := []ReplyMessageType{message}
		push.To = source
		push.ReplyMessages = messages
		jsonBytes, err := json.Marshal(push)

		req, err := http.NewRequest(
			"POST",
			"https://api.line.me/v2/bot/message/push",
			bytes.NewReader(jsonBytes),
		)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+TOKEN)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(body))
		log.Println("/push accessed")
	})
	http.HandleFunc("/reply", handler)

	log.Println("start Listen on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
