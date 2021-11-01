package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
	"strings"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
}

const difficulty int = 2

func (b *Block) fromBytes(data []byte) {
	utils.FromBytes(b, data)
}
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		strBlock := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(strBlock)))
		fmt.Printf("Block string: %s\nHash:%s\nNonce:%d\n", strBlock, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
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
