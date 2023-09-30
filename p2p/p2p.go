package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/viviviviviid/go-coin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { // 웹소켓 연결 허가
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		_, payload, err := conn.ReadMessage()
		utils.HandleErr(err)
		fmt.Printf("%s\n\n", payload)
	}
}
