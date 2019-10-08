package main

import (
	"encoding/json"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	doc       string
	dmp       *diffmatchpatch.DiffMatchPatch
	patchChan chan []diffmatchpatch.Patch
	broker    *Broker
)

type RequestBody struct {
	Patch string `json:"patch"`
}

type EventStream struct {
	Data string `json:"data"`
}

// A Broker holds open client connections,
// listens for incoming events on its Notifier channel
// and broadcast event data to all registered connections
type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func init() {
	doc = ""
	dmp = diffmatchpatch.New()
	patchChan = make(chan []diffmatchpatch.Patch)
	broker = NewBroker()
}

func NewBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

// listen blocks on the different channels of Broker waiting for events
func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:
			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

// subscribeHandler registers a new connection and maintains it to send events
// back to it
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if server sent events are supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's
	// connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	// Block on messageChan to wait for next event
	for m := range messageChan {
		fmt.Fprintf(w, "data:%s\n\n", m)
		// Clear buffer
		flusher.Flush()
	}
}

// patchHandler parses patches sent from a client and applies it to the global
// doc
func patchHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	patchBody := RequestBody{}
	err = json.Unmarshal(reqBody, &patchBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	patches, err := dmp.PatchFromText(patchBody.Patch)
	if err != nil {
		http.Error(w, "Failed to parse patches", http.StatusBadRequest)
		return
	}

	// Send patches to patchChan to be applied to doc
	patchChan <- patches

	// Send to other connections
	broker.Notifier <- reqBody
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, doc)
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", fs)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/patch/", patchHandler)
	http.HandleFunc("/state/", stateHandler)

	go func() {
		for patches := range patchChan {
			doc, _ = dmp.PatchApply(patches, doc)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
