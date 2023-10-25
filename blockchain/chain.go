package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

// const (
// 	blockInterval int = 2 // 2분마다 한개 생성하는것을 목표로 잡음
// 	allowedRange  int = 2 // expectedTime과의 Gap차이 허용 구간
// )

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
	m          sync.Mutex
}

type storage interface {
	FindBlock(hash string) []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	LoadChain() []byte
	DeleteAllBlocks()
}

var b *blockchain
var once sync.Once // sync 패키지
var dbStorage storage = db.DB{}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	dbStorage.SaveChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(port string) *Block {
	block := createBlock(b.NewestHash, b.Height+1, port)
	b.NewestHash = block.Hash
	b.Height = block.Height
	persistBlockchain(b)
	return block
}

func (b *blockchain) AddGenesisBlock() *Block {
	block := createGenesisBlock()
	b.NewestHash = block.Hash
	b.Height = block.Height
	persistBlockchain(b)
	return block
}

func Blocks(b *blockchain) []*Block { // struct를 변화시키지 않으므로, 메서드 형태보다는 함수형태로 선언
	b.m.Lock()
	defer b.m.Unlock()
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

func Txs(b *blockchain) []*Tx { // 모든 트랜잭션을 찾아주는 함수
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transaction...)
	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx { // 특정 트랜잭션 하나를 찾아주는 함수 // 이걸 이용해서 validate 함수 내에서 이전 트잭을 찾아낼 것임
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

// input으로 사용되지 않은 output들을 넘겨주는 함수
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut { // Unspent Tx Output
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool) // 사용한 트랜잭션 output -> map 형태
	for _, block := range Blocks(b) {   // 모든 블럭
		for _, tx := range block.Transaction { // 블럭의 모든 트랜잭션
			for _, input := range tx.TxIns { // 트랜잭션안의 input을 추적
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true // 어떤 트랜잭션이 output을 input으로 사용했는지 bool로 마킹
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok { // ok는 이 map안에 값의 유무 bool
						// input으로 사용하지 않은 트랜잭션 output
						uTxOut := &UTxOut{tx.ID, index, output.Amount, tx.InputData}
						if !isOnMempool(uTxOut) { // UTXO의 output을 확인해서, mempool에 있는지 확인
							uTxOuts = append(uTxOuts, uTxOut)
						}
						// 결론 : unspent transaction output을 생성할때는, 어떤 input에서라도 참조가 되지 않은 경우
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	// TxOutsByAddress로부터 합산된 잔액으로 만들어주는 함수
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		} // 새로 만든 텅빈 블록체인
		checkpoint := dbStorage.LoadChain()
		// search for checkpoint on the db
		if checkpoint == nil {
			b.AddGenesisBlock()
		} else { // checkpoint가 있다면
			// restore b from bytes
			// checkpoint가 있다면 새로 생성하는 것이 아닌 db로부터 블록체인을 복원
			b.restore(checkpoint) // ToBytes를 통해 byte화 된걸 다시 되돌림
		}
	})
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

func (b *blockchain) Replace(newBlocks []*Block) { // 기존 블록체인을, 노드간의 브로드캐스팅을 통해 새로 들어온 블록체인으로 교체 (ex. 내 블록의 높이가 상대의 블록높이보다 낮을때)
	b.m.Lock()
	defer b.m.Unlock()
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.NewestHash = newBlock.Hash

	persistBlockchain(b)
	persistBlock(newBlock)

	for _, tx := range newBlock.Transaction {
		_, ok := m.Txs[tx.ID] // 만약 Txs map에 이 ID를 가진 tx가 있다고 한다면, 이미 다른 노드에 의해 사용된 멤풀의 트잭이므로, 우리 노드의 멤풀에서 삭제
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}

// UTXO 형태로 만들어서 내보내기 UTXO에 TimeStamp 키 넣기
func CheckStaking(address string, targetAddress string, b *blockchain) []*Tx {
	var Txs []*Tx
	var targetAddrTxs []*Tx
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount, tx.InputData}
						if !isOnMempool(uTxOut) {
							Txs = append(Txs, tx)
						}
					}
				}
			}
		}
	}
	for _, tx := range Txs { // 이걸로 들쑤시다가 if로 targetAddress인지 확인하고, timestamp 확인하고,
		for index, txOut := range tx.TxOuts {
			if txOut.Address == targetAddress && index == 0 {
				targetAddrTxs = append(targetAddrTxs, tx)
			}
		}
	}
	fmt.Println(utils.ToString(targetAddrTxs))
	return targetAddrTxs
}
