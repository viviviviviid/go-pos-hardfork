package blockchain

import (
	"errors"
	"time"

	"github.com/viviviviviid/go-coin/utils"
)

type RoleInfo struct {
	MinerAddress            string
	MinerPort               string
	MinerSelectedHeight     int
	ValidatorAddress        string
	ValidatorPort           string
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

func persistBlock(b *Block) {
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

func createBlock(prevHash string, height int, port string) *Block {
	block := &Block{
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	block.Transaction = Mempool().TxToConfirm(port)
	block.Timestamp = int(time.Now().Unix())
	block.Hash = utils.Hash(b)
	persistBlock(block)
	return block
}

func createGenesisBlock() *Block {
	// defer b.Selector()
	roleInfo := &RoleInfo{
		MinerAddress:            "Genesis",
		MinerPort:               "3000",
		MinerSelectedHeight:     1,
		ValidatorAddress:        "Genesis",
		ValidatorPort:           "3000",
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
	persistBlock(block)
	return block
}
