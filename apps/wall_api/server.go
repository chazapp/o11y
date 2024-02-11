package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chazapp/o11y/apps/wall_api/api"
	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/chazapp/o11y/apps/wall_api/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	httpmetrics "github.com/slok/go-http-metrics/middleware/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewWallAPIEngine(db *gorm.DB, wsHub *ws.Hub, allowedOrigins []string) *gin.Engine {
	r := gin.New()
	metricsMiddleware := middleware.New(middleware.Config{
		Recorder: prometheus.NewRecorder(prometheus.Config{}),
	})
	r.Use(otelgin.Middleware("wall-api"))
	//r.Use()
	r.Use(gin.Recovery())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	mr := api.NewMessageRouter(db, wsHub)

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
	r.POST("/message", httpmetrics.Handler("/message", metricsMiddleware), mr.CreateMessage)
	r.GET("/message/:id", httpmetrics.Handler("/message/:id", metricsMiddleware), mr.GetMessage)
	r.GET("/messages", httpmetrics.Handler("/messages", metricsMiddleware), mr.GetMessages)
	r.GET("/ws", httpmetrics.Handler("/ws", metricsMiddleware), wsHub.WsHandler)

	return r
}

func NewOpsEngine() *gin.Engine {
	r := gin.New()
	pprof.Register(r)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", healthCheck)
	return r
}

func API(dbUser, dbPassword, dbHost, dbName string, port int, opsPort int, allowedOrigins []string, otlpEndpoint string) (err error) {
	db := initDB(dbUser, dbPassword, dbHost, dbName, otlpEndpoint)
	wsHub := ws.NewHub()
	go wsHub.Run()
	r := NewWallAPIEngine(db, wsHub, allowedOrigins)
	opsRouter := NewOpsEngine()
	go opsRouter.Run(fmt.Sprintf("0.0.0.0:%d", opsPort))
	return r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

func initDB(dbUser, dbPassword, dbHost, dbName string, otlpEndpoint string) (db *gorm.DB) {
	dbURI := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	if otlpEndpoint != "" {
		if err := db.Use(tracing.NewPlugin(tracing.WithTracerProvider(otel.GetTracerProvider()))); err != nil {
			log.Fatal().Err(err)
		}
	}
	// Auto-migrate the schema
	db.AutoMigrate(&models.WallMessage{})

	// Limit connection pool
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err)
	}
	sqlDb.SetMaxOpenConns(100)
	return db
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
