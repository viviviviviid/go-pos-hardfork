package main

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func main() {
	// genesisBlock := block{"Genesis Block", "", ""}
	// hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
	// // Sum256은 byte형태의 slice를 인자로 받음

	// hexHash := fmt.Sprintf("%x", hash) // Sprint로 return
	// // %x 로 hash 값을 포맷해야 우리가 흔히보는 16진수 해시값이 나옴.

	// genesisBlock.hash = hexHash

	// fmt.Println(genesisBlock)
}
