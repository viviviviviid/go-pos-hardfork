package blockchain

import (
	"fmt"
	"math/rand"

	"github.com/viviviviviid/go-coin/wallet"
)

var (
	r              *RoleInfo
	stakingAddress = wallet.Wallet("3000")
)

const (
	epoch         = 3
	genesisHeight = 1
)

func (r *RoleInfo) selectValidator(b *blockchain, stakingList []*StakingInfo) {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	r.ValidatorAddress = selected.Address
	r.ValidatorPort = selected.Port
	r.ValidatorSelectedHeight = b.Height
}

func (r *RoleInfo) selectMiner(b *blockchain, stakingList []*StakingInfo) {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	r.MinerAddress = selected.Address
	r.MinerPort = selected.Port
	r.MinerSelectedHeight = b.Height
}

func (b *blockchain) Selector() *RoleInfo {
	_, stakingWalletTx, _ := UTxOutsByStakingAddress(stakingAddress.Address, b)
	stakingInfoList := GetStakingList(stakingWalletTx, b)

	r.selectMiner(b, stakingInfoList)
	fmt.Println("gap: ", b.Height-r.ValidatorSelectedHeight)
	if b.Height-r.ValidatorSelectedHeight == epoch {
		r.selectValidator(b, stakingInfoList)
	}

	return r
}

// newBlock := blockchain.Blockchain().AddBlock(port[1:])
// 		p2p.BroadcastNewBlock(newBlock)

// func Broadcast -> p2p에서 작업해야함
