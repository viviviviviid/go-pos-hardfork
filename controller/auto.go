package controller

import (
	"strconv"
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

// 추후 3000번은 pos 용도로 mode=auto로 실행, 나머지는 rest로만 실행
func Auto(aPort int) {
	go rest.Start(aPort)

	time.Sleep(nodeSettingTime * time.Second)
	port := strconv.Itoa(aPort)

	for {
		if port == stakingPort {
			roleInfo, err := blockchain.Blockchain().Selector()
			if err != nil {
				utils.HandleErr(err)
				return
			}
			p2p.PointingMiner(roleInfo)
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(1000 * time.Second)
		}
	}
}
