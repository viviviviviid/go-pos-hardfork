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
	MessageNewProposalNotify
	MessageNewValidatorNotify
	MessageValidateResponse
)

var ValidatedResults []*blockchain.ValidatedInfo

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

func notifyNewProposal(roleInfo *blockchain.RoleInfo, p *peer) {
	m := makeMessage(MessageNewProposalNotify, roleInfo)
	p.inbox <- m
}

func notifyNewValidator(container []interface{}, p *peer) {
	m := makeMessage(MessageNewValidatorNotify, container)
	p.inbox <- m
}

func notifyValidatedResult(validatedInfo *blockchain.ValidatedInfo, p *peer) {
	m := makeMessage(MessageValidateResponse, validatedInfo)
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
	case MessageNewProposalNotify:
		var payload *blockchain.RoleInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("At current height, this %s node has been pointed as a Proposal\n", payload.ProposalPort)
		fmt.Println("Starting to create block as a Proposal")
		newBlock := blockchain.Blockchain().AddBlock(payload.ProposalPort, payload)
		fmt.Println("Just created new block :", utils.ToString(newBlock))
		BroadcastNewBlock(newBlock)
		fmt.Println("Added and broadcasted the block done")
	case MessageNewValidatorNotify:
		var payload []interface{}
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		// roleInfo := payload[0].(*blockchain.RoleInfo)
		port := payload[1].(string)
		fmt.Printf("At current height, this %s node has been pointed as a Validator for 3 blocks\n", port)
		// 제안자한테 블록이 오기까지 대기
		// 제안자에게 블록이 오면, 나또한 블록 구성하기
		// 비교하고, 검증하기
		result := blockchain.ValidateBlock(nil, nil, port)
		fmt.Println(utils.ToString(result))
		// true/false를 3000번에게 보내기
		SendValidatedResult(result)
	case MessageValidateResponse:
		var payload *blockchain.ValidatedInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		ValidatedResults = append(ValidatedResults, payload)
		if len(ValidatedResults) < 3 {
			fmt.Println(utils.ToString(ValidatedResults))
		} else { // 3개가 채워질때까지
			fmt.Println("3 Validator Done: ", utils.ToString(ValidatedResults))
			// 채워졌다면 과반수가 찬성인지 반대인지 확인

			// 1. 과반수가 찬성일 경우 블록제안자에게 true를 보낸다
			// 2. true를 보내면서 블록을 하나 추가할 수 있는 권한을 부여한다. → 음 이건 struct type으로 해결해야하나?
			// 3. 과반수가 반대일 경우 블록제안자에게 false를 보낸다
			// 4. false를 보내면서 Selector로 역할을 새롭게 선정하게 하기위해, pos.go의 for문을 continue 시킨다.

			// 그 뒤 다음블록의 새로운 검증자들의 결과를 받기위해서 초기화
			ValidatedResults = []*blockchain.ValidatedInfo{}
			fmt.Println("validatedResult clear!: ", utils.ToString(ValidatedResults))
		}
	}
}
