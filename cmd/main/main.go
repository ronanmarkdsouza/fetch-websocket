package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var connectedDevices int
var counterMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func printConnectedDevices() {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	fmt.Printf("Connected devices: %d\n", connectedDevices)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Increment the counter and add the client to the map
	counterMutex.Lock()
	connectedDevices++
	counterMutex.Unlock()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Print the updated count to the console
	printConnectedDevices()

	// Listen for client disconnection or errors
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// On disconnect or error, remove the client and decrement the counter
			counterMutex.Lock()
			connectedDevices--
			counterMutex.Unlock()

			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()

			// Print the updated count to the console
			printConnectedDevices()
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("WebSocket server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
