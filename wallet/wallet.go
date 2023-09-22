package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/viviviviviid/go-coin/utils"
)

const (
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	privateKey    string = "3077020101042028f3ae83f5b6726ae5587ddc2238665b31db5c4ab850ddf08028b94df523a66ca00a06082a8648ce3d030107a144034200044f5f301c270532c412f62e1333927337bd0c60e071902e0623e0532725bca5d3a281a4c80c019144710bea802212627120b3e45d128925c68004f68a983e42b4"
	signature     string = "2babd4181f18fa45d3eaa08d345a700ff29201c40b1544130e48a66011b53b736c9acd2a526e7db31232877fcc01ab1d85522013e56fe1dda8d32fa7db0f3474"
)

func Start() {

	// 비공개키를 byte로 복구
	privByte, err := hex.DecodeString(privateKey)
	// []byte()로 안하는 이유 -> 포맷이 hexadecimal bytes인지 확인하고 진행해야함
	// 파일로 저장되기 떄문에 수정될 가능성이 있기 때문
	utils.HandleErr(err)

	private, err := x509.ParseECPrivateKey(privByte) // bytes를 받아서 비공개키를 반환
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	rBytes := sigBytes[:len(sigBytes)/2] // r과 s로 반갈하기
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	hashBytes, err := hex.DecodeString(hashedMessage)

	utils.HandleErr(err)

	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)

	fmt.Println(ok)

}
