package blockchain

import (
	"errors"
	"nomadcoin/utils"
	"nomadcoin/wallet"
	"sync"
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

//save all the signature to Txinput
//signature confirmed by public key
//public keys are also TxOut address
func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, *wallet.Wallet())
	}
}

func validate(t *Tx) bool {
	valid := true
	for _, txIn := range t.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(t.Id, txIn.Signature, address)
		if !valid {
			break
		}
	}
	return valid
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

//TxOut is only used by index of TxOut
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"` //signature를 만들때
	//모든 input에 서명을 함.
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var ErrorNoMoney = errors.New("not Enough money")
var ErrorTxInVaild = errors.New("Tx Invalid")
var ErrorzeroNotAllowed = errors.New("amount zero is invaild")

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
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

func isOnMempool(uTxOut *UTxOut) (exists bool) {
	exists = false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return
}

func makeTx(sender, receiver string, amount int) (*Tx, error) {
	if amount == 0 {
		return nil, ErrorzeroNotAllowed
	}
	if BalanceByAddress(sender, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(sender, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, sender}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{sender, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{receiver, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)

	if !valid {
		return nil, ErrorTxInVaild
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}

	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *mempool) TxToConFirm() []*Tx {
	//mempool에 있는 모든 transaction을 실제 transaction에 넣음
	//mempool에 있는 데이터는 다 지움
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx) //make new empty map
	return txs
}
func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.Id] = tx
}
