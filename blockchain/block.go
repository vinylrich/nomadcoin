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
func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.Encoding(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	fmt.Println("Mining")
	for {
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
	fmt.Println("Mining End")
}
func createBlock(prevHash string, difficulty, height int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}

	block.mine()
	block.Transactions = Mempool.TxToConFirm()
	persistBlock(block)
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
