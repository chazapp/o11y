package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chazapp/o11y/apps/auth/jwt"
	"github.com/chazapp/o11y/apps/auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// port int, host string, db string, jwtPrivateKeyPath string, jwtPublicKeyPath string
type AuthServer struct {
	Port              int
	Host              string
	Domain            string
	DbConn            string
	JwtPrivateKeyPath string
	JwtPublicKeyPath  string
}

type JWKS struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

type AuthRouter struct {
	cfg        *AuthServer
	db         *gorm.DB
	jwkPublic  jose.JSONWebKey
	jwkPrivate jose.JSONWebKey
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthRouter(authServerCfg *AuthServer, testing bool) *AuthRouter {
	jwkPrivate, jwkPublic, err := jwt.LoadJWTKeys(authServerCfg.JwtPrivateKeyPath, authServerCfg.JwtPublicKeyPath)
	if err != nil {
		log.Fatalf("failed to load JWT keys: %v", err)
	}
	var db *gorm.DB
	if !testing {
		db, err = gorm.Open(postgres.Open(authServerCfg.DbConn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		db.AutoMigrate(&models.User{})
	}

	return &AuthRouter{
		cfg:        authServerCfg,
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

	token, expiry, err := jwt.GenerateJWT("24h", user.Email, r.jwkPrivate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("auth", token, int(expiry), "", r.cfg.Domain, true, true)
	c.JSON(200, gin.H{"token": token, "email": user.Email})
}

func (r *AuthRouter) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	var user models.User
	if err := r.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, expiry, err := jwt.GenerateJWT("24h", user.Email, r.jwkPrivate)
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("auth", token, int(expiry), "", r.cfg.Domain, true, true)
	c.JSON(200, gin.H{"token": token, "email": user.Email})
}

func (r *AuthRouter) Me(c *gin.Context) {
	u, exist := c.Get("user")
	fmt.Println("/me ????")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}
	user := u.(models.User)
	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
	})
}

func (r *AuthRouter) AuthMiddleware() gin.HandlerFunc {
	cleanup := func(c *gin.Context) {
		c.SetCookie("auth", "", 0, "", "", true, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
	}

	return func(c *gin.Context) {
		token := ""
		authCookie, _ := c.Cookie("auth")
		authHeader := c.GetHeader("Authorization")

		if authCookie != "" {
			token = authCookie
		} else if authHeader != "" {
			tok, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found {
				cleanup(c)
				return
			}
			token = tok
		} else {
			cleanup(c)
			return
		}

		claims, err := jwt.VerifyJWT(token, r.jwkPublic)
		if err != nil {
			cleanup(c)
			return
		}
		var user models.User
		if err := r.db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			cleanup(c)
			return

		}
		c.Set("user", user)
		c.Next()
	}
}

func (c *AuthServer) Run(ctx context.Context) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/metrics", "/.well-known/jwks.json"},
	}))
	authRouter := NewAuthRouter(c, false)

	engine.POST("/register", authRouter.Register)
	engine.POST("/login", authRouter.Login)
	engine.GET("/me", authRouter.AuthMiddleware(), authRouter.Me)

	engine.GET("/.well-known/jwks.json", func(c *gin.Context) {
		jwks := JWKS{
			Keys: []jose.JSONWebKey{authRouter.jwkPublic},
		}
		c.JSON(200, jwks)
	})

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return engine.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))
}
