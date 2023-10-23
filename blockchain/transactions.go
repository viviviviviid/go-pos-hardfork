package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

// minerReward는 채굴자에게 주어지는 보상입니다.
const (
	minerReward          int    = 50
	genesisBlockRewarder string = "6308e20ddaeae91a48a7e07791d5dabb814bae4a1e44595b0253c6051dc1c260cc6d0747370172c0db48aec400f0dbf7badbeada4f585ecd7ef5115e1dddd433"
)

// mempool은 대기 중인 트랜잭션들을 저장합니다.
type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

// 비어있는 mempool을 생성
var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
	InputData string   `json:"inputData"`
}

type TxIn struct {
	TxID      string `json:"txId"` // TxID와 Index는, 어떤 트랜잭션이 지금 input을 생성한 output을 가지고 있는지 알려줌
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// UTxOut는 사용되지 않은 트랜잭션 출력을 나타냅니다.
type UTxOut struct {
	TxID      string
	Index     int
	Amount    int
	InputData string
}

// getId 메서드는 트랜잭션 ID를 설정합니다. 트랜잭션 struct를 해시화 한걸 id에 삽입
func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

// sign 메서드는 모든 트랜잭션 입력에 대해 서명을 저장합니다.
func (t *Tx) sign() {
	for _, txIn := range t.TxIns { // 이 트랜잭션의 모든 트랜잭션 input들에 대해 서명을 저장
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet()) // 트랜잭션 id에 서명 // t.ID는 Tx struct를 해쉬화한 값
	}
}

// validate 함수는 트랜잭션의 유효성을 검증합니다.
// 트랜잭션 만든 사람을 검증 // 즉 transaction output을 소유한 사람을 검증
// output으로 트잭을 만들 수 있기 때문 -> 왜냐면 output이 다음 트잭의 input이라서
func validate(tx *Tx) bool { // 그래서 output을 보유 중인지 검증해야함
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID) // 여기에서 txIn.TxID는 지금 input으로 쓰이는 output을 만든 트잭. 즉 지금 트잭을 만들어준 이전 트잭
		if prevTx == nil {                        // 이전 트잭이 없다면, 이걸 생성한 사람은 우리 체인의 코인을 갖고있지 않다는 뜻
			valid = false // 즉 유효하지 않아서 loop 탈출
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address) // address로 publicKey를 복구할 수 있기 때문
		if !valid {
			break
		}
	}
	return valid
}

// isOnMempool 함수는 uTxOut가 mempool에 있는지 확인합니다.
func isOnMempool(uTxOut *UTxOut) bool {
	// mempool에 있는 트랜잭션의 input중에, uTxOut와 같은 트랜잭션 ID와 index를 가지고있는 항목이 있는지 검사
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer // 중첩 반복문을 전부 나가고 싶을때, label을 선언해놓고 label을 break 하면 됨
			}
		}
	}
	return exists
}

// 블록 채굴 시
// 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 반환
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // 소유주는 채굴자
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: "From chain",
	}
	tx.getId()
	return &tx
}

func makeGenesisTx() *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // 소유주는 채굴자
	}
	txOuts := []*TxOut{
		{genesisBlockRewarder, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: 1231006505,
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: "Genesis Block",
	}
	tx.getId()
	return &tx
}

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

// makeTx 함수는 일반 트랜잭션을 생성합니다.
func makeTx(from, to string, amount int, inputData string) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0 // UTXO의 잔액 저장할 곳
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 { // change: 거스름돈 // change가 0이 아니라면 거슬러줘야함
		changeTxOut := &TxOut{from, change} // 거스름돈 반환
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount} // 받을사람으르 위한 트랜잭션 output
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: inputData,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

// AddTx 메서드는 mempool에 트랜잭션을 추가
func (m *mempool) AddTx(to string, amount int, inputData string) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount, inputData)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

// TxToConfirm 메서드는 확인할 트랜잭션들을 반환
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx) // 빈 map // nil을 넣으면 삭제하는 것과 같아서
	return txs
}

func (m *mempool) GenesisTxToConfirm() []*Tx {
	coinbase := makeGenesisTx()
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx) // 빈 map // nil을 넣으면 삭제하는 것과 같아서
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs[tx.ID] = tx
}
