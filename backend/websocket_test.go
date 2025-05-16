// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketConnection tests the WebSocket connection functionality.
// It verifies that a client can successfully connect to the WebSocket endpoint
// and that the server correctly adds the client to the clients map.
func TestWebSocketConnection(t *testing.T) {
	fmt.Printf("testWebSocketConnection: %d\n", len(clients))
	client := testSuite.CreateWsClient()
	defer client.Close()

	mutex.Lock()
	assert.Equal(t, 1, len(clients), "One client should be connected")
	mutex.Unlock()
	//TODO: Add a testSuite.ClearWsClients() here
}

// TestBroadcastCount tests the BroadcastCount function.
// It verifies that when BroadcastCount is called, all connected clients
// receive the updated count via WebSocket.
func TestBroadcastCount(t *testing.T) {
	fmt.Printf("testBroadcastCount: %d\n", len(clients))
	connOne := testSuite.CreateWsClient()
	connTwo := testSuite.CreateWsClient()
	defer connOne.Close()
	defer connTwo.Close()

	testCount := 69
	expectedMsg := fmt.Sprintf("<span id=\"counter\">%d</span>", testCount)
	msgChan := make(chan string, 2)

	readMessage := func(conn *websocket.Conn) {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("Error reading message: %v", err)
			return
		}
		msgChan <- string(msg)
	}

	go readMessage(connOne)
	go readMessage(connTwo)

	BroadcastCount(testCount)

	received := []string{}
	for i := range 2 {
		select {
		case msg := <-msgChan:
			received = append(received, msg)
		case <-time.After(2 * time.Second):
			t.Fatalf("Timeout waiting for message %d", i)
		}
	}

	assert.Contains(t, received, expectedMsg)
}

// TestClientDisconnection tests that disconnected clients are removed from the clients map.
// It verifies that when a client closes its connection, it is properly removed from the map.
func TestClientDisconnection(t *testing.T) {
	fmt.Printf("testClientDisconnection: %d\n", len(clients))
	conn := testSuite.CreateWsClient()

	mutex.Lock()
	initialClientCount := len(clients)
	mutex.Unlock()

	conn.Close()

	mutex.Lock()
	assert.Equal(t, initialClientCount-1, len(clients), "Client should be removed from clients map after disconnection")
	mutex.Unlock()
}

// TestMultipleBroadcasts tests that multiple broadcasts work correctly.
// It verifies that when multiple BroadcastCount calls are made in succession,
// all clients receive all messages in the correct order.
// func TestMultipleBroadcasts(t *testing.T) {
// 	teardownTestCase := setupTestCase(t)
// 	defer teardownTestCase(t)

// 	s := httptest.NewServer(http.HandlerFunc(wsHandler))
// 	defer s.Close()

// 	// Convert http://127.0.0.1 to ws://127.0.0.
// 	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		t.Fatalf("Could not connect to WebSocket server: %v", err)
// 	}
// 	defer conn.Close()

// 	expectedMsgs := []string{
// 		"<span id=\"counter\">1</span>",
// 		"<span id=\"counter\">2</span>",
// 		"<span id=\"counter\">3</span>",
// 	}

// 	msgChan := make(chan string, len(expectedMsgs))

// 	go func() {
// 		for i := range len(expectedMsgs) {
// 			_, msg, err := conn.ReadMessage()
// 			if err != nil {
// 				t.Errorf("Error reading message %d: %v", i, err)
// 				return
// 			}

// 			msgChan <- string(msg)
// 		}
// 	}()

// 	BroadcastCount(1)
// 	BroadcastCount(2)
// 	BroadcastCount(3)

// 	var receivedMsgs []string
// 	for i := range len(expectedMsgs) {
// 		select {
// 		case msg := <-msgChan:
// 			receivedMsgs = append(receivedMsgs, msg)
// 		case <-time.After(2 * time.Second):
// 			t.Fatalf("Timeout waiting for mmessage %d", i)
// 		}
// 	}

// 	assert.Equal(t, expectedMsgs, receivedMsgs, "All broadcast messages should be received in the correct order")
// }
