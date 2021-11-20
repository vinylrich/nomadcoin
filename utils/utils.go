package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Encoding(i interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	HandleError(encoder.Encode(i))
	return blockBuffer.Bytes()
}

func Decoding(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, index int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < index {
		return ""
	}
	return r[index]
}
