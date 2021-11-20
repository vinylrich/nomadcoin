package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type peer struct {
	conn *websocket.Conn
}

var Peers map[string]*peer = make(map[string]*peer)

func initPeer(conn *websocket.Conn, address, port string) *peer {
	p := &peer{
		conn,
	}
	key := fmt.Sprintf("%s:%s", address, port)
	go p.read()
	Peers[key] = p
	return p
}

func (p *peer) read() {
	//delete peer if error causes
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}
