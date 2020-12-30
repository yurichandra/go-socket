package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWS open up websocket connection
func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Fail to open WS connection")
	}

	for {
		messageType, p, err := conn.ReadMessage()
		fmt.Println(p)
		if err != nil {
			fmt.Println(err)
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
		}
	}
}
