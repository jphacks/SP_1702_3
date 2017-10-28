package speak

import (
	"bytes"
	"log"
	"net/http"
)

const (
	//SPEAKURL is url
	SPEAKURL = "http://192.168.2.45:1234/speak"
)

//Speak is speaker function
func Speak(text string) error {
	jsonStr := `{"text":"` + text + `"}`
	req, err := http.NewRequest(
		"POST",
		SPEAKURL,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
