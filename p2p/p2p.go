package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/viviviviviid/go-coin/utils"
)

var upgrader = websocket.Upgrader{}

// Upgrade: 프로토콜간의 전환 // 즉 우리는 HTTP에서 WebSocket 통신으로 전환할것임.
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port :3000이 :4000에서 온 request를 upgrade 함
	upgrader.CheckOrigin = func(r *http.Request) bool { // 웹소켓 연결 허가
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil) // 3000번에 저장
	utils.HandleErr(err)
	initPeer(conn, "temp", "temp")
}

func AddPeer(address, port string) { // 서로간에 connection생성, port가 node라고 생각.
	// Port 4000번이 3000으로 upgrade를 요청 // 위 upgrade가 발생하면 우리와 3000번간의 연결이 생성
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws", address, port), nil) // 새로운 URL을 call하면 새로운 connection을 생성 -> 전화기의 다이얼 역할
	// 4000번에 저장
	utils.HandleErr(err)
	initPeer(conn, address, port)
}

// Upgrader은 3000번에 저장되는 conn(4000)
// AddPeer은 4000번에 저장되는 conn(3000)
