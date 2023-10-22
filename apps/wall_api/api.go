package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define the WallMessage model
type WallMessage struct {
	ID                uint `gorm:"primary_key"`
	Username          string
	Message           string
	CreationTimestamp time.Time
}

// Define your database connection
var db *gorm.DB

func API(dbUser, dbPassword, dbHost, dbName string, port int) (err error) {
	// Initialize the database
	initDB(dbUser, dbPassword, dbHost, dbName)
	defer db.Close()

	// Create a Gin router
	r := gin.New()
	r.Use(gin.Recovery())

	// Log using zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Middleware for structured logging
	r.Use(gin.LoggerWithWriter(log.Logger))
	// r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	return fmt.Sprintf("clientIP=%s method=%s statusCode=%d latency=%s path=%s\n",
	// 		param.ClientIP, param.Method, param.StatusCode, param.Latency, param.Path)
	// }))

	// Routes
	r.POST("/message", createMessage)
	r.GET("/message/:id", getMessage)
	r.GET("/messages", getMessages)
	r.GET("/health", healthCheck)
	r.GET("/pprof", gin.WrapF(pprof.Index))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
