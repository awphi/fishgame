package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const addr = "localhost:8081"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return strings.HasPrefix(origin, "http://localhost")
	},
}

func StartServer() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()

		if mt != websocket.BinaryMessage {
			break
		}

		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		action := &PlayerAction{}
		if err := proto.Unmarshal(message, action); err != nil {
			log.Println("Failed to parse address book:", err)
			break
		}

		log.Println("Received action:", action)
		action.Act(c)
	}
}

func send(c *websocket.Conn, update *GameServerUpdate) error {
	out, err := proto.Marshal(update)
	if err != nil {
		log.Println("Failed to encode update:", err)
		return err
	}

	if err := c.WriteMessage(websocket.BinaryMessage, out); err != nil {
		log.Println("Failed to write message:", err)
		return err
	}

	return nil
}

// TODO broadcasts
