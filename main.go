package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

type URL string // string 형태를 가진 URL이라는 type // type을 만들 수 있음

func (u URL) MarshalText() ([]byte, error) { // MarshalText: Field가 json string으로써 어떻게 보여질지 결정하는 Method
	return []byte("hello"), nil // web화면에 json을 뿌려주면 url: "hello" 가 출력
}

func (u URL) String() string {
	return "hi"
}

type URLDescription struct {
	URL         URL    `json:"url"` // json형태로 웹에 출력된다면, 별명상태로 출력 -> 소문자로 출력시키는 방법
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload, omitempty"` // omitempty 옵션은 내용이 없을때 화면에서 생략
}

func (u URLDescription) String() string { // stringer interface는 이렇게 구현해놓은순간부터, URLDescription을 직접 print할경우 return의 내용을 출력해준다.
	return "Hello I'm the URL description" // 어떻게 변수를 넣어야할지 알려주는 가이드라인으로 작성
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "POST",
			Description: "Add Block",
			Payload:     "data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json") // json으로 인지하도록 설정

	// b, err := json.Marshal(data)                        // struct 데이터를 json으로 변환 -> 하지만 byte slice 또는 error로 return됨
	// utils.HandleErr(err)                                // 직접 만든 에러 처리 메서드
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data) // 윗 세줄과 같은 코드
}

func main() {
	// explorer.Start()
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
