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

func PoS(aPort int) {
	go rest.Start(aPort)
	time.Sleep(nodeSettingTime * time.Second)

	for {
		roleInfo, err := blockchain.Blockchain().Selector()
		if err != nil {
			utils.HandleErr(err)
			return
		}
		p2p.PointingMiner(roleInfo)
		time.Sleep(5 * time.Second)
	}
}
