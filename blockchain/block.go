package blockchain

import (
	"errors"
	"time"

	"github.com/viviviviviid/go-coin/utils"
)

type RoleInfo struct {
	ProposalAddress         string
	ProposalPort            string
	ProposalSelectedHeight  int
	ValidatorAddress        []string
	ValidatorPort           []string
	ValidatorSelectedHeight int
}

type Block struct {
	Hash        string `json:"hash"`
	PrevHash    string `json:"prevHash,omitempty"` // omitempty option
	Height      int    `json:"height"`
	Timestamp   int    `json:"timestamp"`
	Transaction []*Tx  `json:"transaction"`
	RoleInfo    *RoleInfo
}

type ValidatedInfo struct {
	ProposalPort string
	Port         string
	Result       bool
}

func PersistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b)) // interface로 인자를 받은 ToBytes는 뭐든 받을 수 있다 = interface
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
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

func ValidateBlock(roleInfo *RoleInfo, proposalBlock *Block, createdBlock *Block, port string) *ValidatedInfo {
	//
	//	검증프로세스
	//
	v := &ValidatedInfo{
		ProposalPort: roleInfo.ProposalPort,
		Port:         port,
		Result:       true, // 검증 프로세스 완성전까지는 true로 제공
	}
	return v
}
