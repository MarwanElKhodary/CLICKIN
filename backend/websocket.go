// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//Link: https://medium.com/wisemonks/implementing-websockets-in-golang-d3e8e219733b
// TODO: Protect against CSWSH - i.e. add wss instead ws
// TODO: Use secure WebSockets
// TODO: Add maximum header to WebSockets

// upgrader is used to upgrade HTTP connections to WebSocket connections.
// It allows all origins for simplicity, but should be restricted in production.
// ! Restrict this in production
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// clients is a map that associates WebSocket connections with boolean values.
var clients = make(map[*websocket.Conn]bool)

// broadcast channel to all connected clients.
var broadcast = make(chan []byte)

// mutex protects the connected clients' map.
var mutex = &sync.Mutex{}

// BroadcastCount sends the current count to the broadcast.
func BroadcastCount(count int) {
	message := fmt.Sprintf("<span id=\"counter\">%d</span>", count)
	// ! Here should you being crypto it? Like hash and unshash.
	broadcast <- []byte(message)
}

// handleMessages broadcasts any updates to all clients.
func handleMessages() {
	for {
		message := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// wsHandler handles WebSocket connections from clients.
// It upgrades the HTTP connection to a WebSocket connection and handles the removal of dead clients.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket: ", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// While there is no need to process incoming messages from the client as it's handled by POST requests,
	// this detects when the client disconnects and handles the removal of the disconnected clients
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
	}
}

// init always gets called before main, which establishes all the goroutines
func init() {
	go handleMessages()
}
