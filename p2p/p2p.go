package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/utils"
	"time"

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
		if err != nil {
			wsConn.Close()
			break
		}
		fmt.Printf("Just Got:%s\n\n", p)
		time.Sleep(5 * time.Second)
		message := fmt.Sprintf("New message: %s", p)
		utils.HandleError(wsConn.WriteMessage(websocket.TextMessage, []byte(message)))
	}

}
