package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"

	"github.com/viviviviviid/go-coin/utils"
)

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet // 이걸 소문자로 써서 자유롭게 공유하는게 아니라, 아래의 Wallet 함수로 드러나게 할 예정

func hasWalletFile() bool {
	_, err := os.Stat(fileName) // 파일이 존재하는지
	return !os.IsNotExist(err)  // os.Stat에서 받아온 err를 확인하고 지갑 파일이 없다면 true
}

func createPriveKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) { // key 저장
	bytes, err := x509.MarshalECPrivateKey(key) // bytes는 복붙가능하기때문에 변환할 필요없이 파일에 박으면 됨
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644) // 0644 : 읽기와 쓰기 허용
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}

		// 지갑의 유뮤 확인
		if hasWalletFile() {
			// 이미 있다면 파일로부터 지갑을 복구

		} else {
			// 없다면 비공개키를 생성해서 파일에 저장
			key := createPriveKey()
			persistKey(key)
			w.privateKey = key
		}

	}
	return w
}
