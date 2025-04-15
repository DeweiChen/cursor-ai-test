package handlers

import (
	"log"
	"root/api/models"
	"root/api/repositories"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketConnection represents a connected client
type WebSocketConnection struct {
	conn     *websocket.Conn
	chatroom string
}

// WebSocketManager manages all WebSocket connections
type WebSocketManager struct {
	connections map[string][]*WebSocketConnection
	mu          sync.RWMutex
	messageRepo repositories.MessageRepository
}

func NewWebSocketManager(messageRepo repositories.MessageRepository) *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string][]*WebSocketConnection),
		messageRepo: messageRepo,
	}
}

// HandleWebSocket handles new WebSocket connections
func (m *WebSocketManager) HandleWebSocket(c *gin.Context) {
	chatroomID := c.Param("id")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create new connection
	wsConn := &WebSocketConnection{
		conn:     conn,
		chatroom: chatroomID,
	}

	// Add connection to manager
	m.addConnection(wsConn)
	defer m.removeConnection(wsConn)

	// Send existing messages
	messages, err := m.messageRepo.GetByChatroomID(chatroomID)
	if err == nil {
		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		}
	}

	// Handle incoming messages
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		// Set chatroom ID and save message
		message.ChatroomID = chatroomID
		if err := m.messageRepo.Create(&message); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Broadcast message to all connections in the chatroom
		m.broadcastMessage(chatroomID, message)
	}
}

func (m *WebSocketManager) addConnection(conn *WebSocketConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.chatroom] = append(m.connections[conn.chatroom], conn)
}

func (m *WebSocketManager) removeConnection(conn *WebSocketConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	connections := m.connections[conn.chatroom]
	for i, c := range connections {
		if c == conn {
			m.connections[conn.chatroom] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func (m *WebSocketManager) broadcastMessage(chatroomID string, message models.Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connections := m.connections[chatroomID]
	for _, conn := range connections {
		if err := conn.conn.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}
