package p2p

type MessageKind int

const (
	MessageNewestBlock       MessageKind = iota // StatusOK = 200 과 같은 스테이터스 변수와 같은 형식스로 진행
	MessageAllBlocksRequest                     // iota 밑에 있어서, 변수들의 숫자가 0부터 1씩 증가하는 형태로 가지게 될것이고
	MessageAllBlocksResponse                    // iota 밑에 있어서, 변수들의 종류도 MessageKind가 될것
)

type Message struct { // 다른 언어와 소통하기에도 적합한 메세지 형식 정의
	Kind    MessageKind
	Payload []byte
}

// func sendNewestBlock(p *peer) {
// 	m := makeMessage()
// }
