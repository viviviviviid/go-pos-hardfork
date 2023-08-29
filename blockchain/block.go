package blockchain

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Data     string `json:"data"` // struct tag
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"` // omitempty option
	Height   int    `json:"height"`
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height) // string
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))    // payload를 byte slice로 hash하고 결과를 hex 형태의 string으로 받음
	return block
}
