package main

import (
	"Aida/service/cloud"
	"Aida/service/common"
	"runtime"
	"Aida/service/conf"
)

func main()  {
	runtime.GOMAXPROCS(6)
	common.HttpToSocket = make(chan string)
	common.SocketToHttp = make(chan string)
	msg := make(chan bool, 0)
	go cloud.StartHttpServer(config.Cfg.HttpPort) // alexa
	go cloud.StartSocketServer(config.Cfg.SocketHost, config.Cfg.BeatingInterval)
	<- msg
}
