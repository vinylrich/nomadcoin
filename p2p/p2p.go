package p2p

import (
	"fmt"
	"net/http"
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

	fmt.Println(r.RemoteAddr)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	initPeer(wsConn, ipaddress, openPort)

}

// Port :4000 is reqeuesting an upgrade form the port :3000
func AddPeer(address, port, openPort string) {
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleError(err)
	peer := initPeer(wsConn, address, port) //conntect peer
	sendNewestBlock(peer)
	//결론은 inbox에 block을 넣어주는 것
}

/*
	Peer["127.0.0.1:3000"]=conn
*/
