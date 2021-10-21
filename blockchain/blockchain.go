package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) getHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getPrevHash() string {
	totalBlocks := len(GenerateBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GenerateBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getPrevHash()}
	newBlock.getHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockchain) ListOfBlocks() []*block {
	return b.blocks
}
func GenerateBlockchain() *blockchain {
	if b == nil {
		//데이터베이스를 사용하면
		//초기화단계에서 blockchain을
		//database에서 가져올 수 있음
		//어떻게 생성되고, 관리되는지
		//b는 한번만 초기화 됨(nil일때)
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
