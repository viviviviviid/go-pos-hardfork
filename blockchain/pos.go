package blockchain

import (
	"math/rand"
)

var (
	r              = &RoleInfo{}
	stakingAddress = "0ba0b66c37ffe7037b114ca5142bb0c6796ad910ead1022d565bee5f86dcc9cc6bc8209cd879cc855ccfbd7ed6113b29ac0ca9ecb4c1a76dafe6a39cbf246dbe"
)

// const (
// 	epoch         = 3
// 	genesisHeight = 1
// )

// func (r *RoleInfo) selectValidator(b *blockchain, stakingList []*StakingInfo) {
// 	randNum := rand.Intn(len(stakingList))
// 	selected := stakingList[randNum]
// 	r.ValidatorAddress = selected.Address
// 	r.ValidatorPort = selected.Port
// 	r.ValidatorSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
// }

func (r *RoleInfo) selectMiner(b *blockchain, stakingList []*StakingInfo) {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	r.MinerAddress = selected.Address
	r.MinerPort = selected.Port
	r.MinerSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

func (b *blockchain) Selector() *RoleInfo {
	_, stakingWalletTx, _ := UTxOutsByStakingAddress(stakingAddress, b)
	stakingInfoList := GetStakingList(stakingWalletTx, b)

	if len(stakingInfoList) == 0 {
		return nil
	}

	r.selectMiner(b, stakingInfoList)
	// if b.Height-r.ValidatorSelectedHeight == epoch {
	// 	r.selectValidator(b, stakingInfoList)
	// }

	return r
}

func Miner() {

}

func Validator() {

}
