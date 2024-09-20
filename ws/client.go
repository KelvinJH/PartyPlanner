package ws

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MAX_ROOMS = 15
)

type ChatClient struct {
	connection *websocket.Conn
	server     *MessageManager
	egress     chan []byte
}

type EventClient struct {
	connection *websocket.Conn
	server     *EventManager
	eventQueue chan []byte
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

func NewChatClient(connection *websocket.Conn, server *MessageManager) *ChatClient {
	return &ChatClient{
		connection: connection,
		server:     server,
		egress:     make(chan []byte),
	}
}

func NewEventClient(connection *websocket.Conn, server *EventManager) *EventClient {
	return &EventClient{
		connection: connection,
		server:     server,
		eventQueue: make(chan []byte),
	}
}

func (client *ChatClient) readMessages() {
	defer func() {
		// Clean up
		client.connection.Close()
		client.server.removeWebsocket(client)
	}()

	for {
		_, content, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error while reading mesasge: ", err)
			}
			break
		}
		fmt.Println("Inside the read message function of golang")

		message := Message{
			content:   string(content),
			timestamp: time.Now().Local().String(),
		}

		client.server.messageBroadcast <- &message
	}
}

func (client *ChatClient) writeMessages() {
	defer func() {
		client.server.removeWebsocket(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					fmt.Println("Connection closed: ", err)
				}
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("failed to send message: ", err)
			}

		}

	}
}

func (client *EventClient) writeMessages() {
	defer func() {
		client.server.removeWebsocket(client)
	}()

	for {
		select {
		case message, ok := <-client.eventQueue:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					fmt.Println("Connection closed: ", err)
				}
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("failed to send event message: ", err)
			}
		}
	}
}
