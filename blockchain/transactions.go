package blockchain

import (
	"errors"
	"time"

	"github.com/viviviviviid/go-coin/utils"
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

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txId"` // TxID와 Index는, 어떤 트랜잭션이 지금 input을 생성한 output을 가지고 있는지 알려줌
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(uTxOut *UTxOut) bool {
	// mempool에 있는 트랜잭션의 input중에, uTxOut와 같은 트랜잭션 ID와 index를 가지고있는 항목이 있는지 검사
	exists := false
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			exists = input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
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
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0 // UTXO의 잔액 저장할 곳
	uTxOuts := Blockchain().UTxOutsByAddress(from)
	for _, uTxOut := range uTxOuts {
		if total > amount {
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
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error { // mempool에 트랜잭션을 추가
	tx, err := makeTx("vivid", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("vivid")
	txs := m.Txs // 블록당 트랜잭션 포함 수가 정해져있지않고, 매번 mempool에 있는 tx들을 전부 가져옴
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
