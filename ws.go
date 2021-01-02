package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client :nodoc
type Client struct {
	channel *Channel
	name    string
	conn    *websocket.Conn
	message chan []byte
}

// serveWS open up websocket connection
func serveWS(channel *Channel, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Fail to open WS connection")
	}

	name := r.URL.Query()["name"]

	client := &Client{
		channel: channel,
		conn:    conn,
		name:    name[0],
		message: make(chan []byte, 256),
	}

	client.channel.register <- client

	go client.read()
	go client.write()
}

// read means read message from websocket connection and
// distribute it to channel.
func (c *Client) read() {
	defer func() {
		c.channel.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				return
			}

			fmt.Print("here")
			fmt.Printf("Fail to read message: %s", err)
		}

		c.channel.broadcast <- message
	}
}

// write means write message from channel and send it
// to websocket connection
func (c *Client) write() {
	for {
		select {
		case message, _ := <-c.message:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				if strings.Contains(err.Error(), "websocket: close") {
					return
				}

				fmt.Printf("Fail to write message: %s", err)
			}

			w.Write(message)

			n := len(c.message)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.message)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
