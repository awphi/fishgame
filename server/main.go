package main

import (
	"fmt"
	"time"

	"github.com/awphi/fishgame/fish"
)

const tickRate = 500 // ms/t

func main() {
	f, err := fish.GenerateFish(0)
	fmt.Println(f.Type, err)

	// TODO goroutine for websockets (either build a list of messages for the game loop to handle or handle directly, locking via mutexes)
	// TODO goroutine with a ticker to update DB instance every 15 mins or so?

	// game loop runs on the main thread so it blocks
	gameLoop()
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
