package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/viviviviviid/go-coin/utils"
)

const port string = ":4000"

type URLDescription struct {
	URL         string
	Method      string
	Description string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
	}
	b, err := json.Marshal(data) // struct 데이터를 json으로 변환 -> 하지만 byte slice 또는 error로 return됨
	utils.HandleErr(err)         // 직접 만든 에러 처리 메서드

	fmt.Printf("%s", b) // byte slice 형태의 json을 decode
}

func main() {
	// explorer.Start()
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
