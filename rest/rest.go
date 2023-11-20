// rest 패키지는 애플리케이션에서 사용되는 REST API와 관련된 기능을 제공합니다.
package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

var port string

var (
	ResNotStaked = map[string]string{
		"message": "Not staked.",
	}
	ResTimeRemained = map[string]string{
		"message": "Staking Time is remained.",
	}
	ResUnstakeDone = map[string]string{
		"message": "Unstaking Transaction is added to mempool.",
	}
	ResStakeDone = map[string]string{
		"message": "Staking Transaction is added to mempool.",
	}
)

const (
	stakingAddress   = "c8546a75af42fd63669afa3d2e72b3567790aa8f2a54da1abb94ec03239c76638f45ada90e6e2a5af42efff001a66d90106fa898ae55d3168b11d9e120a0763d"
	stakingAmount    = 100
	stakingNodePort  = "3000" // PoS 스테이킹 풀 제공자 노드
	unstakingMessage = "unstaked"
)

type url string

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
	Port    string `json:"port"`
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To        string `json:"to"`
	Amount    int    `json:"amount"`
	InputData string `json:"inputData"`
}

type checkTx struct {
	URL         url            `json:"url"`
	Description string         `json:"description"`
	Tx          *blockchain.Tx `json:"tx"`
}

type addPeerPayload struct {
	Address, Port string
}

// 노드 실행 후, localhost:4000에 들어가면 나오는 REST API 가이드 (설명을 읽고 원하는 기능의 URL을 클릭한다)
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/balance"),
			Method:      "GET",
			Description: "See My Wallet's Balance",
		},
		{
			URL:         url("/wallet"),
			Method:      "GET",
			Description: "See My Wallet's Address",
		},
		{
			URL:         url("/peer"),
			Method:      "GET",
			Description: "See All Peer",
		},
		{
			URL:         url("/staking"),
			Method:      "GET",
			Description: "See All Staking Member",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "See All Mempool",
		},
		{
			URL:         url("/randomTransaction"),
			Method:      "POST",
			Description: "Add a Random Transaction to the Mempool",
		},
	}
	utils.HandleErr(json.NewEncoder(rw).Encode(data))
}

// 클라이언트에게 json 응답 (라우터들이 사용할 미들웨어, 해당 라우터의 핸들러 함수가 실행되기 전에 실행)
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json") // json으로 인지하도록 설정
		next.ServeHTTP(rw, r)
	})
}

// 요청된 HTTP url을 출력 (라우터들이 사용할 미들웨어, 해당 라우터의 핸들러 함수가 실행되기 전에 실행)
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

// (/blocks) GET: 전체 블록 조회, POST: 새로운 블록 추가
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock(port[1:], nil)
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

// (/blocks/{hash:[a-f0-9]+}) 특정 블록 조회
func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	// error handling
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		utils.HandleErr(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	} else {
		utils.HandleErr(encoder.Encode(block))
	}
}

// (/status) 체인의 현 상태 확인
func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

// (/balances/{address}) 특정 지갑주소의 잔액을 확인. true가 포함되지 않았다면 잔액에 해당되는 UTXO가 분리되어서 반환
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

// (/balance) 요청한 노드의 지갑 잔액을 확인
func myBalance(rw http.ResponseWriter, r *http.Request) {
	amount := blockchain.BalanceByAddress(wallet.Wallet(port[1:]).Address, blockchain.Blockchain())
	json.NewEncoder(rw).Encode(BalanceResponse{wallet.Wallet(port[1:]).Address, amount})
}

// (/mempool) 멤풀에 속해있는 트랜잭션을 확인
func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

// (/transaction) 트랜잭션 정보를 받아 구성 완료 후 멤풀에 추가
func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount, payload.InputData, port[1:])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

