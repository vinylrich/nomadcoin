package p2p

import (
	"encoding/json"
	"fmt"
	"nomadcoin/blockchain"
	"nomadcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock       MessageKind = iota
	MessageAllBlocksRequest  MessageKind = iota
	MessageAllBlocksResPonse MessageKind = iota
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
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m //블록의 byte 정보를 inbox에 넣어줌
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		fmt.Printf("Peer: %s Sent a message with kind of: %d", p.key, m.Kind)
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
	}

}
