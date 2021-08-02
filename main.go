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

/*
B1
	b1Hash = {data + "x"}
	//data+"x" 가아닌 data+"x!"가 되면
	b2,b3도 모두 바뀜
B2
	b2Hash = {data + b1Hash}
B3
	b3Hash = {data + b2Hash}
*/
//one-way function

func main() {
	genesisBlock := block{"Genesis Block", "", ""}
	//genesisBlock is First Block

	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))

	hexHash := fmt.Sprintf("%x", hash)

	genesisBlock.hash = hexHash

	fmt.Printf("%v", genesisBlock.hash)
}
