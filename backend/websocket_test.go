// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketConnection tests the WebSocket connection functionality.
// It verifies that a client can successfully connect to the WebSocket endpoint
// and that the server correctly adds the client to the clients map.
func TestWebSocketConnection(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	s := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	mutex.Lock()
	assert.Equal(t, len(clients), 1, "One client should be connected")
	mutex.Unlock()
}

// TestBroadcastCount tests the BroadcastCount function.
// It verifies that when BroadcastCount is called, all connected clients
// receive the updated count via WebSocket.
func TestBroadcastCount(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	s := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer conn2.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	testCount := 69
	expectedMsg := fmt.Sprintf("<span id=\"counter\">%d</span>", testCount)

	readMessage := func(conn *websocket.Conn) {
		defer wg.Done()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("Error reading message: %v", err)
			return
		}

		assert.Equal(t, expectedMsg, string(msg), "Client receive the correct broadcast message")
	}

	go readMessage(conn1)
	go readMessage(conn2)

	BroadcastCount(testCount)

	wg.Wait()
}

// TestClientDisconnection tests that disconnected clients are removed from the clients map.
// It verifies that when a client closes its connection, it is properly removed from the map.
func TestClientDisconnection(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	s := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}

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
func TestMultipleBroadcasts(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	s := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	expectedMsgs := []string{
		"<span id=\"counter\">1</span>",
		"<span id=\"counter\">2</span>",
		"<span id=\"counter\">3</span>",
	}

	msgChan := make(chan string, len(expectedMsgs))

	go func() {
		for i := range len(expectedMsgs) {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				t.Errorf("Error reading message %d: %v", i, err)
				return
			}

			msgChan <- string(msg)
		}
	}()

	BroadcastCount(1)
	BroadcastCount(2)
	BroadcastCount(3)

	var receivedMsgs []string
	for i := range len(expectedMsgs) {
		select { // ? wtf does this do?
		case msg := <-msgChan:
			receivedMsgs = append(receivedMsgs, msg)
		case <-time.After(2 * time.Second):
			t.Fatalf("Timeout waiting for mmessage %d", i)
		}
	}

	assert.Equal(t, expectedMsgs, receivedMsgs, "All broadcast messages should be received in the correct order")
}
