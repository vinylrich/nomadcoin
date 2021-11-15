package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/utils"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var wsConns []*websocket.Conn

// Port :3000 will upgrade the reqeust from :4000
func Upgrade(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	fmt.Print(r.RemoteAddr)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	openPort := r.URL.Query().Get("openPort")

	sepedstr := strings.Split(r.RemoteAddr, ":")
	initPeer(wsConn, sepedstr[0], openPort)
}

// Port :4000 is reqeuesting an upgrade form the port :3000
func AddPeer(address, port, openPort string) {
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleError(err)
	initPeer(wsConn, address, port)
}

/*
	Peer["127.0.0.1:3000"]=conn
*/
