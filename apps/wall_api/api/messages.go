package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chazapp/o11/apps/wall_api/metrics"
	"github.com/chazapp/o11/apps/wall_api/models"
	"github.com/chazapp/o11/apps/wall_api/ws"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageRouter struct {
	db    *gorm.DB
	wsHub *ws.Hub
}

func NewMessageRouter(db *gorm.DB, wsHub *ws.Hub) *MessageRouter {
	mr := MessageRouter{
		db:    db,
		wsHub: wsHub,
	}
	return &mr
}

func (r *MessageRouter) CreateMessage(c *gin.Context) {
	var message models.WallMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.CreationTimestamp = time.Now()
	r.db.WithContext(c.Request.Context()).Create(&message)
	r.wsHub.Broadcast <- message
	metrics.ProcessedMessages.Inc()
	c.JSON(http.StatusCreated, message)
}

func (r *MessageRouter) GetMessage(c *gin.Context) {
	id := c.Param("id")
	var message models.WallMessage
	if err := r.db.WithContext(c.Request.Context()).Where("id = ?", id).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}
	c.JSON(http.StatusOK, message)
}

func (r *MessageRouter) GetMessages(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid offset"})
		return
	}
	var messages []models.WallMessage
	r.db.WithContext(c.Request.Context()).Limit(limit).Offset(offset).Find(&messages)

	c.JSON(http.StatusOK, messages)
}
