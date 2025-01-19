package server

import (
	"log"

	"github.com/gorilla/websocket"
)

func (msg *PlayerAction) Act(conn *websocket.Conn) {
	switch msg.GetPayload().(type) {
	case *PlayerAction_Ping:
		// Not particularly ergonomic but for reasons explained here (https://github.com/golang/protobuf/issues/1326) I *think* it's the best we can do
		send(conn, &GameServerUpdate{Payload: &GameServerUpdate_Pong{Pong: &GameServerUpdatePong{}}})
	default:
		log.Println("Unknown message payload")
	}
}
