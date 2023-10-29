package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type WallMessage struct {
	ID                uint      `gorm:"primary_key" json:"id"`
	Username          string    `json:"username"`
	Message           string    `json:"message"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

var db *gorm.DB

var wsHub = newHub()

func API(dbUser, dbPassword, dbHost, dbName string, port int, allowedOrigins []string) (err error) {
	initDB(dbUser, dbPassword, dbHost, dbName)
	defer db.Close()

	r := gin.New()
	r.Use(gin.Recovery())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// ToDo: figure out structured Logging
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("clientIP=%s method=%s statusCode=%d latency=%s path=%s\n",
			param.ClientIP, param.Method, param.StatusCode, param.Latency, param.Path)
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
	}))

	r.POST("/message", createMessage)
	r.GET("/message/:id", getMessage)
	r.GET("/messages", getMessages)

	go wsHub.run()
	r.GET("/ws", wsHandler)

	// Start /pprof, /metric and /health on port 8081
	optsRouter := gin.New()
	pprof.Register(optsRouter)
	optsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	optsRouter.GET("/health", healthCheck)

	go optsRouter.Run("0.0.0.0:8081")

	// Emit received message via WebSocket Broadcast
	// Start the server
	return r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

func initDB(dbUser, dbPassword, dbHost, dbName string) {
	dbURI := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	// Auto-migrate the schema
	db.AutoMigrate(&WallMessage{})
}

func createMessage(c *gin.Context) {
	var message WallMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.CreationTimestamp = time.Now()
	db.Create(&message)
	wsHub.broadcast <- message
	c.JSON(http.StatusCreated, message)
}

func getMessage(c *gin.Context) {
	id := c.Param("id")
	var message WallMessage
	if err := db.Where("id = ?", id).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}
	c.JSON(http.StatusOK, message)
}

func getMessages(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	var messages []WallMessage
	db.Limit(limit).Offset(offset).Find(&messages)

	c.JSON(http.StatusOK, messages)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
