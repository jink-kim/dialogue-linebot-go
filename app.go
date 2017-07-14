package main

import (
    "encoding/json"
    "bytes"
    "log"
    "io/ioutil"
    "net/http"
    "fmt"
    "github.com/line/line-bot-sdk-go/linebot"
)

const (
   Channel_Secret = "Channel_Secret" //LINE Bot Channel Secret.
   Channel_Access_Token = "Channel_Access_Token"//LINE Bot Channel Access Token.
   Docomo_Dialogue_Token = "Docomo_Dialogue_Token" //DOCOMO Dialogue Token.
   Docomo_Dialogue_Url = "https://api.apigw.smt.docomo.ne.jp/dialogue/v1/dialogue?APIKEY=%s" //DOCOMO Request URL.
)

/*Docomo dialogue taking response
   https://dev.smt.docomo.ne.jp/?p=docs.api.page&api_name=dialogue&p_name=api_1#tag01*/
type Dialogue struct {
    Utt     string `json:"utt"`
    Yomi    string `json:"yomi"`
    Mode    string `json:"mode"`
    Da      string `json:"da"`
    Context string `json:"context"`
}

func main() {

    bot, err := linebot.New (
        Channel_Secret,
        Channel_Access_Token,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Setup HTTP Server for receiving requests from LINE platform
    http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {

    //ParseRequest
        events, err := bot.ParseRequest(req)
        if err != nil {
            if err == linebot.ErrInvalidSignature {
                w.WriteHeader(400)
            } else {
                w.WriteHeader(500)
            }
            return
        }
        for _, event := range events {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                if resMessage := getDocomoDialogueMsg(message.Text); resMessage != "" {
                    postMessage := linebot.NewTextMessage(resMessage)
                        if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
                            log.Print(err)
                        }
                }
            }
        }
    }
})

    // For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
    go http.ListenAndServe(":PORT NUM", nil)
    if err := http.ListenAndServeTLS(":PORT NUM", "SSL(PEM) file", "SSL KEY", nil); err != nil {
        log.Fatal(err)
    }

}

//DocomoDialogue API.
func getDocomoDialogueMsg(reqMessage string) (message string) {

   var docomo_url string //DocomoDialogue URL.
   var dia Dialogue      //DocomoDialogue JSON.

   docomo_url = fmt.Sprintf(Docomo_Dialogue_Url, Docomo_Dialogue_Token)

    //Docomo dialogue taking Request.
    client := &http.Client{}
    jsonStr := `{"utt":"` + reqMessage + `"}`
    req, _ := http.NewRequest (
        "POST",
        docomo_url,
        bytes.NewBuffer([]byte(jsonStr)),
    )
    req.Header.Set("Content-Type", "application/json")

    resp, _ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    //JSON Decode.
    err := json.Unmarshal(body, &dia)
    if err != nil {
        panic(err)
    }
    message = dia.Utt

   return
}
