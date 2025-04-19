// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
// It allows all origins for simplicity, but should be restricted in production.
// ! Restrict this in production
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Single websocket connection for simple testing.
var activeConn *websocket.Conn

// BroadcastCount sends the current count to the active websocket connection
func BroadcastCount(count int) {
	if activeConn != nil {
		message := fmt.Sprintf("<span id=\"counter\">%d</span>", count)

		if err := activeConn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("Error updating count with websockets: ", err)
		}
	}
}

// wsHandler handles WebSocket connections from clients.
// It upgrades the HTTP connection to a WebSocket connection, listens for messages,
// and echoes them back to the client.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket: ", err)
		return
	}
	// Save the connection
	activeConn = conn
	fmt.Println("New websocket client connected")

	// Handle disconnection
	defer func() {
		if activeConn == conn {
			activeConn = nil
		}
		conn.Close()
		fmt.Println("Websocket client disconnected")
	}()

	// Listen for incoming messages
	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}
		fmt.Printf("Received message: %s\\n", message)
	}
}
