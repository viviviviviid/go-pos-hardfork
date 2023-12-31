package blockchain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/viviviviviid/go-coin/utils"
)

var (
	ResLeastStaker = "PoS requires at least 4 Stakers to run" // 스테이킹 노드가 최소 4개여야 한다는 응답
)

const (
	slotTime      = 12
	Epoch         = 3 // 이더리움은 32개의 슬롯
	genesisHeight = 1
)

// 검증자 선출
func (r *RoleInfo) selectValidator(b *blockchain, stakingList []*StakingInfo) {
	selectedNumbers := make(map[int]bool)
	var result []int

	for len(result) < 3 {
		randNum := rand.Intn(len(stakingList))
		if !selectedNumbers[randNum] { // 선택한 숫자가 아직 선택되지 않았다면
			selectedNumbers[randNum] = true  // 선택한 숫자를 맵에 추가
			result = append(result, randNum) // 결과 슬라이스에 추가
		}
	}

	for _, num := range result {
		selected := stakingList[num]
		r.ValidatorAddress = append(r.ValidatorAddress, selected.Address)
		r.ValidatorPort = append(r.ValidatorPort, selected.Port)
	}

	r.ValidatorSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

// 제안자 선출
func (r *RoleInfo) selectProposer(b *blockchain, stakingList []*StakingInfo) {
	var selected *StakingInfo
	for {
		check := 0
		randNum := rand.Intn(len(stakingList))
		selected = stakingList[randNum]
		for _, validatorAddress := range r.ValidatorAddress {
			if selected.Address != validatorAddress {
				check++
			} else {
				break
			}
		}
		if check == 3 {
			break
		}
	}
	r.ProposerAddress = selected.Address
	r.ProposerPort = selected.Port
	r.ProposerSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

// 지정된 슬롯과 에포크에 따른 검증자와 제안자 선출
func (b *blockchain) Selector() (*RoleInfo, string) {
	r := &RoleInfo{}

	_, stakingWalletTx, _ := UTxOutsByStakingAddress(utils.StakingAddress, b)
	stakingInfoList := GetStakingList(stakingWalletTx, b)

	if len(stakingInfoList) <= 3 {
		fmt.Println(ResLeastStaker)
		time.Sleep(slotTime * time.Second)
		return nil, ResLeastStaker
	}

	if b.Height%Epoch == 0 {
		r.selectValidator(b, stakingInfoList)
		r.selectProposer(b, stakingInfoList)
	} else {
		block, _ := FindBlock(b.NewestHash)
		r.ValidatorSelectedHeight = block.RoleInfo.ValidatorSelectedHeight
		r.ValidatorAddress = block.RoleInfo.ValidatorAddress
		r.ValidatorPort = block.RoleInfo.ValidatorPort
		r.selectProposer(b, stakingInfoList)
	}

	fmt.Printf("Seleted Roles for the next block:\n%s\n", utils.ToString(r))

	return r, ""
}

// 검증자들의 검증 결과 과반수 계산
func CalculateMajority(v []*ValidatedInfo) bool {
	pass := 0
	fail := 0
	for _, r := range v {
		if r.Result {
			pass++
		} else {
			fail++
		}
	}
	fmt.Printf("PASS: %d \nFAIL: %d\n", pass, fail)
	return pass > fail
}

// 검증 중 RoleInfo 내용 비교
func compareRoleInfo(r1, r2 *RoleInfo) bool {
	return r1.ProposerAddress == r2.ProposerAddress &&
		r1.ProposerPort == r2.ProposerPort &&
		r1.ProposerSelectedHeight == r2.ProposerSelectedHeight &&
		utils.CompareStringSlices(r1.ValidatorAddress, r2.ValidatorAddress) &&
		utils.CompareStringSlices(r1.ValidatorPort, r2.ValidatorPort) &&
		r1.ValidatorSelectedHeight == r2.ValidatorSelectedHeight
}

// 블록 제안 성공 유무
func (b *blockchain) CheckProposalSuccess(lastHeight int) {
	if b.Height == lastHeight {
		fmt.Println("Proposal Rejected.")
	} else if b.Height-lastHeight == 1 {
		fmt.Println("Added and broadcasted the block done.")
	} else {
		fmt.Println("Warning: Block Height was twisted.")
	}
}
