package blockchain

import (
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) fromBytes(data []byte) {
	utils.FromBytes(b, data)
}
func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}
func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
func Blockchain() *blockchain {
	if b == nil {
		//데이터베이스를 사용하면
		//초기화단계에서 blockchain을
		//database에서 가져올 수 있음
		//어떻게 생성되고, 관리되는지
		//b는 한번만 초기화 됨(nil일때)
		once.Do(func() {
			b = &blockchain{"", 0}
			// search for checkpoint on the db
			checkpoint := db.GetBlockchain()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")

			} else {
				// restore b from bytes
				b.fromBytes(checkpoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}
