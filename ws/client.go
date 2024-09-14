package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	MAX_ROOMS = 15
)

type Client struct {
	connection *websocket.Conn
	server     *Manager
	egress     chan []byte
}

type Message struct {
	clientName string
	content    string
	timestamp  string
}

type WSMessage struct {
	Headers interface{} `json:"HEADERS"`
	Message string      `json:"message"`
}

func NewClient(connection *websocket.Conn, server *Manager) *Client {
	return &Client{
		connection: connection,
		server:     server,
		egress:     make(chan []byte),
	}
}

func (client *Client) readMessages() {
	defer func() {
		// Clean up
		client.connection.Close()
		client.server.removeWebsocket(client)
	}()

	for {
		_, content, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error while reading mesasge: ", err)
			}
			break
		}
		log.Println("Inside the read message function of golang")

		message := Message{
			content:   string(content),
			timestamp: time.Now().Local().String(),
		}

		client.server.messageBroadcast <- &message
	}
}

func (client *Client) writeMessages() {
	defer func() {
		client.server.removeWebsocket(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ", err)
				}
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("failed to send message: ", err)
			}

		}

	}
}

func createTestMessages() []Message {
	messages := make([]Message, 3)

	for x := 0; x < len(messages); x++ {
		if x%2 == 0 {
			messages[x] = Message{
				content:   "This is a text message",
				timestamp: "9/6/2024",
			}
		} else {
			messages[x] = Message{
				content:   "This is a text reply",
				timestamp: "9/6/2024",
			}
		}
	}
	return messages
}
