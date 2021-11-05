package blockchain

import (
	"errors"
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

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	//junwoo가 가지고 있는 balance 추출
	var txIns []*TxIn
	var txOuts []*TxOut
	oldTxOuts := Blockchain().TxOutsByAddress(from) //address가 가지고 있는 tx아웃풋 추출
	var total int
	for _, txOut := range oldTxOuts {
		if total > amount { //일정 금액보다 더 많이 채워지면
			//==이 아닌 > 를 쓰는 이유는 트랙잭션이 1,1,1,1 과 같이
			//나눠져서 amount가 있을 수 있기 때문이다
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount} // junwoo 50 100 50
		txIns = append(txIns, txIn)              // input에 update
		total += txIn.Amount                     //만약 20이면 스탑
	}

	change := total - amount //30
	if change != 0 {
		changeTxOut := &TxOut{Owner: from, Amount: change} //잔돈

		txOuts = append(txOuts, changeTxOut) //트렌잭션 새로 만듦
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		TxIns:     txIns,
		TxOuts:    txOuts,
		Timestamp: int(time.Now().Unix()),
	}
	tx.getId()
	return tx, nil
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
