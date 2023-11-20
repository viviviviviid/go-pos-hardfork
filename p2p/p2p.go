// p2p 패키지는 노드간의 연결, 메시지 전송, 메시지 핸들링 등의 함수를 제공합니다.
package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

var upgrader = websocket.Upgrader{}

type validateRequest struct {
	RoleInfo *blockchain.RoleInfo
	Block    *blockchain.Block
	Port     string
}

// Upgrade: 프로토콜간의 전환 (HTTP에서 WebSocket 통신으로 전환)
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")           // 링크의 쿼리문을 추출
	ip := utils.Splitter(r.RemoteAddr, ":", 0)          // RemoteAddr: 우리에게 요청을 보낸 주소를 제공
	upgrader.CheckOrigin = func(r *http.Request) bool { // 웹소켓 연결 허가
		return openPort != "" && ip != "" // 공란으로 잘못보낸다면 업그레이드 안함
	}
	fmt.Printf("%s wants an upgrade \n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}

// peer 추가
func AddPeer(address, port, openPort string, broadcast bool) { // 서로간에 connection생성, port가 node라고 생각.
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil) // 새로운 URL을 call하면 새로운 connection을 생성 -> 전화기의 다이얼 역할
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		BroadcastNewPeer(p) // 새로운 peer가 생겼다고 기존 peers에게 브로드캐스팅
		return
	}
	sendNewestBlock(p)
}

// 새로 선출된 제안자 지목
func PointingProposer(r *blockchain.RoleInfo) {
	for _, p := range Peers.v {
		if r.ProposerPort == p.port {
			notifyNewProposer(r, p)
		}
	}
}

// 새로 선출된 검증자 지목
func PointingValidator(r *blockchain.RoleInfo) {
	for _, p := range Peers.v {
		for _, port := range r.ValidatorPort {
			if port == p.port {
				notifyNewValidator(p)
			}
		}
	}
}

// 제안하고자 하는 블록을 검증자들에게 전달 후 검증 요청
func SendProposalBlock(r *blockchain.RoleInfo, b *blockchain.Block) {
	for _, p := range Peers.v {
		for _, port := range r.ValidatorPort {
			if port == p.port {
				requestFormat := &validateRequest{
					RoleInfo: r,
					Block:    b,
					Port:     p.port,
				}
				requestValidateBlock(requestFormat, p)
			}
		}
	}
}

// PoS 스테이킹 풀 제공자 노드에게 검증결과 전달
func SendValidatedResult(validatedInfo *blockchain.ValidatedInfo) {
	for _, p := range Peers.v {
		if p.port == utils.StakingNodePort { // staking port
			notifyValidatedResult(validatedInfo, p)
		}
	}
}

// PoS 스테이킹 풀 제공자 노드가 제안자에게 제안 결과 전달
func SendProposalResult(proposalResult *blockchain.ValidatedInfo) {
	for _, p := range Peers.v {
		if p.port == proposalResult.ProposerPort { // staking port
			notifyProposalResult(proposalResult, p)
		}
	}
}

// 제안자가 모든 검증결과를 마치고, 새로운 블록을 추가했을때 peer들에게 새로만든 블록을 전파
func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

// 새로 만든 트랜잭션을 peer들에게 전파
func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

// 기존 peer들에게 새로 연결된 peer의 정보를 전달
func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
