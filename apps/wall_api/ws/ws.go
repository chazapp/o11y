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
		Broadcast: make(chan models.WallMessage),
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
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(time.Duration(WsTimeoutSeconds) * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		metrics.WSClients.Dec()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(time.Duration(WsTimeoutSeconds) * time.Second))

			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(time.Duration(WsTimeoutSeconds) * time.Second))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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

	client := &Client{hub: h, conn: conn, send: make(chan models.WallMessage)}
	client.hub.register <- client

	go client.writePump()
}
