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

// peer들이 들어갈 map
var Peers peers = peers{
	v: make(map[string]*peer),
}

// peer에 대한 구조체
type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte // 각각의 peer마다 bytes조각들을 보내는 inbox라는 채널을 줌. channel이므로 특정상황에 국한받지 않음
}

// 현재 연결된 peer들의 리스트 반환
func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

// 연결되어 있던 peer가 예기치 못한 오류로 종료될 시, 우리쪽의 peer 목록에서 삭제
func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key) // golang map 내용 삭제방법
}

// peer에게서 온 메세지 후처리
func (p *peer) read() {
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) // websocket에서 오는 메세지를 받아서, JSON으로  변환 후, Json으로부터 go로 unmarshal 하게 도와줌 (Message의 형식처럼 Kind, Payload로 쪼개져서 저장됨)
		if err != nil {
			break
		}
		handleMsg(&m, p)
	}
}

// peer들에게 메세지 작성
func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox // ok: 채널의 상태가 괜찮은지
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

// peer가 추가될 시, peer 리스트에 추가
func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock() // Peers를 조회하거나 수정할경우 data race가 발생할 수 있는데, 이를 방지하고자 mutex로 잠금 및 잠금해제
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}
	go p.read() // peer로부터 msg를 읽어오는 go 루틴 (끊기지 않고, 다른 코드를 block하지 않고)
	go p.write()
	Peers.v[key] = p
	return p
}
