package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

type peers struct {
	Value map[string]*peer
	mutex *sync.Mutex
}

var Peers peers = peers{
	Value: make(map[string]*peer),
	mutex: &sync.Mutex{},
}

func AllPeers(p *peers) []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var keys []string

	for key := range p.Value {
		keys = append(keys, key)
	}
	return keys
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.mutex.Lock()
	defer Peers.mutex.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}
	go p.read()
	go p.write()
	Peers.Value[key] = p
	return p
}

func (p *peer) close() {
	Peers.mutex.Lock()
	defer Peers.mutex.Unlock()
	p.conn.Close()
	delete(Peers.Value, p.key)

}

func (p *peer) read() {
	defer p.close()
	//delete peer if error causes
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		//Decoding Json to Go
		if err != nil {
			break
		}
		handleMsg(&m, p) //Print msg
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		//write 함수에서는 messages/sendBlock()
		//function에서 준 block byte inbox값을 받음
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}
