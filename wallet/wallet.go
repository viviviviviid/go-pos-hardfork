// wallet 패키지는 지갑과 관련된 함수를 제공합니다.
package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/viviviviviid/go-coin/utils"
)

// 지갑 파일의 기본 이름
const (
	fileName string = ".wallet"
)

// 파일 입출력 관련 메서드를 정의하는 인터페이스
type fileLayer interface {
	hasWalletFile(fileNamebyPort string) bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

// 파일이 존재하는지 확인하는 메서드
func (layer) hasWalletFile(fileNamebyPort string) bool {
	_, err := os.Stat(fileNamebyPort)
	return !os.IsNotExist(err)
}

// 파일에 데이터를 쓰는 메서드
func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// 파일에서 데이터를 읽어오는 메서드
func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}

// 개인 키와 주소 정보를 저장하는 구조체
type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

// 타원곡선 디지털 서명 알고리즘(ECDSA)을 사용하여 개인 키를 생성
func createPrivateKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

// 주어진 파일에 개인 키를 저장
func persistKey(fileNamebyPort string, key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileNamebyPort, bytes, 0644)
	utils.HandleErr(err)
}

// 주어진 파일에서 개인 키를 복원
func restoreKey(fileNamebyPort string) (key *ecdsa.PrivateKey) {
	keyAsBytes, err := files.readFile(fileNamebyPort)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

// 두 개의 big.Int를 결합하고 16진수 문자열로 인코딩
func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

// ECDSA 퍼블릭 키를 생성하는 메서드
func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

// 주어진 지갑 정보로 페이로드를 서명
func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// 대리 서명을 생성하는 메서드로, 타인의 지갑 정보로 서명
func DelegateSign(payload string) string {
	w := DelegateWallet()
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// 16진수 문자열을 big.Int 형태로 복원
func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

// 서명된 데이터의 유효성을 검증
func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

// 대리 서명을 위한 지갑 정보를 생성
func DelegateWallet() *wallet {
	wallet := &wallet{}
	path := "./wallets/" + "3000" + fileName
	wallet.privateKey = restoreKey(path)
	wallet.Address = aFromK(wallet.privateKey)
	return wallet
}

// 주어진 포트에 해당하는 지갑 정보를 반환
func Wallet(port string) *wallet {
	if w == nil {
		w = &wallet{}
		path := "./wallets/" + port + fileName
		if files.hasWalletFile(path) {
			w.privateKey = restoreKey(path)
		} else {
			key := createPrivateKey()
			persistKey(path, key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
