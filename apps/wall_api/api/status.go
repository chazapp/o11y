package api

import (
	"net/http"
	"runtime/debug"

	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/chazapp/o11y/apps/wall_api/ws"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type StatusRouter struct {
	db      *gorm.DB
	wsHub   *ws.Hub
	Version string
}

type Status struct {
	GoVersion     string `json:"goVersion"`
	Version       string `json:"version"`
	ConnectedWS   int    `json:"connectedWS"`
	MessagesCount int    `json:"messagesCount"`
}

func NewStatusRouter(db *gorm.DB, wsHub *ws.Hub, version string) *StatusRouter {
	sr := StatusRouter{
		db:      db,
		wsHub:   wsHub,
		Version: version,
	}

	return &sr
}

func (sr *StatusRouter) GetStatus(c *gin.Context) {
	var messageCount int64
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Error().Msg("Error reading build info")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading build info",
		})
	}

	sr.db.Model(&models.WallMessage{}).Count(&messageCount)
	s := Status{
		GoVersion:     bi.GoVersion,
		Version:       sr.Version,
		ConnectedWS:   sr.wsHub.GetCountConnectedClients(),
		MessagesCount: int(messageCount),
	}
	c.JSON(http.StatusOK, s)
}
