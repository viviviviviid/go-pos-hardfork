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

// 포인터 주소로부터 값을 가져와 문자열로 표시
func ToString(i interface{}) string {
	r, err := json.MarshalIndent(i, "", "    ")
	HandleErr(err)
	return string(r)
}

func FormatTimeFromSeconds(sec int) string {
	if sec < 0 {
		sec = sec * (-1)
	}
	days := sec / 86400
	hours := (sec % 86400) / 3600
	minutes := (sec % 3600) / 60
	seconds := sec % 60
	return fmt.Sprintf("스테이킹 언락까지 %d일 %d시간 %d분 %d초 남았습니다. 인출은 그 이후에 가능합니다.", days, hours, minutes, seconds)
}
