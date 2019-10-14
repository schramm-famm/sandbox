package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func connect(h *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	client := &Client{c, h, make(chan []byte)}

	h.register <- client

	go client.read()
	go client.write()
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	h := newHub()
	go h.run()

	http.Handle("/", fs)
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		connect(h, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
