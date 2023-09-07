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
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// 비트코인은 2016 블록마다 측정 -> 2주간 측정했을때 24 * 14 (2주시간) == 2016 / 60 (한시간마다 1블록이라고 치면)
		// 즉 2주보다 더 걸렸으면 난이도를 낮추고, 덜 걸렸으면 난이도를 높임.
		return b.recalculateDifficulty()
	} else { // 첫번째 블록이 아니면서, 난이도 조절이 필요 없을때
		return b.CurrentDifficulty
	}
}

// input으로 사용되지 않은 output들을 넘겨주는 함수
func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut { // Unspent Tx Output => UTXO ㅋㅋㅋㅋㅋ 이거였네
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool) // 사용한 트랜잭션 output -> map 형태
	for _, block := range b.Blocks() {  // 모든 블럭
		for _, tx := range block.Transaction { // 블럭의 모든 트랜잭션
			for _, input := range tx.TxIns { // 트랜잭션안의 input을 추적
				if input.Owner == address {
					creatorTxs[input.TxID] = true // 어떤 트랜잭션이 output을 input으로 사용했는지 bool로 마킹
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok { // ok는 이 map안에 값의 유무 bool
						// input으로 사용하지 않은 트랜잭션 output
						uTxOuts = append(uTxOuts, &UTxOut{tx.Id, index, output.Amount})
						// 결론 : unspent transaction output을 생성할때는, 어떤 input에서라도 참조가 되지 않은 경우
					}
				}
			}
		}
	}
	return uTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	// TxOutsByAddress로부터 합산된 잔액으로 만들어주는 함수
	txOuts := b.UTxOutsByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
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
