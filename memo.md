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

## 10.9

트랜잭션이 하나 생길때마다 그 전의 txOut은 spent Txout이 되고,
새로운 트랜잭션의 input은 그 전 txOut이 된다

//transaction이 발생할 때
	//case 1 : 5원을 줘야하는데 5개가 있을 때
	//1 transaction 1 input 1 output

	//case 2 : 5원을 줘야하는데 10개가 있을 때
	//잔돈을 거슬러 줘야함 -> 2개의 output 발생시킴
	//이런 조건 하에 아무리 많아도 2개의 output만을 가질 수 있음