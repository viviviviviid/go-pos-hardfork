package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

type Block struct {
	Data     string `json:"data"` // struct tag
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"` // omitempty option
	Height   int    `json:"height"`
}

func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer            // bytes의 Buffer는 bytes를 넣을 수 있는 공간 // read-write 가능
	encoder := gob.NewEncoder(&blockBuffer) // encoder을 만들고
	utils.HandleErr(encoder.Encode(b))      // encode해서 blockBuffer에 넣음
	return blockBuffer.Bytes()
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
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
	block.persist()
	return block
}
