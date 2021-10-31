package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
	Height   int
}

func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleError(encoder.Encode(b))
	return blockBuffer.Bytes()
}
func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
}
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)

	block.Hash = fmt.Sprintf("%s", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
