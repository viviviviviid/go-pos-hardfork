package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

var port string

type url string // string 형태를 가진 URL이라는 type // type을 만들 수 있음

func (u url) MarshalText() ([]byte, error) { // MarshalText: Field가 json string으로써 어떻게 보여질지 결정하는 Method
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
} // URL type에 대한 method가 된 것

type urlDescription struct {
	URL         url    `json:"url"` // json형태로 웹에 출력된다면, 별명상태로 출력 -> 소문자로 출력시키는 방법
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload, omitempty"` // omitempty 옵션은 내용이 없을때 화면에서 생략
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (u urlDescription) String() string { // stringer interface는 이렇게 구현해놓은순간부터, URLDescription을 직접 print할경우 return의 내용을 출력해준다.
	return "Hello I'm the URL description" // 어떻게 변수를 넣어야할지 알려주는 가이드라인으로 작성
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Block",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "POST",
			Description: "See A Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json") // json으로 인지하도록 설정

	// b, err := json.Marshal(data)                        // struct 데이터를 json으로 변환 -> 하지만 byte slice 또는 error로 return됨
	// utils.HandleErr(err)                                // 직접 만든 에러 처리 메서드
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data) // 윗 세줄과 같은 코드
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // http://localhost:4000/blocks 에 들어갔을때
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		// Encode가 Marshall의 일을 해주고, 결과를 ResponseWrite에 작성
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody)) // r.Body에서 read한걸 NewDecoder가 제공해주는 reader에 넣기 // 그래서 decode하고 내용물을 addBlockBody에 넣음
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated) // StatusCreated : 201 (status code)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                     // r인 request에서 Mux가 변수를 추출
	id, err := strconv.Atoi(vars["height"]) // 윗줄에서 추출한 변수 map에서 id를 추출 // Atoi: string to integer
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(id)
	// error handling
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func Start(aPort int) {
	router := mux.NewRouter() // Gorilla Dependecy
	port = fmt.Sprintf(":%d", aPort)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET") // Gorilla Mux 공식문서에 나와있는대로
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
