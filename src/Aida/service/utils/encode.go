package utils

import (
	"bytes"
	"encoding/base64"
)

func Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	encoder.Write(raw)
	encoder.Close()
	return encoded.Bytes()
}


func Decode(coded []byte) []byte {
	var buf bytes.Buffer
	decoded := make([]byte, 215)
	buf.Write(coded)
	decoder := base64.NewDecoder(base64.StdEncoding, &buf)
	decoder.Read(decoded)
	return decoded
}
