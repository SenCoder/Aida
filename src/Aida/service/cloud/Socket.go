package cloud

import (
	"net"
	"log"
	"Aida/service/common"
	"Aida/service/utils"
)

func StartSocketServer(host string, timeinterval int){
	netListen, err := net.Listen("tcp", host)
	utils.CheckError(err)
	defer netListen.Close()
	utils.Log("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		utils.Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn, timeinterval)
	}
}


//handle the connection
func handleConnection(conn net.Conn, timeout int ) {

	defer conn.Close()

	tmpBuffer := make([]byte, 0)
	buffer := make([]byte, 1024)
	messnager := make(chan byte)
	//connOK := true

	overMsg := make(chan bool)

	go func(){
		for {
			log.Println("Socket Server waiting for next Alexa event")
			var msg string

			select {
			case msg = <- common.HttpToSocket:
				log.Println("Socket Server get event:", msg)
			case <-overMsg:
				return
			}
			n, err := conn.Write(utils.Encode([]byte(msg)))
			if err != nil {
				log.Println("Seem to send msg fail:", err)
			}
			log.Println("Socket Server send Alexa event over: n =", n)
		}
	}()

	for {
		n, err := conn.Read(buffer[:(len(buffer)-1)])
		if err != nil {
			utils.Log(conn.RemoteAddr().String(), " connection error:", err)
			//connOK = false
			overMsg <- true
			return
		}
		tmpBuffer = utils.Depack(append(tmpBuffer, buffer[:n]...))
		utils.TaskDeliver(tmpBuffer,conn)

		//start heartbeating
		go utils.HeartBeating(conn,messnager,timeout)
		//check if get message from client
		go utils.GravelChannel(tmpBuffer,messnager)

	}
}

