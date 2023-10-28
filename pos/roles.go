package pos

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/utils"
)

const stakingAddress = "e9d11618e700bad8d5aa12d44531036d5995b21cec01443f8dc27a92f3b22ab3b5879eb20980be45ec6fcec5e153842f42d5cf7b5632d12ccdc4160bf19bd270"

func selectValidator(stakingList []*blockchain.StakingInfo) *blockchain.StakingInfo {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	fmt.Println("validator:", randNum, utils.ToString(selected))
	return nil
}

func selectMiner(stakingList []*blockchain.StakingInfo) *blockchain.StakingInfo {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	fmt.Println("miner: ", randNum, utils.ToString(selected))
	return nil
}

func Selector() {
	tickerA := time.NewTicker(3 * time.Second)
	tickerB := time.NewTicker(15 * time.Second)

	_, stakingWalletTx, _ := blockchain.UTxOutsByStakingAddress(stakingAddress, blockchain.Blockchain())
	stakingInfoList := blockchain.GetStakingList(stakingWalletTx, blockchain.Blockchain())

	for {
		select {
		case <-tickerA.C:
			selectValidator(stakingInfoList)
		case <-tickerB.C:
			selectMiner(stakingInfoList)
		}
	}
}
