package handlers

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	clients = make(map[string]map[*websocket.Conn]bool)
	mutex   = sync.Mutex{}
)

func HandleWebSocketConnection(c *websocket.Conn) {
	defer c.Close()

	fileId := c.Params("fileId")

	mutex.Lock()
	if clients[fileId] == nil {
		clients[fileId] = make(map[*websocket.Conn]bool)
	}
	clients[fileId][c] = true
	mutex.Unlock()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		broadcastMessage(fileId, c, msg)
	}

	mutex.Lock()
	delete(clients[fileId], c)
	mutex.Unlock()
}

func broadcastMessage(fileId string, sender *websocket.Conn, msg []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range clients[fileId] {
		if conn != sender {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				conn.Close()
				delete(clients[fileId], conn)
			}
		}
	}
}
