package blockchain

// import (
// 	"time"

// 	"github.com/viviviviviid/go-coin/utils"
// )

// type Tx struct {
// 	Id        string
// 	Timestamp int
// 	TxIns     []*TxIn
// 	TxOuts    []*TxOut
// }

// const (
// 	minerReward int = 50
// )

// type TxIn struct {
// 	Owner  string
// 	Amount int
// }

// type TxOut struct {
// 	Owner  string
// 	Amount int
// }

// func (t *Tx) getId() {
// 	t.Id = utils.Hash(t)
// }

// func makeCoinbaseTx(address string) *Tx { // 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 반환
// 	txIns := []*TxIn{
// 		{"COINBASE", minerReward}, // 소유주는 채굴자
// 	}
// 	txOuts := []*TxOut{
// 		{address, minerReward},
// 	}
// 	tx := Tx{
// 		Id:        "",
// 		Timestamp: int(time.Now().Unix()),
// 		TxIns:     txIns,
// 		TxOuts:    txOuts,
// 	}
// 	tx.getId() // 이게 어떻게 생각해야하는거지
// 	return &tx
// }
