package pos

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

var (
	once sync.Once // sync 패키지
	v    *validatorInfo
	m    *minerInfo
)

const (
	stakingAddress = "e9d11618e700bad8d5aa12d44531036d5995b21cec01443f8dc27a92f3b22ab3b5879eb20980be45ec6fcec5e153842f42d5cf7b5632d12ccdc4160bf19bd270"
	epoch          = 3
	genesisHeight  = 1
)

type minerInfo struct {
	Address        string
	Port           string
	SelectedHeight int
}

type validatorInfo struct {
	Address        string
	Port           string
	SelectedHeight int
}

func (v *validatorInfo) selectValidator(stakingList []*blockchain.StakingInfo) *blockchain.StakingInfo {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	v.Address = selected.Address
	v.Port = selected.Port
	v.SelectedHeight = blockchain.Blockchain().Height
	fmt.Println("validatorInfo: ", randNum, utils.ToString(v))
	return nil
}

func (m *minerInfo) selectMiner(stakingList []*blockchain.StakingInfo) *blockchain.StakingInfo {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	m.Address = selected.Address
	m.Port = selected.Port
	m.SelectedHeight = blockchain.Blockchain().Height
	fmt.Println("minerInfo: ", randNum, utils.ToString(m))
	return nil
}

// func initRoles() (*minerInfo, *validatorInfo) {

// 	return nil, nil
// }

func Selector() {
	for {
		once.Do(func() {
			// if blockchain.Blockchain().Height == genesisHeight {
			fmt.Println("here?")
			m = &minerInfo{
				Address:        "",
				Port:           "",
				SelectedHeight: blockchain.Blockchain().Height,
			}
			v = &validatorInfo{
				Address:        "",
				Port:           "",
				SelectedHeight: blockchain.Blockchain().Height,
			}
			// }
		})

		if m.SelectedHeight != blockchain.Blockchain().Height {
			fmt.Println(blockchain.Blockchain().Height)
			fmt.Println("validator height: ", v.SelectedHeight)

			_, stakingWalletTx, _ := blockchain.UTxOutsByStakingAddress(stakingAddress, blockchain.Blockchain())
			stakingInfoList := blockchain.GetStakingList(stakingWalletTx, blockchain.Blockchain())

			fmt.Println("gap: ", blockchain.Blockchain().Height-v.SelectedHeight)
			m.selectMiner(stakingInfoList)
			if blockchain.Blockchain().Height-v.SelectedHeight == epoch {
				v.selectValidator(stakingInfoList)
			}
		}
	}
}

// func validator() {

// }

// func miner() {

// }
