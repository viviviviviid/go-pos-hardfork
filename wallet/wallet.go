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
	address    string
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

func restoreKey() (key *ecdsa.PrivateKey) { // *ecdsa.PrivateKey 형식의 key를 선언 및 초기화
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes) // 이미 함수의 반환 구조에서 초기화되었으므로 key를 갱신만 해줘도 됨.
	utils.HandleErr(err)
	return // 함수의 반환 구조에서 뭘 반환할지 알려줬으므로, return 다음에 뭔가를 안써줘도 됨
} // return에 비어있는지 아닌지 확인해야하므로 긴 함수에서는 귀찮음이 가중될 수 있음 -> 알고만 있기

func aFromK(key *ecdsa.PrivateKey) string {

}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}

		// 지갑의 유뮤 확인
		if hasWalletFile() {
			// 이미 있다면 파일로부터 지갑을 복구
			w.privateKey = restoreKey()

		} else {
			// 없다면 비공개키를 생성해서 파일에 저장
			key := createPriveKey()
			persistKey(key)
			w.privateKey = key
		}
		w.address = aFromK(w.privateKey)

	}
	return w
}
