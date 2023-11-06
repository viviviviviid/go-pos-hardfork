package pos

import (
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/rest"
	"github.com/viviviviviid/go-coin/utils"
)

const (
	stakingPort     = "3000"
	nodeSettingTime = 20
)

// 3000 port has a staking address of this PoS chain. So this port is selecting Proposal and Validator
func PoS(aPort int) {
	go rest.Start(aPort)
	time.Sleep(nodeSettingTime * time.Second)

	for {
		height := blockchain.Blockchain().Height
		roleInfo, err := blockchain.Blockchain().Selector()
		if err != nil {
			utils.HandleErr(err)
			return
		}
		p2p.PointingProposal(roleInfo)
		p2p.PointingValidator(roleInfo)

		for blockchain.Blockchain().Height == height {
			// 블록추가될때까지 대기
		}
	}
}
