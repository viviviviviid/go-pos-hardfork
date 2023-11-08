package pos

import (
	"fmt"
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/rest"
	"github.com/viviviviviid/go-coin/utils"
)

const (
	stakingPort     = "3000"
	nodeSettingTime = 20
	slotTime        = 10
)

// 3000 port has a staking address of this PoS chain. So this port is selecting Proposal and Validator
func PoS(aPort int) {
	go rest.Start(aPort)
	time.Sleep(nodeSettingTime * time.Second)
	for {
		lastHeight := blockchain.Blockchain().Height
		roleInfo, err := blockchain.Blockchain().Selector()
		if err != nil {
			utils.HandleErr(err)
			return
		}
		p2p.PointingProposal(roleInfo)
		p2p.PointingValidator(roleInfo)
		time.Sleep(slotTime * time.Second)
		if blockchain.Blockchain().Height == lastHeight {
			fmt.Println("Proposal Rejected.")
		} else if blockchain.Blockchain().Height-lastHeight == 1 {
			fmt.Println("Added and broadcasted the block done.")
		} else {
			fmt.Println("Warning: Block Height was twisted.")
		}
	}
}
