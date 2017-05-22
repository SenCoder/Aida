package utils

import (
	"testing"
	"fmt"
)

func TestEncode(t *testing.T) {
	fmt.Println(Encode([]byte("hello alexa")))
}


func TestDecode(t *testing.T) {
	fmt.Println(string(Decode(Encode([]byte("hello alexa")))))
}