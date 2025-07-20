package http

import (
	"context"
	"fmt"
	"log"

	"github.com/chazapp/o11y/apps/auth/jwt"
	"github.com/chazapp/o11y/apps/auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JWKS struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

type AuthRouter struct {
	db         *gorm.DB
	jwkPublic  jose.JSONWebKey
	jwkPrivate jose.JSONWebKey
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthRouter(dbConn, jwtPrivateKeyPath, jwtPublicKeyPath string, testing bool) *AuthRouter {
	jwkPrivate, jwkPublic, err := jwt.LoadJWTKeys(jwtPrivateKeyPath, jwtPublicKeyPath)
	if err != nil {
		log.Fatalf("failed to load JWT keys: %v", err)
	}
	var db *gorm.DB
	if !testing {
		db, err := gorm.Open(postgres.Open(dbConn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		db.AutoMigrate(&models.User{})
	}

	return &AuthRouter{
		db:         db,
		jwkPublic:  jwkPublic,
		jwkPrivate: jwkPrivate,
	}
}

func (r *AuthRouter) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := r.db.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.GenerateJWT("24h", user.Email, r.jwkPrivate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(200, gin.H{"token": token, "email": user.Email})
}

func (r *AuthRouter) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Find user
	var user models.User
	if err := r.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT with 24h expiration
	token, err := jwt.GenerateJWT("24h", user.Email, r.jwkPrivate)
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token, "email": user.Email})
}

func Run(ctx context.Context, port int, host string, db string, jwtPrivateKeyPath string, jwtPublicKeyPath string) error {
	engine := gin.Default()
	authRouter := NewAuthRouter(db, jwtPrivateKeyPath, jwtPublicKeyPath, false)

	engine.POST("/register", authRouter.Register)
	engine.POST("/login", authRouter.Login)

	// Serve the JWKS (JSON Web Key Set) containing the public key
	engine.GET("/.well-known/jwks.json", func(c *gin.Context) {
		jwks := JWKS{
			Keys: []jose.JSONWebKey{authRouter.jwkPublic},
		}
		c.JSON(200, jwks)
	})

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return engine.Run(fmt.Sprintf("%s:%d", host, port))
}
