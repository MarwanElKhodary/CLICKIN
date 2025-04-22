// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
