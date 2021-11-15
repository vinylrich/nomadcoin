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
	Peers[key] = p
	return p
}
