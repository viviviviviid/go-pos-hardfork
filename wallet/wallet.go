package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/viviviviviid/go-coin/utils"
)

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
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

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload) // []bytes()를 안쓰는 이유는 길이 관련으로 오류가 생기는걸 확인하기위해
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]  // 중간까지
	secondHalfBytes := bytes[len(bytes)/2:] // 중간부터 끝까지
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{ // 퍼블릭키 만들기
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
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
		w.Address = aFromK(w.privateKey)

	}
	return w
}
