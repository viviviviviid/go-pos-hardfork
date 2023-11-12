package blockchain

import (
	"errors"
	"fmt"
	"time"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

type RoleInfo struct {
	ProposalAddress         string   `json:"proposalAddress"`
	ProposalPort            string   `json:"proposalPort"`
	ProposalSelectedHeight  int      `json:"proposalSelectedHeight"`
	ValidatorAddress        []string `json:"validatorAddress"`
	ValidatorPort           []string `json:"validatorPort"`
	ValidatorSelectedHeight int      `json:"validatorSelectedHeight"`
}

type Block struct {
	Hash        string               `json:"hash"`
	PrevHash    string               `json:"prevHash,omitempty"` // omitempty option
	Height      int                  `json:"height"`
	Timestamp   int                  `json:"timestamp"`
	Transaction []*Tx                `json:"transaction"`
	RoleInfo    *RoleInfo            `json:"roleinfo"`
	Signature   []*ValidateSignature `json:"signature"`
}

type ValidateSignature struct {
	Port      string `json:"port"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type ValidatedInfo struct {
	ProposalPort  string             `json:"proposalPort"`
	ProposalBlock *Block             `json:"proposalBlock"`
	Port          string             `json:"port"`
	Result        bool               `json:"result"`
	Signature     *ValidateSignature `json:"signature"`
}

func PersistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b)) // interface로 인자를 받은 ToBytes는 뭐든 받을 수 있다 = interface
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

// sign 메서드는 블록에 대해 서명을 저장합니다.
func BlockSign(b *Block, port string) *ValidateSignature {
	sig := &ValidateSignature{
		Port:      port,
		Address:   wallet.Wallet(port).Address,
		Signature: wallet.Sign(b.Hash, wallet.Wallet(port)), // 블록 id에 서명 // b.ID는 Block struct를 해쉬화한 값
	}
	return sig
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.FindBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{} // 빈 struct 만들고
	block.restore(blockBytes)
	return block, nil
}

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
	block.Timestamp = 1231006505 // bitcoin genesis block's timestamp
	block.RoleInfo = roleInfo
	block.Hash = utils.Hash(b)
	PersistBlock(block)
	return block
}

// 블록의 유효성을 검증하는 함수
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
