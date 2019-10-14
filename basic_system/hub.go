package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	doc        string
}

type Message struct {
	content []byte
	sender  *Client
}

var dmp *diffmatchpatch.DiffMatchPatch = diffmatchpatch.New()

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		doc:        "",
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			err := client.conn.WriteMessage(websocket.TextMessage, []byte(h.doc))
			if err != nil {
				return
			}

			fmt.Println("Registered a client")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			patches, err := dmp.PatchFromText(string(message.content))
			if err != nil {
				return
			}
			h.doc, _ = dmp.PatchApply(patches, h.doc)

			for client := range h.clients {
				if client != message.sender {
					client.send <- message.content
				}
			}
		}
	}
}
