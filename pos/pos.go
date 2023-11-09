package pos

import (
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/rest"
)

const (
	stakingPort     = "3000"
	nodeSettingTime = 20
	slotTime        = 10
)

func PoS(aPort int) {
	go rest.Start(aPort)
	time.Sleep(nodeSettingTime * time.Second)
	for {
		lastHeight := blockchain.Blockchain().Height
		roleInfo, err := blockchain.Blockchain().Selector() // 블록 제안자와 검증자 선정 (제안자: 1블록마다, 검증자 3블록마다 재선정)
		if err != "" {
			continue
		}
		p2p.PointingProposal(roleInfo)
		p2p.PointingValidator(roleInfo)
		time.Sleep(slotTime * time.Second)
		blockchain.Blockchain().CheckProposalSuccess(lastHeight)
	}
}
