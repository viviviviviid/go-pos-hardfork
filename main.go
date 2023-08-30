package main

import (
	"github.com/viviviviviid/go-coin/blockchain"
	"github.com/viviviviviid/go-coin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}

// Mux : Multiplexer
// 하나의 Mux가 3000번과 4000번의 "/" 을 동시에 보고있기때문에 오류가 발생함
// port가 서로 달라서 분리된것 같지만, http.ListenAndServe 함수에서 두번째 인자에 nil을 넣으면,
// DefaultServeMux 가 default로 사용되기 때문
// 즉 우리만의 Mux를 만들어서 사용해야함
