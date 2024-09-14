package ws

import (
	"fmt"
	"net/http"
	"sync"
	"partyplanner/service"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients          map[*Client]bool
	messages         []*Message
	messageBroadcast chan *Message
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:          make(map[*Client]bool),
		messageBroadcast: make(chan *Message),
	}
}

func (manager *Manager) addClient(client *Client) {

	fmt.Printf("Connection from client : %s\n", client.connection.RemoteAddr())
	// Lock the client map
	manager.Lock()
	// Unocks when it is done
	defer manager.Unlock()

	manager.clients[client] = true
}

func (manager *Manager) removeWebsocket(client *Client) {
	fmt.Println("Disconnecting client: ", client.connection.RemoteAddr())
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.clients[client]; ok {
		client.connection.Close()
		delete(manager.clients, client)
	}
}

func (manager *Manager) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error while connecting: ", err)
	}
	client := NewClient(connection, manager)
	manager.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (manager *Manager) Run() {
	for {
		select {
		case message := <-manager.messageBroadcast:
			manager.messages = append(manager.messages, message)
			for client := range manager.clients {
				select {
				case client.egress <- service.CreateChatMessages("Test User", message.content, message.timestamp):
				default:
					close(client.egress)
					delete(manager.clients, client)
				}
			}
		}
	}
}
