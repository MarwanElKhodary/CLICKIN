// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type TestSuite struct {
	Db      *sql.DB
	SqlMock sqlmock.Sqlmock
	Router  *gin.Engine
	Repo    *Repository
	Serv    *Service
	Handler *Handler

	WsServer  *httptest.Server
	WsUrl     string
	WsClients []*websocket.Conn
	WsMutex   sync.Mutex
}

var testSuite *TestSuite

// InitWebSocketServer initializes the WebSocket test server
func (ts *TestSuite) InitWebSocketServer() {
	ts.WsServer = httptest.NewServer(http.HandlerFunc(wsHandler))
	// Convert http://127.0.0.1 to ws://127.0.0.
	ts.WsUrl = "ws" + strings.TrimPrefix(ts.WsServer.URL, "http")
	ts.WsClients = make([]*websocket.Conn, 0)
}

// AddWsClient creates and returns new WebSocket client connection
// The new client is added to TestSuite.WsClients
func (ts *TestSuite) AddWsClient() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(ts.WsUrl, nil)
	if err != nil {
		panic("Failed to create a mock database: " + err.Error())
	}

	ts.WsMutex.Lock()
	ts.WsClients = append(ts.WsClients, conn)
	ts.WsMutex.Unlock()

	return conn
}

// TestMain is called once before any tests are run
func TestMain(m *testing.M) {
	var teardown func()
	testSuite, teardown = setupTestSuite()

	code := m.Run()

	teardown()

	os.Exit(code)
}

// setupTestSuite initializes test dependencies and returns a TestSuite and teardown function.
// It creates a mock database, repository, service, and handler, and sets up routes.
func setupTestSuite() (*TestSuite, func()) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		panic("Failed to create a mock database: " + err.Error())
	}

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	router := gin.Default()
	handler.SetupRoutes(router)

	suite := &TestSuite{
		Db:      db,
		SqlMock: mockSql,
		Router:  router,
		Repo:    repo,
		Serv:    service,
		Handler: handler,
	}

	suite.InitWebSocketServer()

	return suite, func() {
		db.Close()
	}
}
