package blockchain

import (
	"nomadcoin/utils"
	"time"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOut"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string
	Amount int //to give miner
}

type TxOut struct {
	Owner  string
	Amount int
}

func makeConinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"Coinbase", minerReward},
	}
	txOut := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOut,
	}
	tx.getId()
	return &tx
}
