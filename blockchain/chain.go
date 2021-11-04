package blockchain

import (
	"nomadcoin/db"
	"nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currendifficulty"`
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowRange         int = 2
)

var b *blockchain
var once sync.Once

func (b *blockchain) decoding(data []byte) {
	utils.Decoding(b, data)
}
func (b *blockchain) persist() {
	db.SaveBlockchain(utils.Encoding(b))
}
func (b *blockchain) AllBlocks() []*Block {
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.AllBlocks()
	newestBlock := allBlocks[0]                              //첫 번째 블록
	lastRecalculatedBlock := allBlocks[difficultyInterval-1] //5번째 블록
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	//첫 번째 블록과 5번째 블록이 만들어진 시간을 뺌(분단위)
	expectedTime := difficultyInterval * blockInterval //10분마다 5개 단위수
	if actualTime < (expectedTime - allowRange) {      //actualTime이 8분보다 적으면
		//즉(처음 블록부터 끝 블록까지 만드는데 걸리는 시간이 8분 미만이면)
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowRange) {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 { //5개마다 재설정
		b.CurrentDifficulty = b.recalculateDifficulty()
	}
	return b.CurrentDifficulty

}

func (b *blockchain) AllTxOuts() []*TxOut {
	var txOuts []*TxOut
	blocks := b.AllBlocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	allTx := b.AllTxOuts()
	for _, txOut := range allTx { //모든 트렌젝션 인덱스별로 추출
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut) //owner가 address면 owntx에 추출한 개별 tx 넣음
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.TxOutsByAddress(address)
	var sumOwnedBalance int
	for _, txOut := range txOuts {
		sumOwnedBalance += txOut.Amount
	}
	return sumOwnedBalance
}

//singleton pattern
func Blockchain() *blockchain {
	if b == nil {
		//데이터베이스를 사용하면
		//초기화단계에서 blockchain을
		//database에서 가져올 수 있음
		//어떻게 생성되고, 관리되는지
		//b는 한번만 초기화 됨(nil일때)
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			// search for checkpoint on the db
			checkpoint := db.GetBlockchain()
			if checkpoint == nil {
				b.AddBlock()

			} else {
				// restore b from bytes
				b.decoding(checkpoint)
			}
		})
	}
	return b
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}
