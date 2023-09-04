package blockchain

import (
	"fmt"
	"sync"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5 // 5 블록마다 걸린 시간을 측정할것임
	blockInterval      int = 2 // 2분마다 한개 생성하는것을 목표로 잡음
	allowedRange       int = 2 // expectedTime과의 Gap차이 허용 구간
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once // sync 패키지

func (b *blockchain) restore(data []byte) { // ToBytes를 통해 byte화 된걸 다시 되돌림
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
	fmt.Println(b)
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash // hashCursor: 우리가 찾을 target hash
	for {
		block, _ := FindBlock(hashCursor) // prevHash를 찾다보면 제네시스 블록은 무조건 찾을 수 있으므로 err은 무시 // 제네시스에는 prevHash가 없기때문
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash // target hash를 prevHash로 변경함으로써 계속 파고들수있음
		} else {
			break
		}
	}
	return blocks // 모든 블록이 담긴 slice를 반환
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]                                                         // chain.go에서 Blocks를 보면, 우리는 최근 해시부터 찾아들어갔다는 것을 확인할 수 있다. 즉 0번 인덱스를 조회해야 최근 블록내용이 나온다.
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]                            // 가장 최근 업데이트된 블록
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60) // unix라서 60을 나눠줌으로 분단위
	expectedTime := difficultyInterval * blockInterval                                  // 우린 블록당 2분으로 예상을 했고, 5블록마다 측정한다면 이 둘의 곱은 10분이어야함.

	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 1 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 { // 우리는 매 5블록마다 걸린 시간 측정
		// 비트코인은 2016 블록마다 측정 -> 2주간 측정했을때 24 * 14 (2주시간) == 2016 / 60 (한시간마다 1블록이라고 치면)
		// 즉 2주보다 더 걸렸으면 난이도를 낮추고, 덜 걸렸으면 난이도를 높임.

		return b.recalculateDifficulty()
	} else { // 첫번째 블록이 아니면서, 난이도 조절이 필요 없을때
		return b.CurrentDifficulty
	}

}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			} // 새로 만든 텅빈 블록체인
			checkpoint := db.Checkpoint()
			// search for checkpoint on the db
			if checkpoint == nil {
				b.AddBlock()
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
