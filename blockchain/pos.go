package blockchain

import (
	"fmt"
	"math/rand"

	"github.com/viviviviviid/go-coin/utils"
)

var (
	r              = &RoleInfo{}
	stakingAddress = "0ed84571488f4474f83291fcb29f73348983df8ac535873d44acb7cdb38035a547720ab7f64d2fce2811fd7b3b8db7b9100e8c054f88970aa415ddced6a12beb"
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

	fmt.Println(utils.ToString(stakingInfoList))

	if len(stakingInfoList) == 0 {
		fmt.Println("Anyone not staked, now")
		return nil
	}

	r.selectMiner(b, stakingInfoList)
	// if b.Height-r.ValidatorSelectedHeight == epoch {
	// 	r.selectValidator(b, stakingInfoList)
	// }

	return r
}
