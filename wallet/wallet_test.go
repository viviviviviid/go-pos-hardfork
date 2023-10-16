package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testPrivKey string = "30770201010420035da6e40df9b8f07c39e7c26b18739accaa7d62e874fb748971e0669904b338a00a06082a8648ce3d030107a144034200044ca6415b21fff2ad7641fe16df9a41808c96e2c6a09912da19620c7d5eee8ded5b54c71bc15f26a99e23f9db01d7af29c82adc7abff32f1767ac8f505109c640"
	testPayload string = "00bc9a6fdcf6d884312e8422b986371972398267e069b39185c40a361ed659d7"
	testSig     string = "9d55d4923bedff540ff8adc725a0c92ce8896a112ef211396a01dc198225e4c5dbe4bef6526c25ff012a377cc78a3c920287b72601f24df4f138163c2f9eafb1"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool { // 실제 Wallet()에서 hasWalletFile을 호춣하고있지만 테스트에서는 이 내용으로 바뀐 hasWalletFile을 호출함
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey) // 원래는 지갑 파일의 bytes를 가져오지만, 테스트이므로, 테스트 지갑의 프라이빗키를 가져와서 마샬로 bytes화 해서 넘겨줌
}

func TestWallet(t *testing.T) {
	t.Run("New Wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return false
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return true
			},
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testPrivKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should retrun a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"10bc9a6fdcf6d884312e8422b986371972398267e069b39185c40a361ed659d7", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify test Signature and test payload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex")
	}
}
