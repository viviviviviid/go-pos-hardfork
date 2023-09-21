package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/viviviviviid/go-coin/utils"
)

func Start() {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // 타원곡선 디지털서명 알고리즘 라이브러리
	utils.HandleErr(err)

	message := "I love you"

	hashedMessage := utils.Hash(message)

	hashAsByte, err := hex.DecodeString(hashedMessage) // 16진수 string -> byte
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsByte) // 현재는 서명이 r과 s로 나뉘어 있음
	utils.HandleErr(err)

	fmt.Printf("R: %d\nS: %d", r, s)
}

// fmt.Println("Private Key: ", privateKey.D)
// // privateKey struct안에 public키가 들어있음 // 정의 보면됨
// fmt.Println("Public Key, x, y: ", privateKey.X, privateKey.Y)
