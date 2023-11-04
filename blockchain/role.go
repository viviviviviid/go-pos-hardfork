package blockchain

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/viviviviviid/go-coin/utils"
)

var (
	ErrLeastStaker = errors.New("PoS requires at least 4 Stakers to run")
	stakingAddress = "c8546a75af42fd63669afa3d2e72b3567790aa8f2a54da1abb94ec03239c76638f45ada90e6e2a5af42efff001a66d90106fa898ae55d3168b11d9e120a0763d"
)

const (
	epoch         = 3 // 실제는 32
	genesisHeight = 1
)

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

func (r *RoleInfo) selectProposal(b *blockchain, stakingList []*StakingInfo) {
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
	r.ProposalAddress = selected.Address
	r.ProposalPort = selected.Port
	r.ProposalSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

func (b *blockchain) Selector() (*RoleInfo, error) {
	r := &RoleInfo{}

	_, stakingWalletTx, _ := UTxOutsByStakingAddress(stakingAddress, b)
	stakingInfoList := GetStakingList(stakingWalletTx, b)

	if len(stakingInfoList) <= 3 {
		return nil, ErrLeastStaker
	}

	if b.Height%epoch == 0 {
		r.selectValidator(b, stakingInfoList)
		r.selectProposal(b, stakingInfoList)
	} else {
		block, _ := FindBlock(b.NewestHash)
		r.ValidatorSelectedHeight = block.RoleInfo.ValidatorSelectedHeight
		r.ValidatorAddress = block.RoleInfo.ValidatorAddress
		r.ValidatorPort = block.RoleInfo.ValidatorPort
		r.selectProposal(b, stakingInfoList)
	}

	fmt.Println(utils.ToString(r))

	return r, nil
}
