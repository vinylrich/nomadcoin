## 10.3 Balances

1. 모든 거래 출력값 받기 (txOuts)
2. address로 받기
3. input 받기


## 10.7 confirm Transcation
```go
func (m *mempool) txToConFirm()[]*Tx {
	coinbase:=makeConinbaseTx("junwoo")//
}
```

junwoo라고 써져 있는 부분 wallet만들기

Tx1
    TxIns[COINBASE]
    TxOuts[$5(you)] <---- Spent TxOut
Tx2
    TxIn[Tx1.TxOuts[0]
    TxOuts[$5(me)] <---- Spent TxOut

Tx3
    TxIns[Tx2.TxOuts[0]]
    TxOuts[$3(you), $2(me)] <----- uTxOut x2