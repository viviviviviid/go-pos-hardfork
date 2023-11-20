package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/viviviviviid/go-coin/utils"
	"github.com/viviviviviid/go-coin/wallet"
)

const (
	proposalReward  int = 50      // 제안자 보상
	validatorReward int = 10      // 검증자 보상
	MonthToSec      int = 2592000 // 1달의 초단위 변환
	WeekToSec       int = 604800  // 1주의 초단위 변환
	DayToSec        int = 86400   // 1일의 초단위 변환
	SlotSec         int = 12      // 슬롯의 초단위 변환
)

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

// 대기 중인 트랜잭션들을 저장
func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

// 트랜잭션에 대한 구조체
type Tx struct {
	ID        string   `json:"id"`        // 트랜잭션의 해시 값
	Timestamp int      `json:"timestamp"` // 트랜잭션의 타임스탬프
	TxIns     []*TxIn  `json:"txIns"`     // 트랜잭션 Input
	TxOuts    []*TxOut `json:"txOuts"`    // 트랜잭션 Outputs
	InputData string   `json:"inputData"` // 트랜잭션에 추가적으로 기입한 문자열
}

// 트랜잭션 Input에 대한 구조체
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

// 트랜잭션 Output에 대한 구조체
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// UTXO에 대한 구조체 (사용하지 않은 TxOut)
type UTxOut struct {
	TxID      string
	Index     int
	Amount    int
	InputData string
}

// 트랜잭션 내용을 해시화 한 뒤 ID에 저장
func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

// 트랜잭션 Input에 서명 저장
func (t *Tx) sign(port string) {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet(port))
	}
}

// unstaking 시, PoS 스테이킹 풀 제공자 노드의 대리서명을 이용하여 스테이킹 자금 인출
func (t *Tx) delegateSign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.DelegateSign(t.ID)
	}
}

// 트랜잭션의 유효성을 검증: UTXO로 구성된 트랜잭션인가
func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

// UTXO가 mempool에 있는지 확인
func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer: // label로 중첩반복문 탈출
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

// 블록 채굴 시, 채굴자를 주소로 삼는 코인베이스 거래내역을 생성
func makeCoinbaseTx(roleInfo *RoleInfo) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // 소유주는 채굴자
	}
	txOuts := []*TxOut{
		{roleInfo.ProposerAddress, proposalReward},
		{roleInfo.ValidatorAddress[0], validatorReward},
		{roleInfo.ValidatorAddress[1], validatorReward},
		{roleInfo.ValidatorAddress[2], validatorReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: "Proof of Stake",
	}
	tx.getId()
	return &tx
}

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

// 일반 트랜잭션을 생성
func makeTx(from, to string, amount int, inputData string, port string) (*Tx, error) {
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
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: inputData,
	}
	tx.getId()
	tx.sign(port)
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

// 사용하고자 하는 UTXO가 포함된 트랜잭션 생성
func makeTxbyUTXO(from, to, inputData, mainPort string, amount int, sInfo *StakingInfo, indexes []int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn

	txIn := &TxIn{sInfo.ID, indexes[0], from}
	txIns = append(txIns, txIn)

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
		InputData: inputData,
	}
	tx.getId()
	tx.delegateSign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

// mempool에 트랜잭션을 추가
func (m *mempool) AddTx(to string, amount int, inputData string, port string) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet(port).Address, to, amount, inputData, port)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

// Unstaking을 위해, 해당 주소의 UTXO로 트랜잭션구성
func (m *mempool) AddTxFromStakingAddress(from, to, inputData, mainPort string, amount int, sInfo *StakingInfo, indexes []int) (*Tx, error) {
	tx, err := makeTxbyUTXO(from, to, inputData, mainPort, amount, sInfo, indexes)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

// 확인할 트랜잭션들을 반환
func (m *mempool) TxToConfirm(port string, roleInfo *RoleInfo) []*Tx {
	coinbase := makeCoinbaseTx(roleInfo)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

// 노드간 전파된 peer의 트랜잭션 추가
func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs[tx.ID] = tx
}

// Unstaking 시, 락업 기간이 남아있는지 확인
func CheckLockupPeriod(timeStamp int) (ok bool, gapTime int) {
	gapTime = int(time.Now().Unix()) - timeStamp
	if gapTime > MonthToSec { // 한달 락업 기간이 지났는지 확인
		return true, 0
	}
	return false, gapTime - MonthToSec
}

// 검증자가 트랜잭션 검증 시 트랜잭션 비교
func compareTransactions(txs1, txs2 []*Tx) bool {
	if len(txs1) != len(txs2) {
		return false
	}
	for i := range txs1 {
		if !compareSingleTransaction(txs1[i], txs2[i]) {
			return false
		}
	}
	return true
}

// 검증자가 트랜잭션 검증 시 세부 트랜잭션 비교
func compareSingleTransaction(tx1, tx2 *Tx) bool {
	if tx1.InputData != tx2.InputData {
		return false
	}
	if !compareTxIns(tx1.TxIns, tx2.TxIns) {
		return false
	}
	if !compareTxOuts(tx1.TxOuts, tx2.TxOuts) {
		return false
	}
	return true
}

// 검증자가 트랜잭션 검증 시 트랜잭션 Input 비교
func compareTxIns(ins1, ins2 []*TxIn) bool {
	if len(ins1) != len(ins2) {
		return false
	}
	for i, in1 := range ins1 {
		in2 := ins2[i]
		if in1.TxID != in2.TxID || in1.Index != in2.Index || in1.Signature != in2.Signature {
			return false
		}
	}
	return true
}

// 검증자가 트랜잭션 검증 시 트랜잭션 Output 비교
func compareTxOuts(outs1, outs2 []*TxOut) bool {
	if len(outs1) != len(outs2) {
		return false
	}
	for i, out1 := range outs1 {
		out2 := outs2[i]
		if out1.Address != out2.Address || out1.Amount != out2.Amount {
			return false
		}
	}
	return true
}
