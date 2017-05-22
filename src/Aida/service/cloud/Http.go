package cloud

import (
	"net/http"
	"Aida/service/common"
	"github.com/b00giZm/golexa"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"math/rand"
)

const (
	INTENT_SWITCH_CHANNEL = "MytvSwitchChannel"
	INTENT_SET_VOLUME = "MytvSetVolume"
	INTENT_WATCH_MOVIE = "MytvWatchMovie"
	INTENT_OPERATE_APP = "MytvOperateApp"
	INTENT_LISTEN_MUSIC = "MytvListenMusic"
)

const (
	CHANNEL_NAME = "ChannelName"
	CHANNEL_NUMBER = "ChannelNumber"
	VOLUME_NUMBER = "VolumeNumber"
	MOVIE_NAME = "MovieName"
	APP_NAME = "AppName"
	MUSIC_NAME = "MusicName"
)

var GREETING = [...]string {
	"welcome to my tv alexa skill, what can i do for you",
	"welcome to my tv, it's really my honor to serve you",
	"nice to meet you, i'm ready to serve you",
	"nice to meet you again, i'm very pleased to serve you"}


type Slot struct {
	Name string
	Value string
}

type Intent struct {
	Name string
	Slots []Slot
}

var appMatch map[string]string

func init() {
	appMatch = make(map[string]string, 0)
	//appMatch["Launcher"] = "Launcher"
}


func match(raw string) string {
	return raw
}

func handler(writer http.ResponseWriter, request *http.Request) {

	//Chan := make(chan string)
	app := golexa.Default()
	app.OnLaunch(func(alexa *golexa.Alexa, req *golexa.Request, session *golexa.Session) *golexa.Response {

		greetIndex := rand.Intn(100) % len(GREETING)
		randomGreet := GREETING[greetIndex]

		return alexa.Response().AddPlainTextSpeech(randomGreet)
	})

	app.OnIntent(func(alexa *golexa.Alexa, intent *golexa.Intent, req *golexa.Request, session *golexa.Session) *golexa.Response {

		var reply string
		var msg string
		intentName := req.Intent.Name

		switch {
		case intentName == INTENT_SWITCH_CHANNEL:
			fmt.Println("MytvSwitchChannel")
			reply += "ok, my tv will switch to channel "
		case intentName == INTENT_SET_VOLUME:
			fmt.Println("MytvSetVolume ")
			reply += "ok, my tv will set volume to "
		case intentName == INTENT_WATCH_MOVIE:
			fmt.Println("MytvWatchMovie ")
			reply += "ok, my tv will search "
		case intentName == INTENT_OPERATE_APP:
			fmt.Println(INTENT_OPERATE_APP)
			reply += "ok, my tv will open "
		case intentName == INTENT_LISTEN_MUSIC:
			fmt.Println(INTENT_LISTEN_MUSIC)
			reply += "ok, my tv will play "
		default:
			log.Println("unknown intent")
			reply = "Sorry, I am confused. Are you really speaking English ?"
			goto GETOUT
		}

		for slot := range req.Intent.Slots {
			if len(req.Intent.Slots[slot].Value) > 0 {
				slotName := req.Intent.Slots[slot].Name
				switch {
				case slotName == CHANNEL_NAME:
					reply += req.Intent.Slots[slot].Value
					msg += common.SWITCH_CHANNEL_BY_NAME
					msg += ":"
					msg += req.Intent.Slots[slot].Value
				case slotName == CHANNEL_NUMBER:
					reply += req.Intent.Slots[slot].Value
					msg += common.SWITCH_CHANNEL_BY_NUMBER
					msg += ":"
					msg += req.Intent.Slots[slot].Value
				case slotName == VOLUME_NUMBER:
					reply += req.Intent.Slots[slot].Value
					msg += common.SET_VOLUME
					msg += ":"
					msg += req.Intent.Slots[slot].Value
				case slotName == MOVIE_NAME:
					reply += req.Intent.Slots[slot].Value
					msg += common.WATCH_MOVIE
					msg += ":"
					msg += req.Intent.Slots[slot].Value
				case slotName == APP_NAME:
					reply += req.Intent.Slots[slot].Value
					msg += common.OPERATE_APP
					msg += ":"
					msg += req.Intent.Slots[slot].Value

				case slotName == MUSIC_NAME:
					reply += req.Intent.Slots[slot].Value
					msg += common.LISTEN_MUSIC
					msg += ":"
					msg += req.Intent.Slots[slot].Value

				default:
					reply = "Are you really speaking English ?"
					log.Println("unknown slot name")
					goto GETOUT
				}
			}
		}

		if len(msg) <= 0 {
			reply = "Are you really speaking English ?"
			goto GETOUT
		}

		go func(){
			common.HttpToSocket <- msg
			log.Println("http send msg to socket success")
		}()

		GETOUT:
		return alexa.Response().AddPlainTextSpeech(reply)
	})

	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("handler post error: ", err.Error())
	}

	fmt.Println("handler post result: ", string(result))
	request.Body.Close()

	msg := json.RawMessage(result)
	response,_ := app.Process(msg)
	re,_ := json.Marshal(response)
	writer.Write(re)
}

/*
func handler(writer http.ResponseWriter, request *http.Request) {
	msg := "event:====hello, Aida &* %$#@ ,.-=yfhjfhjfhfgghdhgffhghjghjffggdgdhdgjfjghjghjhffgdgryutughjfhrythj我们"
	common.HttpToSocket <- msg
	writer.Write([]byte("hello"))
}

*/
func StartHttpServer(port string)  {
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}


func StartHttpsServer(port string)  {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(port, "server.crt", "server.key", nil)
}
