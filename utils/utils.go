package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	HandleError(encoder.Encode(i))
	return blockBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}
