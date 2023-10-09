package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	v map[string]*peer // value
	m sync.Mutex       // data race를 막는 방법 // mutex는 mutex가 위치한 struct를 잠금
	// 바로 위의 map을 보호하기 위해 mutex가 있는 struct 내부에 넣고 진행
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte // 각각의 peer마다 bytes조각들을 보내는 inbox라는 채널을 줌 // 이제 특정 함수에 국한된 것이 아닌 언제어디서나 보낼 수 있음
}

func AllPeers(p *peers) []string {
	p.m.Lock()         // mutex로 잠그고
	defer p.m.Unlock() // 다 끝나면 언락
	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

func (p *peer) close() {
	Peers.m.Lock()         // mutex를 이용해서 변수를 잠금 -> unlock할때까지 접근 불가
	defer Peers.m.Unlock() // 함수 다 끝나고 unlock
	p.conn.Close()
	delete(Peers.v, p.key) // golang map 내용 삭제방법
}

func (p *peer) read() {
	// 에러 발생 시 peer 제거
	defer p.close() // defer: 함수 종료후 실행되는 코드라인
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) // websocket에서 오는 메세지를 받아서, JSON으로  변환 후, Json으로부터 go로 unmarshal 하게 도와줌
		// Message의 형식처럼 Kind, Payload로 쪼개져서 저장됨
		if err != nil {
			break
		}
		handleMsg(&m, p)
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
	Peers.v[key] = p
	return p
}
