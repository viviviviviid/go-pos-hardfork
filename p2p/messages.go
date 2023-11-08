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
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
	MessageNewProposalNotify
	MessageNewValidatorNotify
	MessageValidateRequest
	MessageValidateResponse
	MessageProposalResponse
)

var validatedResults []*blockchain.ValidatedInfo

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, block)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func requestValidateBlock(v *validateRequest, p *peer) {
	m := makeMessage(MessageValidateRequest, v)
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

func notifyNewValidator(p *peer) {
	m := makeMessage(MessageNewValidatorNotify, nil)
	p.inbox <- m
}

func notifyValidatedResult(validatedInfo *blockchain.ValidatedInfo, p *peer) {
	m := makeMessage(MessageValidateResponse, validatedInfo)
	p.inbox <- m
}

func notifyProposalResult(proposalResult *blockchain.ValidatedInfo, p *peer) {
	m := makeMessage(MessageProposalResponse, proposalResult)
	p.inbox <- m
}

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

	case MessageNewProposalNotify:
		var payload *blockchain.RoleInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("At %d height, this %s node has been pointed as a Proposal\n", blockchain.Blockchain().Height+1, payload.ProposalPort)
		newBlock := blockchain.CreateBlock(blockchain.Blockchain().NewestHash, blockchain.Blockchain().Height+1, payload.ProposalPort, payload, false)
		fmt.Println("Just created new block :", utils.ToString(newBlock))
		SendProposalBlock(payload, newBlock)

	case MessageNewValidatorNotify:
		fmt.Printf("At %d height, this node has been pointed as a Validator for 3 blocks\n", blockchain.Blockchain().Height+1)

	case MessageValidateRequest:
		var payload *validateRequest
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		newBlock := blockchain.CreateBlock(blockchain.Blockchain().NewestHash, blockchain.Blockchain().Height+1, payload.RoleInfo.ProposalPort, payload.RoleInfo, false)
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
			fmt.Println("All Validator's message is arrived")
			result := blockchain.CalculateMajority(validatedResults)
			proposalResult := &blockchain.ValidatedInfo{
				ProposalPort:  payload.ProposalPort,
				ProposalBlock: payload.ProposalBlock,
				Port:          StakingPort,
				Result:        result,
				Signature:     blockchain.BlockSign(payload.ProposalBlock, StakingPort),
			}
			proposalResult.ProposalBlock.Signature = validateSignature
			SendProposalResult(proposalResult)
			validatedResults = []*blockchain.ValidatedInfo{} // 그 뒤 다음블록의 새로운 검증자들의 결과를 받기위해서 초기화
		}

	case MessageProposalResponse:
		var payload *blockchain.ValidatedInfo
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		if payload.Port == StakingPort && payload.Result {
			blockchain.PersistBlock(payload.ProposalBlock)
			blockchain.Blockchain().UpdateBlockchain(payload.ProposalBlock)
			BroadcastNewBlock(payload.ProposalBlock)
			fmt.Println("Added and broadcasted the block done")
		} else if payload.ProposalPort == StakingPort && !payload.Result {
			fmt.Println("Proposal Rejected")
		}
	}
}
