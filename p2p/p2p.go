package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

var upgrader = websocket.Upgrader{}

// Upgrade: 프로토콜간의 전환 // 즉 우리는 HTTP에서 WebSocket 통신으로 전환할것임.
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port :3000이 :4000에서 온 request를 upgrade 함
	openPort := r.URL.Query().Get("openPort")           // 링크의 쿼리문을 가져와줌
	ip := utils.Splitter(r.RemoteAddr, ":", 0)          // RemoteAddr: 우리에게 요청을 보낸 주소를 재공 // 127.0.0.1:57039
	upgrader.CheckOrigin = func(r *http.Request) bool { // 웹소켓 연결 허가
		return openPort != "" && ip != "" // 공란으로 잘못보낸다면 업그레이드 안함
	}
	fmt.Printf("%s wants an upgrade \n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil) // 3000번에 저장
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)

}

func AddPeer(address, port, openPort string, broadcast bool) { // 서로간에 connection생성, port가 node라고 생각.
	// Port 4000번이 3000으로 upgrade를 요청 // 위 upgrade가 발생하면 우리와 3000번간의 연결이 생성
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil) // 새로운 URL을 call하면 새로운 connection을 생성 -> 전화기의 다이얼 역할
	utils.HandleErr(err)
	// 4000번에 저장
	p := initPeer(conn, address, port)
	if broadcast {
		BroadcastNewPeer(p) // 새로운 peer가 생겼다고 기존 peers에게 브로드캐스팅
		return
	}
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v { // 모든 피어들에게 전달
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}

// Upgrader은 3000번에 저장되는 conn(4000)
// AddPeer은 4000번에 저장되는 conn(3000)

// 우리 4000번이 쉬고있는 3000번에게 Upgrader로 upgrade를 요청
// 나중에 누군가가 우리의 peers에 연결하길 원할 수 있으니, openPort 변수로 열려있는 포트를 전달해줘야함
// 추후 2000번포트가 3000번 포트와 연결을 맺는데, 네트워크를 형성해야하므로 3000번은 현재 연결되어있는 peers들을 알려줌 (ex 우리 4000번)
// 그렇게 되면, 4000번에 업그레이드 요청을 하도록 만듬
// 이건 openPort 형태 "ws://%s:%s/ws?openPort=%s"

// Upgrader내의 ip는 -> 업그레이드를 요청하는 ip 또는 주소
