package handlers

import (
	"net/http"
	"root/api/models"
	"root/api/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageRepo  repositories.MessageRepository
	chatroomRepo repositories.ChatroomRepository
}

func NewMessageHandler(messageRepo repositories.MessageRepository, chatroomRepo repositories.ChatroomRepository) *MessageHandler {
	return &MessageHandler{
		messageRepo:  messageRepo,
		chatroomRepo: chatroomRepo,
	}
}

// CreateMessage handles the creation of a new message in a chatroom
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	chatroomID := c.Param("id")

	// Verify chatroom exists
	if _, err := h.chatroomRepo.GetByID(chatroomID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chatroom not found"})
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.ID = uuid.New().String()
	message.ChatroomID = chatroomID
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	if err := h.messageRepo.Create(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessagesByChatroomID retrieves all messages for a specific chatroom
func (h *MessageHandler) GetMessagesByChatroomID(c *gin.Context) {
	chatroomID := c.Param("id")

	// Verify chatroom exists
	if _, err := h.chatroomRepo.GetByID(chatroomID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chatroom not found"})
		return
	}

	messages, err := h.messageRepo.GetByChatroomID(chatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
