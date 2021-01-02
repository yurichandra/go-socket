package main

import "fmt"

// Channel :nodoc
type Channel struct {
	client     map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newChannel() *Channel {
	return &Channel{
		client:     make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (c *Channel) run() {
	for {
		select {
		case client := <-c.register:
			c.client[client] = true
			fmt.Printf("%s entered the channel", client.name)
		case client := <-c.unregister:
			delete(c.client, client)
			close(client.message)
			fmt.Printf("%s leave the channel", client.name)
		case message := <-c.broadcast:
			for client := range c.client {
				client.message <- message
			}
		}
	}
}
