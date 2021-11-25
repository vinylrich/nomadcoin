package p2p

import (
	"encoding/json"
	"fmt"
	"nomadcoin/blockchain"
	"nomadcoin/utils"
	"strings"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResPonse
	MessageBroadCastNewBlock
	MessageBroadCastNewTx
	MessageBroadCastNewPeer
)

//iota is serial
type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}
func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleError(err)
	//msg의 payload
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m //블록의 byte 정보를 inbox에 넣어줌
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResPonse, blockchain.AllBlocks(blockchain.Blockchain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageBroadCastNewBlock, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageBroadCastNewTx, tx)
	p.inbox <- m
}
func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageBroadCastNewPeer, address)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		fmt.Printf("Received the newest block from %s\n", p.key)
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleError(err)
		//if not a.block equals to b.block
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else { //if 3000 hasnt Block
			//4000 -> 3000
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p)
			//send 4000 out block
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResPonse:
		fmt.Printf("Receive all the blocks from %s", p.key)
		var payload []*blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)

	case MessageBroadCastNewBlock:
		var payload *blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageBroadCastNewTx:
		var payload *blockchain.Tx
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageBroadCastNewPeer:
		var payload string
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}

}
