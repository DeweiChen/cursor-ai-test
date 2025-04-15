package handlers

import (
	"net/http"
	"root/api/models"
	"root/api/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatroomHandler struct {
	repo repositories.ChatroomRepository
}

func NewChatroomHandler(repo repositories.ChatroomRepository) *ChatroomHandler {
	return &ChatroomHandler{repo: repo}
}

// CreateChatroom handles the creation of a new chatroom
func (h *ChatroomHandler) CreateChatroom(c *gin.Context) {
	var chatroom models.Chatroom
	if err := c.ShouldBindJSON(&chatroom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatroom.ID = uuid.New().String()
	chatroom.CreatedAt = time.Now()
	chatroom.UpdatedAt = time.Now()

	if err := h.repo.Create(&chatroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chatroom)
}

// GetChatroom retrieves a chatroom by ID
func (h *ChatroomHandler) GetChatroom(c *gin.Context) {
	id := c.Param("id")
	chatroom, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatroom)
}

// GetAllChatrooms retrieves all chatrooms
func (h *ChatroomHandler) GetAllChatrooms(c *gin.Context) {
	chatrooms, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatrooms)
}

// UpdateChatroom updates an existing chatroom
func (h *ChatroomHandler) UpdateChatroom(c *gin.Context) {
	id := c.Param("id")
	var chatroom models.Chatroom
	if err := c.ShouldBindJSON(&chatroom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatroom.ID = id
	if err := h.repo.Update(&chatroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatroom)
}

// DeleteChatroom deletes a chatroom by ID
func (h *ChatroomHandler) DeleteChatroom(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
