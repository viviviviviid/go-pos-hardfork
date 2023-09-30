package p2p

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/viviviviviid/go-coin/utils"
)

var conns []*websocket.Conn
var upgrader = websocket.Upgrader{}

// Upgrade: 프로토콜간의 전환 // 즉 우리는 HTTP에서 WebSocket 통신으로 전환할것임.
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { // 웹소켓 연결 허가
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn) // 새로운 브라우저로 열리면 그 connection을 connections에 넣어둠
	utils.HandleErr(err)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		for _, aConn := range conns { // aConn은 Chrome, FireFox등 브라우저들에 대한 connection
			if aConn != conn { // 현재 열려있는 브라우저가 아닌 브라우저의 커넥션일 경우만 메세지를 보냄
				utils.HandleErr(aConn.WriteMessage(websocket.TextMessage, payload))
			}

		}

	}
}
