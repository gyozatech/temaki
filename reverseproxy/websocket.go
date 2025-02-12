package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// functions to manage websockets and upgrade the http/https connection to ws/wss

func isWebSocketRequest(r *http.Request) bool {
	connectionHeader := strings.ToLower(r.Header.Get("Connection"))
	upgradeHeader := strings.ToLower(r.Header.Get("Upgrade"))

	return strings.Contains(connectionHeader, "upgrade") && upgradeHeader == "websocket"
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, targetURL string) {
	remote, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL %s", err), http.StatusBadGateway)
		return
	}

	// Ensure correct WebSocket scheme
	/*if strings.HasPrefix(remote.Scheme, "http") {
		remote.Scheme = strings.Replace(remote.Scheme, "http", "ws", 1)
	}*/

	// Upgrade client connection to WebSocket
	clientConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer clientConn.Close()

	// Dial WebSocket server
	serverConn, _, err := websocket.DefaultDialer.Dial(remote.String(), nil)
	if err != nil {
		log.Println("Failed to connect to WebSocket server:", err)
		return
	}
	defer serverConn.Close()

	// Start bidirectional message forwarding
	errCh := make(chan error, 2)

	// Client -> Server
	go func() {
		errCh <- copyWebSocketMessages(clientConn, serverConn)
	}()

	// Server -> Client
	go func() {
		errCh <- copyWebSocketMessages(serverConn, clientConn)
	}()

	// Wait for errors
	<-errCh
}

func copyWebSocketMessages(src, dest *websocket.Conn) error {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			return err
		}

		err = dest.WriteMessage(messageType, message)
		if err != nil {
			return err
		}
	}
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for security)
	},
}
