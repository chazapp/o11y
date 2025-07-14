package ws

import (
	"net/http"
	"time"

	"github.com/chazapp/o11y/apps/wall_api/metrics"
	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var WsReadBufferSize = 1024
var WsWriteBufferSize = 1024
var WsTimeoutSeconds = 10
var ChannelSize = 255

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  WsReadBufferSize,
	WriteBufferSize: WsWriteBufferSize,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan models.WallMessage
}

type Hub struct {
	clients   map[*Client]bool
	Broadcast chan models.WallMessage
	register  chan *Client
}

func NewHub() *Hub {
	metrics.WSClients.Set(0)

	return &Hub{
		Broadcast: make(chan models.WallMessage, ChannelSize),
		register:  make(chan *Client),
		clients:   make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			metrics.WSClients.Inc()

		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
				}
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(time.Duration(WsTimeoutSeconds) * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
		metrics.WSClients.Dec()
		delete(c.hub.clients, c)
		log.Error().Msg("Removed wwebsocket via defer")
	}()
	log.Info().Msg("Started write pump")
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(time.Duration(WsTimeoutSeconds) * time.Second)); err != nil {
				log.Err(err)
				return
			}
			if !ok {
				log.Error().Msg("Hub  closing write pump")
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				log.Err(err)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(time.Duration(WsTimeoutSeconds) * time.Second)); err != nil {
				log.Err(err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Err(err)
				return
			}
		}
	}
}

func (h *Hub) WsHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Err(err)
		return
	}

	client := &Client{hub: h, conn: conn, send: make(chan models.WallMessage, ChannelSize)}
	client.hub.register <- client

	go client.writePump()
}

func (h *Hub) GetCountConnectedClients() int {
	return len(h.clients)
}
