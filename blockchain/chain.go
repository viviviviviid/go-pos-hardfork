package blockchain

import (
	"fmt"
	"sync"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once // sync 패키지

func (b *blockchain) restore(data []byte) { // ToBytes를 통해 byte화 된걸 다시 되돌림
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil { // 블록체인 최초 실행
		once.Do(func() {
			b = &blockchain{"", 0} // 새로 만든 텅빈 블록체인
			checkpoint := db.Checkpoint()
			// search for checkpoint on the db
			if checkpoint == nil {
				b.AddBlock("Genesis Block\n")
			} else { // checkpoint가 있다면
				// restore b from bytes
				// checkpoint가 있다면 새로 생성하는 것이 아닌 db로부터 블록체인을 복원
				b.restore(checkpoint) // ToBytes를 통해 byte화 된걸 다시 되돌림
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}
