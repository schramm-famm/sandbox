package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	ID    string `json:"id"`
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
	Notifier chan map[string][]byte

	// New client connections
	newClients chan map[string]chan []byte

	// Closed client connections
	closingClients chan string

	// Client connections registry
	clients map[string]chan []byte
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
		Notifier:       make(chan map[string][]byte),
		newClients:     make(chan map[string]chan []byte),
		closingClients: make(chan string),
		clients:        make(map[string]chan []byte),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

// listen blocks on the different channels of Broker waiting for events
func (broker *Broker) listen() {
	for {
		select {
		case client := <-broker.newClients:
			for id, c := range client {
				// A new client has connected.
				// Register their message channel
				broker.clients[id] = c
			}
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case id := <-broker.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, id)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:
			// Send event to all connected clients
			for senderID, msg := range event {
				for clientID, clientMessageChan := range broker.clients {
					if clientID != senderID {
						clientMessageChan <- msg
					}
				}
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
		log.Println("Streaming unsupported")
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

	var id uuid.UUID
	for {
		id = uuid.New()
		if _, ok := broker.clients[id.String()]; !ok {
			break
		}
	}

	// Signal the broker that we have a new connection
	broker.newClients <- map[string]chan []byte{id.String(): messageChan}

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- id.String()
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.closingClients <- id.String()
	}()

	idBody := RequestBody{ID: id.String()}
	idJSON, err := json.Marshal(idBody)
	if err != nil {
		log.Println("Unable to create ID for connection:", err)
		http.Error(w, "Unable to create ID for connection", http.StatusInternalServerError)
		return
	}
	// Send assigned id to client
	fmt.Fprintf(w, "data:%s\n\n", string(idJSON))
	// Clear buffer
	flusher.Flush()

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
	broker.Notifier <- map[string][]byte{patchBody.ID: reqBody}
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
