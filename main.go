package main

import (
	"fmt"
	"nomadcoin/blockchain"
)

func main() {
	chain := blockchain.GenerateBlockchain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")

	for _, block := range chain.ListOfBlocks() {
		fmt.Println("Data: " + block.Data)
		fmt.Println("Hash: " + block.Hash)
		fmt.Println("PrevHash: " + block.PrevHash + "\n")
	}
}
