package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

type MessageKind int

const (
	MessageNewestBlock       MessageKind = iota // StatusOK = 200 과 같은 스테이터스 변수와 같은 형식스로 진행
	MessageAllBlocksRequest                     // iota 밑에 있어서, 변수들의 숫자가 0부터 1씩 증가하는 형태로 가지게 될것이고
	MessageAllBlocksResponse                    // iota 밑에 있어서, 변수들의 종류도 MessageKind가 될것
)

type Message struct { // 다른 언어와 소통하기에도 적합한 메세지 형식 정의
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
} // 이중으로 JSON화 하는 이유? : Payload의 타입이 []byte라서
// Payload안에는 block을 포함한 다양한 내용이 들어갈 수 있기 때문에, 범용성을 위해 []byte로 지정

func sendNewestBlock(p *peer) {
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, block) // JSON 바이트로 인코딩 된 메세지 반환
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) { // 들어오는 메세지의 유형에 따라 어떻게 처리할지 분류 및 처리
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
	}
}
