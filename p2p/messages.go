package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

// 메세지 번호
type MessageKind int

// 메세지 식별자
const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
	MessageNewProposerNotify
	MessageNewValidatorNotify
	MessageValidateRequest
	MessageValidateResponse
	MessageProposalResponse
)

var validatedResults []*blockchain.ValidatedInfo

// 메세지 구조체
type Message struct {
	Kind    MessageKind
	Payload []byte
}

// peer에게 보낼 메세지 생성
func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}

// 새로 연결된 peer에게 저장된 데이터를 비교하기 위해, 최근 블록 전송
func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, block)
	p.inbox <- m
}

// 상대 peer가 더 높은 블록 높이를 가지고 있을경우, 대체하기 위해 모든 블록 요청
func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

// requestAllBlocks의 응답으로 peer에게 모든 블록 전송
func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

// 검증자에게 제안하고자 하는 블록을 보내 검증 요청
func requestValidateBlock(v *validateRequest, p *peer) {
	m := makeMessage(MessageValidateRequest, v)
	p.inbox <- m
}

// 제안자가 모든 검증과정을 거치고 블록을 추가했을때, peer들에게 새로 추가된 블록을 저장하라고 알림
func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

// 트랜잭션이 생성되었을때, peer들에게 새로 추가된 트랜잭션을 저장하라고 알림
func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

// 새로운 peer와 연결되었을때, 기존 연결되어있던 peer들에게 새로운 peer가 연결되었다고 알림
func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

// 새롭게 뽑힌 블록 제안자에게, 제안자로 선출되었다고 알림
func notifyNewProposer(roleInfo *blockchain.RoleInfo, p *peer) {
	m := makeMessage(MessageNewProposerNotify, roleInfo)
	p.inbox <- m
}

// 새롭게 뽑힌 검증자에게, 검증자로 선출되었다고 알림
func notifyNewValidator(p *peer) {
	m := makeMessage(MessageNewValidatorNotify, nil)
	p.inbox <- m
}

// PoS 스테이킹 풀 제공자 노드에게 블록 검증 결과를 알림
func notifyValidatedResult(validatedInfo *blockchain.ValidatedInfo, p *peer) {
	m := makeMessage(MessageValidateResponse, validatedInfo)
	p.inbox <- m
}

// PoS 스테이킹 풀 제공자 노드가 블록 검증결과를 종합하고 과반수를 매겨 제안자에게 제안결과를 알림
func notifyProposalResult(proposalResult *blockchain.ValidatedInfo, p *peer) {
	m := makeMessage(MessageProposalResponse, proposalResult)
	p.inbox <- m
}

// 메세지를 수신과 관련된 핸들러
func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height > b.Height { // 우리 노드의 최신블록보다 블록높이가 높은지 확인 -> 뒤처지는지 앞서는지
			fmt.Printf("Request all block from %s\n", p.key)
			requestAllBlocks(p)
		} else if payload.Height < b.Height {
			fmt.Printf("Sending newest block from %s\n", p.key)
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

	case MessageNewProposerNotify:
		var payload *blockchain.RoleInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("At %d height, this %s node has been pointed as a Proposer\n", blockchain.Blockchain().Height+1, payload.ProposerPort)
		newBlock := blockchain.CreateBlock(blockchain.Blockchain().NewestHash, blockchain.Blockchain().Height+1, payload.ProposerPort, payload, false)
		fmt.Println("Just created new block :", utils.ToString(newBlock))
		SendProposalBlock(payload, newBlock)

	case MessageNewValidatorNotify:
		fmt.Printf("At %d height, this node has been pointed as a Validator for 3 blocks\n", blockchain.Blockchain().Height+1)

	case MessageValidateRequest:
		var payload *validateRequest
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		newBlock := blockchain.CreateBlock(blockchain.Blockchain().NewestHash, blockchain.Blockchain().Height+1, payload.RoleInfo.ProposerPort, payload.RoleInfo, false)
		fmt.Println("comparisonBlock: ", utils.ToString(newBlock))
		result := blockchain.ValidateBlock(payload.RoleInfo, payload.Block, newBlock, payload.Port)
		fmt.Println("validate result: ", utils.ToString(result.Result))
		SendValidatedResult(result)

	case MessageValidateResponse:
		var payload *blockchain.ValidatedInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		validatedResults = append(validatedResults, payload)
		var validateSignature []*blockchain.ValidateSignature
		if len(validatedResults) == 3 {
			for _, v := range validatedResults {
				fmt.Print(v.Port, " ")
				validateSignature = append(validateSignature, v.Signature)
			}
			fmt.Println("\nAll Validator's message is arrived")
			result := blockchain.CalculateMajority(validatedResults)
			proposalResult := &blockchain.ValidatedInfo{
				ProposerPort:  payload.ProposerPort,
				ProposalBlock: payload.ProposalBlock,
				Port:          utils.StakingNodePort,
				Result:        result,
				Signature:     blockchain.BlockSign(payload.ProposalBlock, utils.StakingNodePort),
			}
			proposalResult.ProposalBlock.Signature = validateSignature
			SendProposalResult(proposalResult)
			validatedResults = []*blockchain.ValidatedInfo{} // 그 뒤 다음블록의 새로운 검증자들의 결과를 받기위해서 초기화
			// if proposalResult.Result == false {              // 악의적인 노드로 판명났다면, 해당 제안자의 스테이킹 자금을 소각
			// 3000번 본인포트로부터, 소각 주소로 보내 트랜잭션 구성 및 멤풀에 추가
			// 트랜잭션 구성시, 제안자 노드가 스테이킹했던 트랜잭션 UTXO를 이용해야함. (unstake 내용 확인)
			// 스테이킹 인포부터 불러와야함
			// }
		}

	case MessageProposalResponse:
		var payload *blockchain.ValidatedInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		if payload.Port == utils.StakingNodePort && payload.Result {
			blockchain.PersistBlock(payload.ProposalBlock)
			blockchain.Blockchain().UpdateBlockchain(payload.ProposalBlock)
			BroadcastNewBlock(payload.ProposalBlock)
		}
	}
}
