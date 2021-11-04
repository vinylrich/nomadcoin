package blockchain

import (
	"errors"
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
	"strings"
	"time"
)

type Block struct {
	Transactions []*Tx  `json:"transactions"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
}

func (b *Block) fromBytes(data []byte) {
	utils.Decoding(b, data)
}
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.Encoding(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		fmt.Printf("Target:%s\nHash:%s\nNonce:%d\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}
func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:         "",
		PrevHash:     prevHash,
		Height:       height,
		Difficulty:   Blockchain().difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeConinbaseTx("junwoo")},
	}

	block.mine()
	block.persist()
	return block
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}

	block := &Block{}
	block.fromBytes(blockBytes)

	return block, nil
}
