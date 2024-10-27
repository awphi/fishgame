package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/awphi/fishgame/fish"
	"github.com/gorilla/websocket"
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
	// TODO goroutine for websockets (either build a list of messages for the game loop to handle or handle directly, locking via mutexes)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
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
	fmt.Println("tick", t)
}
