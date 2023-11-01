package controller

import (
	"strconv"
	"time"

	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/p2p"
	"github.com/viviviviviid/go-coin/rest"
)

var roleInfo *blockchain.RoleInfo

const (
	stakingPort     = "3000"
	nodeSettingTime = 20
)

func Auto(aPort int) {
	go rest.Start(aPort)

	time.Sleep(nodeSettingTime * time.Second)
	port := strconv.Itoa(aPort)

	if port == stakingPort {
		roleInfo = blockchain.Blockchain().Selector()
		p2p.PointingMiner(roleInfo)
	}
	time.Sleep(1000 * time.Second)

}

// func postRequest(url string, payload []byte) {
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("Response status:", resp.Status)
// }
