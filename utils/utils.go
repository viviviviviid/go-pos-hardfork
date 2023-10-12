// Package utils contains functions to be used across the application.
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var logFn = log.Panic

func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

func ToBytes(i interface{}) []byte { // interface: 함수에게 뭐든 받으라고 하는 것
	var aBuffer bytes.Buffer            // bytes의 Buffer는 bytes를 넣을 수 있는 공간 // read-write 가능
	encoder := gob.NewEncoder(&aBuffer) // encoder을 만들고
	HandleErr(encoder.Encode(i))        // encode해서 blockBuffer에 넣음
	return aBuffer.Bytes()
}

// FromBytes take an interface and data and then will encode the data to the interface
func FromBytes(i interface{}, data []byte) { // ex (interface{}: 블록의 포인터, data: data) -> data를 포인터로 복원
	encoder := gob.NewDecoder(bytes.NewReader(data)) // 디코더 생성
	HandleErr(encoder.Decode(i))
}

// Hash takes an interface, hashes it and returns the hex encoding of the hash.
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i) // v: default formmater
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, index int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < index { // 원하는 인덱스보다 길이가 작으면
		return ""
	}
	return r[index]
}

func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
