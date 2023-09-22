package wallet

import (
	"crypto/x509"
	"encoding/hex"
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

	restoredKey, err := x509.ParseECPrivateKey(privByte) // bytes를 받아서 비공개키를 반환
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	rBytes := sigBytes[:len(sigBytes)/2] // r과 s로 반갈하기
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

}

// 한 일 : 비공개키의 문자열을 가져와서, hex 패키지의 DecodeString를 이용하여 변환하였고
// 그 뒤, 저번에 비공개키를 byte로 변환했던 그 패키지를 다시 사용해서, byte를 다시 비공개키로 변경
// 비공개키를 복구한 다음에, 서명을 복구, 서명은 두 slice의 조합
// 서명의 byte를 받아왔고, r과 s로 분리
// big.Int 로 생성해서 byte값으로 set했음
