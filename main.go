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
		return b.blocks[len(b.blocks)-1].hash
	}
	return ""
}
func (b *blockchain) getHash(hash string) string {
	hashed := sha256.Sum256([]byte(hash))
	return fmt.Sprintf("%x", hashed)
}
func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	newBlock.hash = b.getHash(data + newBlock.prevHash)

	b.blocks = append(b.blocks, newBlock)
}
func (b *blockchain) listofBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("PrevHash: %s\n", block.prevHash)
	}
}

/*
b1hash=data+""

b2hash=data+b1hash
*/
func main() {
	chain := blockchain{}
	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listofBlocks()
}
