package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/viviviviviid/go-coin/db"
	"github.com/viviviviviid/go-coin/utils"
)

type blockchain struct {
	NewestHash string     `json:"newestHash"` // 블록체인 중 최근 블록 해시 값
	Height     int        `json:"height"`     // 블록체인의 현 블록 높이
	m          sync.Mutex // data race를 방지하기 위한 라이브러리
}

type StakingInfo struct {
	ID        string `json:"id"`        // 스테이킹 트랜잭션의 해시 값
	Address   string `json:"address"`   // 스테이커 주소
	Port      string `json:"port"`      // 스테이커 노드 포트
	TimeStamp int    `json:"timestamp"` // 스테이킹 트랜잭션의 타임스탬프
}

type storage interface {
	FindBlock(hash string) []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	LoadChain() []byte
	DeleteAllBlocks()
}

var (
	stakingQuantity = 100
	b               *blockchain
	once            sync.Once // sync 패키지
	dbStorage       storage   = db.DB{}
)

// bytes 형태의 블록체인 값 복구
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// 풀노드의 db에 현 블록체인 상태 저장
func persistBlockchain(b *blockchain) {
	dbStorage.SaveChain(utils.ToBytes(b))
}

// 블록체인 상태 업데이트
func (b *blockchain) UpdateBlockchain(block *Block) {
	b.NewestHash = block.Hash
	b.Height = block.Height
	persistBlockchain(b)
}

// 블록체인에 블록 추가
func (b *blockchain) AddBlock(port string, roleInfo *RoleInfo) *Block {
	block := CreateBlock(b.NewestHash, b.Height+1, port, roleInfo, true)
	b.UpdateBlockchain(block)
	return block
}

// 최초 상태의 블록체인에 제네시스 블록 추가 (위 AddBlock과 구분한 이유는 비트코인의 타임스탬프 등 여러가지 조건을 넣고 싶어서)
func (b *blockchain) AddGenesisBlock() *Block {
	block := createGenesisBlock()
	b.UpdateBlockchain(block)
	return block
}

// 전체 블록 탐색 후 반환
func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash // target hash를 prevHash로 변경함으로써 계속 파고들수있음
		} else {
			break
		}
	}
	return blocks // 모든 블록이 담긴 slice 반환
}

// 전체 트랜잭션 반환
func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transaction...)
	}
	return txs
}

// 특정 트랜잭션 정보 반환
func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

// 트랜잭션의 input으로 사용되지 않은 UTXO들을 반환
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut { // Unspent Tx Output
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
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
						uTxOut := &UTxOut{tx.ID, index, output.Amount, tx.InputData}
						if !isOnMempool(uTxOut) { // UTXO의 output을 확인해서, mempool에 있는지 확인
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

// 특정 주소의 코인 잔액
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

// 풀노드 재 구동 시 저장된 블록체인 상태 불러오기
func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := dbStorage.LoadChain()

		if checkpoint == nil {
			b.AddGenesisBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}

// 블록체인 상태 반환
func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

// 노드간 브로드캐스팅을 통해, 블록 높이 비교 후 대체
func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		PersistBlock(block)
	}
}

// 노드간 새로 추가된 블록을 저장
func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.NewestHash = newBlock.Hash

	persistBlockchain(b)
	PersistBlock(newBlock)

	for _, tx := range newBlock.Transaction {
		_, ok := m.Txs[tx.ID] // 만약 Txs map에 이 ID를 가진 tx가 있다고 한다면, 이미 다른 노드에 의해 사용된 멤풀의 트잭이므로, 우리 노드의 멤풀에서 삭제
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}

// 스테이커의 스테이킹과 관련된 UTXO 반환
func UTxOutsByStakingAddress(stakingAddress string, b *blockchain) ([]*UTxOut, []*Tx, []int) {
	var uTxOuts []*UTxOut
	var Txs []*Tx
	var indexes []int

	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == stakingAddress {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == stakingAddress && output.Amount == stakingQuantity {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount, tx.InputData}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
							indexes = append(indexes, index)
							Txs = append(Txs, tx)
						}
					}
				}
			}
		}
	}
	return uTxOuts, Txs, indexes
}

// 스테이커 리스트 반환
func GetStakingList(Txs []*Tx, b *blockchain) []*StakingInfo {
	var sInfos []*StakingInfo
	var stakerAddr string
	for _, tx := range Txs {
		for _, input := range tx.TxIns {
			stakerAddr = FindTx(b, input.TxID).TxOuts[input.Index].Address
		}
		sInfo := &StakingInfo{tx.ID, stakerAddr, tx.InputData, tx.Timestamp}
		sInfos = append(sInfos, sInfo)
	}
	return sInfos
}

// 스테이킹 유무 확인
func CheckStaking(stakingInfoList []*StakingInfo, targetAddress string) *StakingInfo {
	var sInfo *StakingInfo
	for _, info := range stakingInfoList {
		if info.Address == targetAddress {
			sInfo = info
		}
	}
	return sInfo
}
