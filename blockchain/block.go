package blockchain

import (
	"errors"
	"fmt"
	"time"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

type RoleInfo struct {
	ProposalAddress         string   `json:"proposalAddress"`         // 제안자의 주소
	ProposalPort            string   `json:"proposalPort"`            // 제안자의 노드 포트
	ProposalSelectedHeight  int      `json:"proposalSelectedHeight"`  // 제안자가 선정된 블록 높이
	ValidatorAddress        []string `json:"validatorAddress"`        // 검증자 주소
	ValidatorPort           []string `json:"validatorPort"`           // 검증자의 노드 포트
	ValidatorSelectedHeight int      `json:"validatorSelectedHeight"` // 검증자가 선정된 블록 높이
}

type Block struct {
	Hash        string               `json:"hash"`               // 블록의 해시 값
	PrevHash    string               `json:"prevHash,omitempty"` // 직전 블록의 해시 값
	Height      int                  `json:"height"`             // 블록 높이
	Timestamp   int                  `json:"timestamp"`          // 블록 생성 타임스탬프
	Transaction []*Tx                `json:"transaction"`        // 블록내의 트랜잭션
	RoleInfo    *RoleInfo            `json:"roleinfo"`           // 블록 추가를 위해 구성된 제안자, 검증자 정보
	Signature   []*ValidateSignature `json:"signature"`          // 블록의 유효성을 확인한 검증자들의 서명
}

type ValidateSignature struct {
	Port      string `json:"port"`      // 검증자의 노트 포트
	Address   string `json:"address"`   // 검증자의 주소
	Signature string `json:"signature"` // 검증자가 블록을 서명한 값
}

type ValidatedInfo struct {
	ProposalPort  string             `json:"proposalPort"`  // 제안자의 노드 포트
	ProposalBlock *Block             `json:"proposalBlock"` // 제안자가 제안한 블록
	Port          string             `json:"port"`          // 검증자 노드 포트
	Result        bool               `json:"result"`        // 검증 결과
	Signature     *ValidateSignature `json:"signature"`     // 검증자 서명 정보
}

var ErrNotFound = errors.New("block not found")

// 풀노드의 db에 최신 블록을 업데이트
func PersistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b))
}

// bytes 형태의 블록정보를 json으로 복구
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

// 블록 해시에 서명
func BlockSign(b *Block, port string) *ValidateSignature {
	sig := &ValidateSignature{
		Port:      port,
		Address:   wallet.Wallet(port).Address,
		Signature: wallet.Sign(b.Hash, wallet.Wallet(port)),
	}
	return sig
}

// 블록 해시로 특정 블록을 조회
func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.FindBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

// 블록 구성 함수
func CreateBlock(prevHash string, height int, port string, roleInfo *RoleInfo, update bool) *Block {
	block := &Block{
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	block.Transaction = Mempool().TxToConfirm(port, roleInfo)
	block.Timestamp = int(time.Now().Unix())
	if roleInfo == nil {
		roleInfo = &RoleInfo{
			ProposalAddress:         "Staking Address",
			ProposalPort:            "Staking Port",
			ProposalSelectedHeight:  b.Height,
			ValidatorAddress:        []string{"Staking Address", "Staking Address", "Staking Address"},
			ValidatorPort:           []string{"Staking Port", "Staking Port", "Staking Port"},
			ValidatorSelectedHeight: b.Height,
		}
	}
	block.RoleInfo = roleInfo
	block.Hash = utils.Hash(b)
	if update {
		PersistBlock(block)
	}
	return block
}

// 제네시스 블록 구성 함수
func createGenesisBlock() *Block {
	roleInfo := &RoleInfo{
		ProposalAddress:         "Genesis",
		ProposalPort:            "3000",
		ProposalSelectedHeight:  1,
		ValidatorAddress:        []string{"Genesis", "Genesis", "Genesis"},
		ValidatorPort:           []string{"3000", "3000", "3000"},
		ValidatorSelectedHeight: 1,
	}
	block := &Block{
		Hash:     "",
		PrevHash: "",
		Height:   1,
	}
	block.Transaction = Mempool().GenesisTxToConfirm()
	block.Timestamp = 1231006505 // 비트코인 제네시스 블록의 실제 타임스탬프
	block.RoleInfo = roleInfo
	block.Hash = utils.Hash(b)
	PersistBlock(block)
	return block
}

// 블록 유효성 검증
func ValidateBlock(roleInfo *RoleInfo, proposalBlock *Block, createdBlock *Block, port string) *ValidatedInfo {
	var result = true
	var sig *ValidateSignature

	if proposalBlock.PrevHash != createdBlock.PrevHash {
		fmt.Println("Not pass: prev")
		result = false
	}
	if proposalBlock.Height != createdBlock.Height {
		fmt.Println("Not pass: height")
		result = false
	}
	if !compareTransactions(proposalBlock.Transaction, createdBlock.Transaction) {
		fmt.Println("Not pass: tx")
		result = false
	}
	if !compareRoleInfo(proposalBlock.RoleInfo, createdBlock.RoleInfo) {
		fmt.Println("Not pass: roleinfo")
		result = false
	}
	if result {
		sig = BlockSign(proposalBlock, port)
	}
	v := &ValidatedInfo{
		ProposalPort:  roleInfo.ProposalPort,
		ProposalBlock: proposalBlock,
		Port:          port,
		Result:        result,
		Signature:     sig,
	}

	return v
}
