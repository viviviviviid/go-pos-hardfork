package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte // 각각의 peer마다 bytes조각들을 보내는 inbox라는 채널을 줌 // 이제 특정 함수에 국한된 것이 아닌 언제어디서나 보낼 수 있음
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key) // golang map 내용 삭제방법

}

func (p *peer) read() {
	// 에러 발생 시 peer 제거
	defer p.close() // defer: 함수 종료후 실행되는 코드라인
	for {
		_, m, err := p.conn.ReadMessage() // blocking and read msg
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	defer p.close() // defer: 함수 종료후 실행되는 코드라인
	for {
		m, ok := <-p.inbox // ok: 채널의 상태가 괜찮은지
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}

	go p.read() // peer로부터 msg를 읽어오는 go 루틴 // 끊기지 않고, 다른 코드를 block하지 않고
	go p.write()
	Peers[key] = p
	return p
}
