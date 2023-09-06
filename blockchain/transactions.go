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
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx { // 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 반환
	txIns := []*TxIn{
		{"COINBASE", minerReward}, // 소유주는 채굴자
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount { // 잔금이 보내고 싶은 금액보다 적다면
		return nil, errors.New("not enough money")
	}
}

func (m *mempool) AddTx(to string, amount int) error { // mempool에 트랜잭션을 추가
	tx, err := makeTx("vivid", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}