// (/randomTransaction) 랜덤한 트랜잭션을 멤풀에 추가. 테스트용
func randomTransaction(rw http.ResponseWriter, r *http.Request) {
	var randomQuantity = rand.Intn(50)
	var randomTo = utils.Hash(randomQuantity)
	tx, err := blockchain.Mempool().AddTx(randomTo, randomQuantity, "Random Transaction", port[1:])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(&checkTx{
		URL:         url("/mempool"),
		Description: "Check Mempool before Block Confirm. If you missed this chance, then you can find the Transaction on /blocks",
		Tx:          tx,
	},
	)
}

// (/wallet) 요청한 노드의 지갑 정보 확인
func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet(port[1:]).Address
	json.NewEncoder(rw).Encode(myWalletResponse{Port: port[1:], Address: address})
}

// (/peer) GET: 연결된 peer 노드 리스트 출력, POST: 새로운 노드와 peer 맺기
func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

// (/stake) PoS에 참여하기 위해, 스테이킹 트랜잭션을 멤풀에 추가 - (수량: 100, 기간: 1달)
func stake(rw http.ResponseWriter, r *http.Request) {
	tx, err := blockchain.Mempool().AddTx(stakingAddress, stakingAmount, port[1:], port[1:])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
	utils.HandleErr(json.NewEncoder(rw).Encode(ResStakeDone))
}

// (/unstake) 언스테이킹 트랜잭션을 멤풀에 추가 (락업 기한이 지났다면, 스테이킹 노드에게서 인출 요청)
func unstake(rw http.ResponseWriter, r *http.Request) {
	myAddress := wallet.Wallet(port[1:]).Address
	_, stakingWalletTx, indexes := blockchain.UTxOutsByStakingAddress(stakingAddress, blockchain.Blockchain())
	stakingInfoList := blockchain.GetStakingList(stakingWalletTx, blockchain.Blockchain())
	myStakingInfo := blockchain.CheckStaking(stakingInfoList, myAddress)

	if myStakingInfo == nil {
		json.NewEncoder(rw).Encode(ResNotStaked)
		return
	}

	ok, remainTime := blockchain.CheckLockupPeriod(myStakingInfo.TimeStamp)
	if !ok {
		ResTimeRemained["message"] = utils.FormatTimeFromSeconds(remainTime)
		utils.HandleErr(json.NewEncoder(rw).Encode(ResTimeRemained))
		return
	}

	tx, err := blockchain.Mempool().AddTxFromStakingAddress(stakingAddress, myAddress, "unstaking ordered", stakingNodePort, stakingAmount, myStakingInfo, indexes)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
	utils.HandleErr(json.NewEncoder(rw).Encode(ResUnstakeDone))
}

// (/staking) 현재 스테이킹 중인 노드와 지갑정보를 조회
func checkStaking(rw http.ResponseWriter, r *http.Request) {
	_, stakingWalletTx, _ := blockchain.UTxOutsByStakingAddress(stakingAddress, blockchain.Blockchain())
	stakingInfoList := blockchain.GetStakingList(stakingWalletTx, blockchain.Blockchain())
	utils.HandleErr(json.NewEncoder(rw).Encode(stakingInfoList))
}

// 라우터를 초기화하고 HTTP 서버를 시작
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()                               // Gorilla Dependecy
	router.Use(jsonContentTypeMiddleware, loggerMiddleware) // 모든 라우터가 이 middleware사용
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET") // hash: hexadecimal 타입 // [a-f0-9] 이렇게해야 둘다 받을 수 있음
	router.HandleFunc("/balance", myBalance).Methods("GET")
	router.HandleFunc("/balances/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/randomTransaction", randomTransaction).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peer", peers).Methods("GET", "POST")
	router.HandleFunc("/stake", stake).Methods("POST")
	router.HandleFunc("/unstake", unstake).Methods("POST")
	router.HandleFunc("/staking", checkStaking).Methods("GET")
	// Gorilla Mux 공식문서에 나와있는대로
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
