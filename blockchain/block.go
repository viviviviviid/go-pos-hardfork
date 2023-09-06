package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

type Block struct {
	Hash        string `json:"hash"`
	PrevHash    string `json:"prevHash,omitempty"` // omitempty option
	Height      int    `json:"height"`
	Difficulty  int    `json:"difficulty"`
	Nonce       int    `json:"nonce"`
	Timestamp   int    `json:"timestamp"`
	Transaction []*Tx  `json:"transaction"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b)) // interface로 인자를 받은 ToBytes는 뭐든 받을 수 있다 = interface
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{} // 빈 struct 만들고
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("\n\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
		Nonce:      0,
	}
	block.mine()
	block.Transaction = Mempool.TxToConfirm()
	// 위 block에 바로 안 넣은 이유 : 바로 윗줄 채굴이 종료되고나서 컨펌되어야하기때문
	block.persist()
	return block
}
