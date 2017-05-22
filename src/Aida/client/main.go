/*
Aida: Artificial Intelligent Digital Assistant
 */

package main

import (
	"runtime"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("=====================")
	fmt.Println("Hello, Aida")

	sendInputEvent()
	//StartService()
	//defer StopService()
	fmt.Println("=====================")
}
