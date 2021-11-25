package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/blockchain"
	"nomadcoin/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Port :3000 will upgrade the reqeust from :4000
func Upgrade(w http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ipaddress := utils.Splitter(r.RemoteAddr, ":", 0)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ipaddress != ""
	}
	fmt.Printf("%s to wants an upgrade\n", openPort)
	fmt.Println(r.RemoteAddr)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	initPeer(wsConn, ipaddress, openPort)

}

// Port :4000 is reqeuesting an upgrade form the port :3000
func AddPeer(address, port, openPort string, broadcast bool) {
	fmt.Printf("%s to connect to port %s\n", openPort, port)
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleError(err)
	//Starting initPeer function, read,write with go routine
	peer := initPeer(wsConn, address, port) //conntect peer
	if broadcast {
		broadcastNewPeer(peer)
		return
	}
	sendNewestBlock(peer)
	//결론은 inbox에 block을 넣어주는 것
}

/*
	Peer["127.0.0.1:3000"]=conn
*/

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.Value {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.Value {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.Value {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
