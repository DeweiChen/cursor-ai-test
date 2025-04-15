package repositories

import (
	"root/api/models"
	"sync"
	"time"
)

// MessageRepository defines the interface for message operations
type MessageRepository interface {
	Create(message *models.Message) error
	GetByChatroomID(chatroomID string) ([]*models.Message, error)
}

// MemoryMessageRepository implements MessageRepository using in-memory storage
type MemoryMessageRepository struct {
	messages map[string]*models.Message
	mu       sync.RWMutex
}

// NewMemoryMessageRepository creates a new instance of MemoryMessageRepository
func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{
		messages: make(map[string]*models.Message),
	}
}

// Create adds a new message to the repository
func (r *MemoryMessageRepository) Create(message *models.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	r.messages[message.ID] = message
	return nil
}

// GetByChatroomID retrieves all messages for a specific chatroom
func (r *MemoryMessageRepository) GetByChatroomID(chatroomID string) ([]*models.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var messages []*models.Message
	for _, msg := range r.messages {
		if msg.ChatroomID == chatroomID {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}
