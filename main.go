package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash // 블록배열의 마지막 item의 hash
	}
	return ""
}

func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash)) // Sum256은 byte형태의 slice를 인자로 받음
	newBlock.hash = fmt.Sprintf("%x", hash)                          // Sprint로 return // %x 로 hash 값을 포맷해야 우리가 흔히보는 16진수 해시값이 나옴.
	b.blocks = append(b.blocks, newBlock)                            // 블록 추가
}

func (b *blockchain) listBlock() {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("Prev Hash: %s\n", block.prevHash)
	}
}

func main() {
	chain := blockchain{}
	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listBlock()
}
