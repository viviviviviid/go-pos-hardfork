// utils 패키지에는 애플리케이션 전반에서 사용할 함수들이 포함되어 있습니다.
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

const StakingAddress = "c8546a75af42fd63669afa3d2e72b3567790aa8f2a54da1abb94ec03239c76638f45ada90e6e2a5af42efff001a66d90106fa898ae55d3168b11d9e120a0763d"

// 에러 핸들링
func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

// bytes로 변환
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer            // bytes의 Buffer는 bytes를 넣을 수 있는 공간 // read-write 가능
	encoder := gob.NewEncoder(&aBuffer) // encoder을 만들고
	HandleErr(encoder.Encode(i))        // encode해서 blockBuffer에 넣음
	return aBuffer.Bytes()
}

// 인터페이스와 데이터를 가져와 데이터를 해당 인터페이스로 디코딩 후 저장
func FromBytes(i interface{}, data []byte) { // ex (interface{}: 블록의 포인터, data: data) -> data를 포인터로 복원
	encoder := gob.NewDecoder(bytes.NewReader(data)) // 디코더 생성
	HandleErr(encoder.Decode(i))
}

// 인터페이스를 가져와 해당 내용을 해싱한 후 해시의 16진수 인코딩을 반환
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i) // v: default formmater
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// 문자열을 원하는 주문에 맞춰 쪼개주는 함수
func Splitter(s string, sep string, index int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < index {
		return ""
	}
	return r[index]
}

// 인터페이스를 JSON으로 변환
func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}

// 인터페이스를 문자열로 변환
func ToString(i interface{}) string {
	r, err := json.MarshalIndent(i, "", "    ")
	HandleErr(err)
	return string(r)
}

// 초단위를 일/시/분/초로 변환
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

// 문자열 슬라이스끼리 비교
func CompareStringSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
