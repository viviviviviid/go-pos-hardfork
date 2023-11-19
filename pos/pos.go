// pos 패키지는 Proof of Stake (PoS) 알고리즘과 관련된 함수를 제공합니다.
package pos

import (
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/rest"
)

const (
	nodeSettingTime = 20 // 10개 노드 동시 실행시 세팅에 걸리는 시간으로, 간섭이 없도록 Lock
	slotTime        = 12 // 블록 하나가 추가되는 이상적인 시간 (이더리움 기준)
)

// PoS의 기둥이 되는 함수. 슬롯과 에포크마다 스테이킹 리스트를 토대로 검증자와 제안자를 선정 후, 제안 성공 여부를 따지는 로직을 반복한다.
func PoS(aPort int) {
	go rest.Start(aPort)
	time.Sleep(nodeSettingTime * time.Second)
	for {
		lastHeight := blockchain.Blockchain().Height
		roleInfo, err := blockchain.Blockchain().Selector()
		if err != "" {
			continue
		}
		p2p.PointingProposal(roleInfo)
		p2p.PointingValidator(roleInfo)
		time.Sleep(slotTime * time.Second)
		blockchain.Blockchain().CheckProposalSuccess(lastHeight)
	}
}
