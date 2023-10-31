package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/utils"
)

func Auto(aPort int) {
	port := strconv.Itoa(aPort)
	if port == "3000" {
		time.Sleep(3 * time.Second)
		for len(p2p.AllPeers(&p2p.Peers)) != 10 {
			for peerPort := 4009; peerPort >= 4000; peerPort-- {
				p2p.AddPeer("127.0.0.1", strconv.Itoa(peerPort), port, true)
			}
		}
	}

	for {
		fmt.Println("Start!")
		roleInfo := blockchain.Blockchain().Selector()
		if port != roleInfo.MinerPort { // 현재 노드의 포트가 r의 miner 포트와 동일할떄 정상적으로 진행.
			fmt.Println("Re Select")
			time.Sleep(3 * time.Second)
			continue
		}
		newBlock := blockchain.Blockchain().AddBlock(port, roleInfo)
		p2p.BroadcastNewBlock(newBlock)
		hash := blockchain.Blockchain().NewestHash
		block, _ := blockchain.FindBlock(hash)
		fmt.Println("Newest Block: ", utils.ToString(block))
		time.Sleep(1 * time.Second)
		fmt.Println("Next Selector!")
		for i := 5; i > 0; i-- {
			fmt.Printf("%d seconds remaining...\n", i)
			time.Sleep(1 * time.Second)
		}
	}

	// 	최초
	// 	체인: 역할 선택

	// r := blockchain.Blockchain().Selector()

	// blockchain.Miner()

	// 	체인: 노드에게 정보전달
	// 	func p2p.rolePointing(r) {
	// 		// 노드에게 r 정보 전달 내용
	// 		miner() // 아래에 명시해둔 내용
	// 	}

	// 	func p2p.miner(r) {

	// 		newBlock := blockchain.Blockchain().AddBlock(port[1:], r) // 여기서 AddBlock이 r의 내용을 집어넣도록 변경해줘야함
	// 	}

	// 	노드: 생성자가 블록을 생성한다. 현재 내용이 들어가야함 (10초 딜레이 주기)
	// 	p2p.@@@@@@(r)내부에서 rest의 blocks에서 진행할 내용을 실행
	// 	race 안일어나도록 조심

	// 	그 이후

	// 	체인: 블록이 추가가 되었으면 곧바로 역할 선택 및 노드에게 정보를 전달해 준다.

	// 	노드: 노드는 그 정보를 가지고 블록을 생성한다.

}
