package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func handle(writer http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		mt, msg, err := conn.ReadMessage()

		if mt != websocket.BinaryMessage {
			break
		}

		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		action := &PlayerAction{}
		if err := proto.Unmarshal(msg, action); err != nil {
			log.Println("Failed to parse address book:", err)
			break
		}

		log.Println("Received action:", action)
		action.Act(conn)
	}
}

func send(conn *websocket.Conn, update protoreflect.ProtoMessage) error {
	out, err := proto.Marshal(update)
	if err != nil {
		log.Println("Failed to encode update:", err)
		return err
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
		log.Println("Failed to write message:", err)
		return err
	}

	return nil
}

// TODO broadcasts
