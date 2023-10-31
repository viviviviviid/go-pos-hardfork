package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

type MessageKind int

const (
	MessageNewestBlock       MessageKind = iota // StatusOK = 200 과 같은 스테이터스 변수와 같은 형식스로 진행
	MessageAllBlocksRequest                     // iota 밑에 있어서, 변수들의 숫자가 0부터 1씩 증가하는 형태로 가지게 될것이고
	MessageAllBlocksResponse                    // iota 밑에 있어서, 변수들의 종류도 MessageKind가 될것
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
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
// 그래서 Unmarshal도 두번해줘야함

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, block) // JSON 바이트로 인코딩 된 메세지 반환
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) { // 들어오는 메세지의 유형에 따라 어떻게 처리할지 분류 및 처리
	switch m.Kind {
	case MessageNewestBlock: // 3000번 입장에서 4000번으로부터의 메세지를 받고 있는 상황
		fmt.Printf("Received the newest block from %s\n", p.key)
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height > b.Height { // 우리 노드의 최신블록보다 블록높이가 높은지 확인 -> 뒤처지는지 앞서는지
			fmt.Printf("Request all block from %s\n", p.key)
			// 4000번에게 블록전체를 요청
			requestAllBlocks(p)
		} else if payload.Height < b.Height {
			fmt.Printf("Sending newest block from %s\n", p.key)
			// 4000번에게 우리의 블록들을 전달
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}
