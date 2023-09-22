package wallet

import (
	"crypto/ecdsa"
	"os"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet // 이걸 소문자로 써서 자유롭게 공유하는게 아니라, 아래의 Wallet 함수로 드러나게 할 예정

func hasWalletFile() bool {
	_, err := os.Stat("nomadcoin.wallet") // 파일이 존재하는지
	return !os.IsNotExist(err)            // os.Stat에서 받아온 err를 확인하고 지갑 파일이 없다면 true
}

func Wallet() *wallet {
	if w == nil {
		// 지갑의 유뮤 확인
		if hasWalletFile() {
			// 없다면 비공개키를 생성해서 파일에 저장
		} else {
			// 이미 있다면 파일로부터 지갑을 복구
		}
	}
	return w
}
