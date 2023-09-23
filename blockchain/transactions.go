package blockchain

import (
	"errors"
	"time"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

// 비어있는 mempool을 생성
var Mempool *mempool = &mempool{}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
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

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns { // 이 트랜잭션의 모든 트랜잭션 input들에 대해 서명을 저장
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet()) // 트랜잭션 id에 서명 // t.ID는 Tx struct를 해쉬화한 값
	}
}

// 트랜잭션 만든 사람을 검증 // 즉 transaction output을 소유한 사람을 검증
// output으로 트잭을 만들 수 있기 때문 -> 왜냐면 output이 다음 트잭의 input이라서
func validate(tx *Tx) bool { // 그래서 output을 보유 중인지 검증해야함
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTxs(Blockchain(), txIn.TxID) // 여기에서 txIn.TxID는 지금 input으로 쓰이는 output을 만든 트잭. 즉 지금 트잭을 만들어준 이전 트잭
	}
	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	// mempool에 있는 트랜잭션의 input중에, uTxOut와 같은 트랜잭션 ID와 index를 가지고있는 항목이 있는지 검사
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer // 중첩 반복문을 전부 나가고 싶을때, label을 선언해놓고 label을 break 하면 됨
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx { // 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 반환
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
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enough money")
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
	}
	tx.getId()
	tx.sign()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error { // mempool에 트랜잭션을 추가
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs // 블록당 트랜잭션 포함 수가 정해져있지않고, 매번 mempool에 있는 tx들을 전부 가져옴
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
