package repositories

import (
	"errors"
	"root/api/models"
	"sync"
	"time"
)

// ChatroomRepository defines the interface for chatroom operations
type ChatroomRepository interface {
	Create(chatroom *models.Chatroom) error
	GetByID(id string) (*models.Chatroom, error)
	GetAll() ([]*models.Chatroom, error)
	Update(chatroom *models.Chatroom) error
	Delete(id string) error
}

// MemoryChatroomRepository implements ChatroomRepository using in-memory storage
type MemoryChatroomRepository struct {
	chatrooms map[string]*models.Chatroom
	mu        sync.RWMutex
}

// NewMemoryChatroomRepository creates a new instance of MemoryChatroomRepository
func NewMemoryChatroomRepository() *MemoryChatroomRepository {
	return &MemoryChatroomRepository{
		chatrooms: make(map[string]*models.Chatroom),
	}
}

// Create adds a new chatroom to the repository
func (r *MemoryChatroomRepository) Create(chatroom *models.Chatroom) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.chatrooms[chatroom.ID]; exists {
		return errors.New("chatroom already exists")
	}

	chatroom.CreatedAt = time.Now()
	chatroom.UpdatedAt = time.Now()
	r.chatrooms[chatroom.ID] = chatroom
	return nil
}

// GetByID retrieves a chatroom by its ID
func (r *MemoryChatroomRepository) GetByID(id string) (*models.Chatroom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chatroom, exists := r.chatrooms[id]
	if !exists {
		return nil, errors.New("chatroom not found")
	}
	return chatroom, nil
}

// GetAll retrieves all chatrooms
func (r *MemoryChatroomRepository) GetAll() ([]*models.Chatroom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chatrooms := make([]*models.Chatroom, 0, len(r.chatrooms))
	for _, chatroom := range r.chatrooms {
		chatrooms = append(chatrooms, chatroom)
	}
	return chatrooms, nil
}

// Update modifies an existing chatroom
func (r *MemoryChatroomRepository) Update(chatroom *models.Chatroom) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.chatrooms[chatroom.ID]; !exists {
		return errors.New("chatroom not found")
	}

	chatroom.UpdatedAt = time.Now()
	r.chatrooms[chatroom.ID] = chatroom
	return nil
}

// Delete removes a chatroom by its ID
func (r *MemoryChatroomRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.chatrooms[id]; !exists {
		return errors.New("chatroom not found")
	}

	delete(r.chatrooms, id)
	return nil
}
