package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/viviviviviid/go-coin/utils"
)

const (
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	// privateKey    string = "3077020101042028f3ae83f5b6726ae5587ddc2238665b31db5c4ab850ddf08028b94df523a66ca00a06082a8648ce3d030107a144034200044f5f301c270532c412f62e1333927337bd0c60e071902e0623e0532725bca5d3a281a4c80c019144710bea802212627120b3e45d128925c68004f68a983e42b4"
	// signature     string = "2babd4181f18fa45d3eaa08d345a700ff29201c40b1544130e48a66011b53b736c9acd2a526e7db31232877fcc01ab1d85522013e56fe1dda8d32fa7db0f3474"
)

func Start() {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // 타원곡선 디지털서명 알고리즘 라이브러리
	utils.HandleErr(err)

	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey) // x509는 key를 parsing하는 패키지

	fmt.Printf("%x\n", keyAsBytes)

	hashAsBytes, err := hex.DecodeString(hashedMessage) // 16진수 string -> byte
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes) // 현재는 서명이 r과 s로 나뉘어 있음
	utils.HandleErr(err)

	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("%x\n", signature)
}

// fmt.Println("Private Key: ", privateKey.D)
// // privateKey struct안에 public키가 들어있음 // 정의 보면됨
// fmt.Println("Public Key, x, y: ", privateKey.X, privateKey.Y)
