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

//아직 input에서 사용되지 않은 output
//address와 똑같은 owner 중에서 txid와 input id와 같은 것 중에서
//false인지 true인지 판별하여 이게 output이 이미 있는지 확인하는 것
func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	//transaction이 발생할 때
	//case 1 : 5원을 줘야하는데 5개가 있을 때
	//1 transaction 1 input 1 output

	//case 2 : 5원을 줘야하는데 10개가 있을 때
	//잔돈을 거슬러 줘야함 -> 2개의 output 발생시킴
	//이런 조건 하에 아무리 많아도 2개의 output만을 가질 수 있음

	var uTxOuts []*UTxOut
	//아직 사용되지 않은 output
	createrTxs := make(map[string]bool)
	for _, block := range b.AllBlocks() { //모든 블록 불러옴
		for _, tx := range block.Transactions { //블록 안의 모든 트랜잭션
			for _, input := range tx.TxIns { //트랜잭션 INPUT
				if input.Owner == address { //
					//output을 생성하지 않은 input을 찾아야함
					createrTxs[input.TxID] = true //아웃풋을 사용한
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := createrTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.UTxOutsByAddress(address)
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
