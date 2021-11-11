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

트랜잭션 input은 output을 찾기 위한 표지판

//함수가 구조체를 변화시킨다면 함수는 메서드여야함
//struct가 변화하지 않는다면 function

1. 서명 검증 -공개키,비공개키
2. wallet persistency
3. 서명, 증명 후 트랜잭션에 적용


1) we hash the msg.
"I LOVE YOU" -> hash() -> "hashed_message"

2) generate key pair
keyPair (privateK,publicK) wallet->(private key -> file)

3) sign the hash

("hashed_message" + privateK) -> "signature" //privateK로는 서명
//public K로는 검증

4) verify

("hased_message +"signature" + publicK) -> true / false

```go
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"owner"` //signature를 만들때
	//모든 input에 서명을 함.
}

type TxOut struct {
	Address string `json:"owner"`
	Amount  int    `json:"amount"`
}
```
txin을 가져와 우리가 할 수 있는건 txout address는 txinput의 서명을 검증할 수 있어야 함

1. 우리가 tx에 서명할 때 wallet을 private key로 서명함
2. 근데 우리가 거짓말을 하고 있을 수도 있음
3. tx input을 만들기 위해 필요한 tx output을 소유하지 않을 가능성이 있음
4. validate에서 하고 있는 건 txout의 address 즉, public key를 가지고 오는 것


그니까 트랜잭션 아웃풋은 트랜잭션 인풋을 참조하고 있음

인풋에는 private key로 만들어진 signature가 있고 txinput이 tx output을 검증하기 위해서는 txout의 address인 public key로 private key를 검증해야함 


사람들은 나의 public key(address)에 돈을 보내고

내가 그 돈을 사용할 때, 서명을 하게 됨

input은 돈을 사용한다는 거고, 내 private key로 서명할 것임

근데 사람들은 내가 그 돈을 가지고 있다는 것을 믿지 않음

근데 사람들은 나의 public key(address)에 돈을 보냈기 때문에

서명 안에 있는 private key와 public key를 사용해서 돈이 있는지 없는지 검증함


TxOut1,TxOut2

Tx
    TxIn[
        (TxOut1) has address - public key
        (TxOut2) has address = public key
    ]
    Sign: X

    verify: TxIn.Sign(using private key) + TxOut1.Address =true/false


p2p

우리가 컴퓨터(server)를 끄면 database는 살아있겠지만 블록체인 네트워크는 꺼짐 


socket을 사용해서 실시간으로 서로를 업데이트 

ex) 5개 노드가 작동하고 있으면

첫 노드가 새로운 block을 찾으면 그 새로운 block이 다른 4개의 노드에게 알리고, 4개의 노드는 상태를 바꿔야 함

3번째 node를 메모리 영역에서 tx하면 memory영역은 모든 네트워크에 중계되어야함

이 시스템은 동시에 실행되기 때문에 channel을 사용할 때 race condition(경쟁 상태) 라는 버그가 생겨날 수 있음