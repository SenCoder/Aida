package common

//var HttpToSocket chan Msg
//var SocketToHttp chan Msg


var HttpToSocket chan string
var SocketToHttp chan string

//type Msg struct {
//	Meta   map[string]interface{} `json:"meta"`
//	Content interface{}            `json:"content"`
//	Debug string            	`json:"debug"`
//}


type Msg struct {
	Meta   map[string]interface{} `json:"meta"`
	IntentType	string	`json:"intentType"`
	Value 		string	`json:"value"`
}

const (

	SWITCH_CHANNEL_BY_NAME = "SWITCH_CHANNEL_BY_NAME"
	SWITCH_CHANNEL_BY_NUMBER = "SWITCH_CHANNEL_BY_NUMBER"
	SET_VOLUME = "SET_VOLUME"
	WATCH_MOVIE = "WATCH_MOVIE"
	OPERATE_APP = "OPERATE_APP"
	LISTEN_MUSIC = "LISTEN_MUSIC"
)