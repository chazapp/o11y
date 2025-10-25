package http

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/chazapp/apps/url-shortner/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type URLShortnerRouter struct {
	db *gorm.DB
}

type CreateShortLinkRequest struct {
	URL string `json:"url"`
}

func NewURLShortnerRouter(dbConn string) *URLShortnerRouter {
	db, err := gorm.Open(postgres.Open(dbConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.ShortLink{})
	return &URLShortnerRouter{
		db: db,
	}
}

func Run(ctx context.Context, port int, host string, db string) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/metrics"},
	}))
	urlShortnerRouter := NewURLShortnerRouter(db)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/", urlShortnerRouter.CreateShortLink)
	engine.GET("/s/:short_url", urlShortnerRouter.GetShortLink)

	return engine.Run(fmt.Sprintf("%s:%d", host, port))
}

func generateShortURL() string {
	// Placeholder function to generate a short URL
	const charset = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 15)
	for i := range result {
		result[i] = charset[rand.IntN(len(charset))]
	}
	return string(result)
}

func getAuthorFromJWTClaims(c *gin.Context) string {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return "Anon"
	}
	return "NotAnon"
}

func (r *URLShortnerRouter) CreateShortLink(c *gin.Context) {
	var req CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	shortLink := models.ShortLink{
		URL:      req.URL,
		ShortURL: generateShortURL(),
		Author:   getAuthorFromJWTClaims(c),
	}
	if err := r.db.Create(&shortLink).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create short link"})
		return
	}
	c.JSON(201, gin.H{"short_url": shortLink.ShortURL})
}

func (r *URLShortnerRouter) GetShortLink(c *gin.Context) {
	shortURL := c.Param("short_url")
	var shortLink models.ShortLink
	if err := r.db.Where("short_url = ?", shortURL).First(&shortLink).Error; err != nil {
		c.JSON(404, gin.H{"error": "Short link not found"})
		return
	}
	c.Header("Location", shortLink.URL)
	c.JSON(301, gin.H{"url": shortLink.URL})
}
