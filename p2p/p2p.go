package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	for {
		_, p, err := wsConn.ReadMessage()
		utils.HandleError(err)
		fmt.Printf("%s\n\n", p)
	}
}
