package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
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
	Payload     string `json:"payload,omitempty"` // omitempty 옵션은 내용이 없을때 화면에서 생략
}

type BalanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type addPeerPayload struct {
	Address, Port string
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
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
			URL:         url("/blocks/{hash}"),
			Method:      "POST",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}
	utils.HandleErr(json.NewEncoder(rw).Encode(data))
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // http://localhost:4000/blocks 에 들어갔을때
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
		// Encode가 Marshall의 일을 해주고, 결과를 ResponseWrite에 작성
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated) // StatusCreated : 201 (status code)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)  // r인 request에서 Mux가 변수를 추출
	hash := vars["hash"] // 윗줄에서 추출한 변수 map에서 id를 추출
	block, err := blockchain.FindBlock(hash)
	// error handling
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		utils.HandleErr(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	} else {
		utils.HandleErr(encoder.Encode(block))
	}

}

func status(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain()))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler { //
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json") // json으로 인지하도록 설정
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total") // url에 total이 있는지 확인
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(BalanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var paylaod addPeerPayload
		json.NewDecoder(r.Body).Decode(&paylaod)
		p2p.AddPeer(paylaod.Address, paylaod.Port, port)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.Peers)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()                               // Gorilla Dependecy
	router.Use(jsonContentTypeMiddleware, loggerMiddleware) // 모든 라우터가 이 middleware사용
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET") // hash: hexadecimal 타입 // [a-f0-9] 이렇게해야 둘다 받을 수 있음
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	// Gorilla Mux 공식문서에 나와있는대로
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
