package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/utils"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var wsConns []*websocket.Conn

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
	time.Sleep(20 * time.Second)
	wsConn.WriteMessage(websocket.TextMessage, []byte("Hello from Port 3000!"))
}

// Port :4000 is reqeuesting an upgrade form the port :3000
func AddPeer(address, port, openPort string) {
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleError(err)
	initPeer(wsConn, address, port) //conntect peer
	time.Sleep(10 * time.Second)
	wsConn.WriteMessage(websocket.TextMessage, []byte("Hello from Port 4000!"))
}

/*
	Peer["127.0.0.1:3000"]=conn
*/
