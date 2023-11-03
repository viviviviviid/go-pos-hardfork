package blockchain

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/viviviviviid/go-coin/utils"
)

var (
	r              = &RoleInfo{}
	ErrLeastStaker = errors.New("PoS requires at least 4 Stakers to run")
	stakingAddress = "0ed84571488f4474f83291fcb29f73348983df8ac535873d44acb7cdb38035a547720ab7f64d2fce2811fd7b3b8db7b9100e8c054f88970aa415ddced6a12beb"
)

const (
	epoch         = 3 // 실제는 32
	genesisHeight = 1
)

func (r *RoleInfo) selectValidator(b *blockchain, stakingList []*StakingInfo) {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	r.ValidatorAddress = selected.Address
	r.ValidatorPort = selected.Port
	r.ValidatorSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

func (r *RoleInfo) selectProposal(b *blockchain, stakingList []*StakingInfo) {
	randNum := rand.Intn(len(stakingList))
	selected := stakingList[randNum]
	r.ProposalAddress = selected.Address
	r.ProposalPort = selected.Port
	r.ProposalSelectedHeight = b.Height + 1 // b.Height는 현재 높이이고, 이제 추가할 블록의 높이는 +1로 해야함
}

func (b *blockchain) Selector() (*RoleInfo, error) {
	_, stakingWalletTx, _ := UTxOutsByStakingAddress(stakingAddress, b)
	stakingInfoList := GetStakingList(stakingWalletTx, b)

	if len(stakingInfoList) <= 3 {
		return nil, ErrLeastStaker
	}

	r.selectProposal(b, stakingInfoList)

	if b.Height%epoch == 0 {
		r.selectValidator(b, stakingInfoList)
	} else {
		block, _ := FindBlock(b.NewestHash)
		r.ValidatorSelectedHeight = block.RoleInfo.ValidatorSelectedHeight
		r.ValidatorAddress = block.RoleInfo.ValidatorAddress
		r.ValidatorPort = block.RoleInfo.ValidatorPort
	}

	fmt.Println(utils.ToString(r))

	return r, nil
}
