package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/awphi/fishgame/fish"
	"github.com/awphi/fishgame/gen"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const addr = "localhost:8081"
const tickRate = 500 // ms/t

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return strings.HasPrefix(origin, "http://localhost")
	},
}

func main() {
	f, err := fish.GenerateFish(0)
	fmt.Println(f.Type, err)

	// TODO goroutine with a ticker to update DB instance every 15 mins or so?

	go startServer()
	// game loop runs on the main thread so it blocks
	gameLoop()
}

func startServer() {
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

		// TODO pull this IO stuff out - make utils for wrapping updates and unwrapping actions and finish implementing ID-based ping

		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		action := &gen.PlayerAction{}
		if err := proto.Unmarshal(message, action); err != nil {
			log.Println("Failed to parse address book:", err)
			break
		}

		log.Println("Received action:", action)

		// for now just always respond with a ping with ID 1008
		payload := &gen.GameServerUpdatePong{Id: 1008}
		update := &gen.GameServerUpdate{Payload: &gen.GameServerUpdate_Pong{Pong: payload}}

		out, err := proto.Marshal(update)
		if err != nil {
			log.Println("Failed to encode update:", err)
		}

		if err := c.WriteMessage(websocket.BinaryMessage, out); err != nil {
			log.Println("Failed to write message:", err)
			break
		}
	}
}

func gameLoop() {
	ticker := time.NewTicker(tickRate * time.Millisecond)
	defer ticker.Stop()
	for tickTime := range ticker.C {
		tick(tickTime)
	}
}

func tick(t time.Time) {
	//fmt.Println("tick", t)
}
