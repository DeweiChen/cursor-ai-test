package handlers

import (
	"log"
	"root/api/models"
	"root/api/repositories"
	"sync"
	"time"

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
	connections  map[string][]*WebSocketConnection
	mu           sync.RWMutex
	messageRepo  repositories.MessageRepository
	lastActivity map[string]time.Time
}

func NewWebSocketManager(messageRepo repositories.MessageRepository) *WebSocketManager {
	manager := &WebSocketManager{
		connections:  make(map[string][]*WebSocketConnection),
		messageRepo:  messageRepo,
		lastActivity: make(map[string]time.Time),
	}
	go manager.startCleanupRoutine()
	return manager
}

func (m *WebSocketManager) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for chatroomID, lastActive := range m.lastActivity {
			connections, exists := m.connections[chatroomID]
			if (!exists || len(connections) == 0) && now.Sub(lastActive) > 10*time.Second {
				delete(m.connections, chatroomID)
				delete(m.lastActivity, chatroomID)
				log.Printf("Cleaned up inactive chatroom: %s", chatroomID)
			}
		}
		m.mu.Unlock()
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

	// Update last activity time
	m.updateLastActivity(chatroomID)

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

		// Update last activity time
		m.updateLastActivity(chatroomID)

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

func (m *WebSocketManager) updateLastActivity(chatroomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastActivity[chatroomID] = time.Now()
}

func (m *WebSocketManager) addConnection(conn *WebSocketConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.chatroom] = append(m.connections[conn.chatroom], conn)
	m.lastActivity[conn.chatroom] = time.Now()
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

	// Update last activity time when removing connection
	m.lastActivity[conn.chatroom] = time.Now()
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

func (m *WebSocketManager) HasActiveConnections(chatroomID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connections, exists := m.connections[chatroomID]
	if exists && len(connections) > 0 {
		return true
	}

	lastActive, exists := m.lastActivity[chatroomID]
	if !exists {
		return false
	}

	return time.Since(lastActive) <= 10*time.Second
}
