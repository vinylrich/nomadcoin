package main

import (
	"fmt"
	"nomadcoin/blockchain"
)

/*
B1
	b1Hash = {data + "x"}
	//data+"x" 가아닌 data+"x!"가 되면
	b2,b3도 모두 바뀜
B2
	b2Hash = {data + b1Hash}
B3
	b3Hash = {data + b2Hash}

1번째 block
data를 그대로 hash하고 prevHash는 비워져있음

2번째 block
prevHash, 즉, 전에 있던 hash를 hashing하고 1번째

*/

func main() {
	chain := blockchain.GetBlockchain()

	chain.AddBlock("Second Block")
	chain.AddBlock("Third Blok")

	for _, block := range chain.ListOfBlock() {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)

	}
}
