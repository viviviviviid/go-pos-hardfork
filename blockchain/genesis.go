package blockchain

import "github.com/viviviviviid/go-coin/utils"

// 제네시스 블록 구성 함수
func createGenesisBlock() *Block {
	roleInfo := &RoleInfo{
		ProposerAddress:         "Genesis",
		ProposerPort:            "3000",
		ProposerSelectedHeight:  1,
		ValidatorAddress:        []string{"Genesis", "Genesis", "Genesis"},
		ValidatorPort:           []string{"3000", "3000", "3000"},
		ValidatorSelectedHeight: 1,
	}
	block := &Block{
		Hash:     "",
		PrevHash: "",
		Height:   1,
	}
	block.Transaction = Mempool().GenesisTxToConfirm()
	block.Timestamp = 1231006505 // 비트코인 제네시스 블록의 실제 타임스탬프
	block.RoleInfo = roleInfo
	block.Hash = utils.Hash(b)
	PersistBlock(block)
	return block
}

// 최초 상태의 블록체인에 제네시스 블록 추가 (위 AddBlock과 구분한 이유는 비트코인의 타임스탬프 등 여러가지 조건을 넣고 싶어서)
func (b *blockchain) AddGenesisBlock() *Block {
	block := createGenesisBlock()
	b.UpdateBlockchain(block)
	return block
}

// 제네시스 트랜잭션 생성
func makeGenesisTx() *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // 소유주는 채굴자
	}
	txOuts := []*TxOut{
		{"", proposalReward},
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

// 제네시스 블록 확인
func (m *mempool) GenesisTxToConfirm() []*Tx {
	coinbase := makeGenesisTx()
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}
