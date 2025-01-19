package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const addr = "localhost:8081"
const maxPlayers = 128

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return strings.HasPrefix(origin, "http://localhost")
	},
}

var connections = map[*websocket.Conn]time.Time{} // TODO loop over this every 10 seconds and purge connections that have expired

func StartServer() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func addConnection(conn *websocket.Conn) {
	connections[conn] = time.Now()
	log.Println("Adding connection", conn.RemoteAddr())
	log.Println(len(connections), "/", maxPlayers, " player slots in use")
}

func removeConnection(conn *websocket.Conn) {
	conn.Close()
	delete(connections, conn)
	log.Println("Removing connection", conn.RemoteAddr())
	log.Println(len(connections), "/", maxPlayers, " player slots in use")
}

func handle(writer http.ResponseWriter, req *http.Request) {
	if len(connections) >= maxPlayers {
		// TODO should we send something useful instead of ignoring the connection?
		return
	}

	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		return
	}

	addConnection(conn)

	defer removeConnection(conn)

	for {
		mt, msg, err := conn.ReadMessage()

		if mt != websocket.BinaryMessage {
			return
		}

		if err != nil {
			log.Println("Failed to read message:", err)
			return
		}

		action := &PlayerAction{}
		if err := proto.Unmarshal(msg, action); err != nil {
			log.Println("Failed to parse action:", err)
			return
		}

		log.Println("Received action:", action)
		connections[conn] = time.Now()
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
