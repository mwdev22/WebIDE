package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Connection *websocket.Conn
	UserId     string
}

var (
	clients = make(map[string]map[*websocket.Conn]Client)
	mutex   = sync.Mutex{}
)

func HandleWebSocketConnection(c *websocket.Conn) {
	defer c.Close()

	fileId := c.Params("fileId")
	userId := c.Query("userId")

	mutex.Lock()
	if clients[fileId] == nil {
		clients[fileId] = make(map[*websocket.Conn]Client)
	}
	clients[fileId][c] = Client{Connection: c, UserId: userId}
	mutex.Unlock()

	for {
		_, fileContent, err := c.ReadMessage()
		if err != nil {
			break
		}
		go broadcastMessage(fileId, c, fileContent)
	}

	mutex.Lock()
	delete(clients[fileId], c)
	mutex.Unlock()
}

func broadcastMessage(fileId string, sender *websocket.Conn, msg []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn, client := range clients[fileId] {
		if conn != sender {
			message := map[string]interface{}{
				"fileId": fileId,
				"userId": client.UserId,
				"data":   string(msg),
			}

			messageJSON, err := json.Marshal(message)
			if err != nil {
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				conn.Close()
				delete(clients[fileId], conn)
			}
		}
	}
}
