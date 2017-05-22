package main

import (
	"Aida/service/utils"
	"Aida/service/cloud"
	"Aida/service/common"
	"runtime"
)

func main()  {
	runtime.GOMAXPROCS(6)
	common.HttpToSocket = make(chan string)
	common.SocketToHttp = make(chan string)
	msg := make(chan bool, 0)
	go cloud.StartHttpServer(utils.Cfg.HttpPort) // alexa
	go cloud.StartSocketServer(utils.Cfg.SocketHost, utils.Cfg.BeatingInterval)
	<- msg
}
