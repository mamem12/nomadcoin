package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(v any) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(v))
	return aBuffer.Bytes()
}

func FromBytes(v any, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(v))
}

func Hash(v any) string {
	s := fmt.Sprintf("%v", v)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
