package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"partyplanner/bus"
	"partyplanner/service"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type MessageManager struct {
	clients          map[*ChatClient]bool
	messageBroadcast chan *Message
	sync.RWMutex
}

type EventManager struct {
	clients  map[*EventClient]bool
	eventBus *bus.EventBus
	sync.RWMutex
}

func NewMessageManager() *MessageManager {
	return &MessageManager{
		clients:          make(map[*ChatClient]bool),
		messageBroadcast: make(chan *Message),
	}
}

func NewEventManager(bus *bus.EventBus) *EventManager {
	return &EventManager{
		clients:  make(map[*EventClient]bool),
		eventBus: bus,
	}
}

func (manager *MessageManager) addClient(client *ChatClient) {

	fmt.Printf("Connection from client : %s\n", client.connection.RemoteAddr())
	// Lock the client map
	manager.Lock()
	// Unocks when it is done
	defer manager.Unlock()

	manager.clients[client] = true
}

func (manager *MessageManager) removeWebsocket(client *ChatClient) {
	fmt.Println("Disconnecting client: ", client.connection.RemoteAddr())
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.clients[client]; ok {
		client.connection.Close()
		delete(manager.clients, client)
	}
}

func (manager *EventManager) removeWebsocket(client *EventClient) {
	fmt.Println("Disconnecting client: ", client.connection.RemoteAddr())
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.clients[client]; ok {
		client.connection.Close()
		delete(manager.clients, client)
	}
}

func (manager *MessageManager) ServeChat(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error while connecting: ", err)
	}
	client := NewChatClient(connection, manager)
	manager.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (manager *EventManager) ServeEvents(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error while connecting: ", err)
	}
	client := NewEventClient(connection, manager)
	manager.Lock()
	defer manager.Unlock()
	fmt.Printf("Connection from event client : %s\n", client.connection.RemoteAddr())

	manager.clients[client] = true

	go client.writeMessages()
}

func (manager *MessageManager) Run() {
	for {
		select {
		case message := <-manager.messageBroadcast:
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

func (manager *EventManager) ListenForEvents() {
	for {
		select {
		case event := <-manager.eventBus.Events:
			for client := range manager.clients {
				select {
				case client.eventQueue <- event:
				default:
					close(client.eventQueue)
					delete(manager.clients, client)
				}
			}
		}
	}
}