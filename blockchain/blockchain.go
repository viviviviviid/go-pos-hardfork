package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"` // struct tag
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"` // omitempty option
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once // sync 패키지

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash // 마지막 블록 해시 반환
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() { // 한 번만 실행함
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})

	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return GetBlockchain().blocks
}

func (b *blockchain) GetBlock(height int) *Block {
	return b.blocks[height-1] // zero-index, 0부터 셈
}
