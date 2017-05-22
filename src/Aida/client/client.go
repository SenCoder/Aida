package main

/*

#include <linux/input.h>
#include <linux/uinput.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>

int sendInputEvent(const __u16 type, const __u16 code, const __s32 value)
{

	printf("cgo alexa ===========\n");
	int uinputfd;

	struct input_event event;
	memset(&event, 0, sizeof(event));
	event.type = type;
	event.code = code;
	event.value = value;

	uinputfd = open("/dev/input/uinput", O_WRONLY | O_NDELAY);
	if (uinputfd < 0)
	{
		printf("Open dev/input/uinput error\n");
		return -1;
	}

	if (write(uinputfd, &event, sizeof(event)) != sizeof(event))
	{
		printf("Error on send input event\n");
		close(uinputfd);
		return -1;
	}
	printf("sendInputEvent return normal ===========\n");
	close(uinputfd);
	return 0;
}
 */
import "C"
import (
	"net"
	"fmt"
	"os"
	"log"
	"time"
	"strconv"
	"Aida/service/utils"
)

const (
	SERVER_ADDRESS = "192.168.1.114:1024"
)

var conn *net.TCPConn
var isStart bool

func init() {
	isStart = false
}


func sendInputEvent() {
	fmt.Println("client.sendInputEvent")
	//C.sendInputEvent(1, 0x09d00, 1)
	fmt.Println("=======================")
}

func StartService() (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", SERVER_ADDRESS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	defer conn.Close()
	log.Println("connect success")
	isStart = true

	go heartBeat()
	receiveMsg()
	return nil
}


func StopService() {
	isStart = false
}

// 发送心跳包
func heartBeat() error {
	count := 0
	for isStart {
		time.Sleep(time.Second * 5)
		conn.Write([]byte("I am on line:" + strconv.Itoa(count)))
	}
	return nil
}

// 收取 socket server 信息
func receiveMsg() error {
	buffer := make([]byte, 1024)
	for isStart {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), "connection error:", err)
			isStart = false
		} else {
			msg := string(utils.Decode(buffer))
			log.Println("Msg:", msg, n)
			err :=  msgHandle(msg)
			if err != nil {

			}
		}
	}
	return nil
}


func msgHandle(msg string) error {
	fmt.Println("msg >>", msg)
	return nil
}


func startHttp() error {

	return nil
}