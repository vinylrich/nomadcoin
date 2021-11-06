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

type UTxOut struct {
	TxID   string
	Index  int
	amount int
}
type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

func makeConinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "Coinbase"},
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

func makeTx(from, to string, amount int) (*Tx, error) {

}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("junwoo", to, amount)
	if err != nil {
		return err
	}

	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConFirm() []*Tx {
	//mempool에 있는 모든 transaction을 실제 transaction에 넣음
	//mempool에 있는 데이터는 다 지움
	coinbase := makeConinbaseTx("junwoo")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
