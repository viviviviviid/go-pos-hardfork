package blockchain

import (
	"errors"
	"time"

	"github.com/viviviviviid/go-coin/utils"
)

type Block struct {
	Hash        string `json:"hash"`
	PrevHash    string `json:"prevHash,omitempty"` // omitempty option
	Height      int    `json:"height"`
	Timestamp   int    `json:"timestamp"`
	Transaction []*Tx  `json:"transaction"`
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

func (b *Block) mine() {
	b.Timestamp = int(time.Now().Unix())
	b.Hash = utils.Hash(b)
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	block.Transaction = Mempool().TxToConfirm()
	block.mine()
	persistBlock(block)
	return block
}
