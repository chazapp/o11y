package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	SQLMaxConnections = 100
)

func NewWallAPIEngine(db *gorm.DB, wsHub *ws.Hub, allowedOrigins []string) *gin.Engine {
	r := gin.New()
	r.Use(otelgin.Middleware("wall-api"))
	r.Use(gin.Recovery())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	metricsMiddleware := middleware.New(middleware.Config{
		Recorder: prometheus.NewRecorder(prometheus.Config{}),
	})

	mr := api.NewMessageRouter(db, wsHub)
	sr := api.NewStatusRouter(db, wsHub, Version)

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		currentSpan := trace.SpanFromContext(param.Request.Context())
		currentTrace := currentSpan.SpanContext().TraceID()
		logData := map[string]interface{}{
			"time":       param.TimeStamp.Format(time.RFC3339),
			"clientIP":   param.ClientIP,
			"method":     param.Method,
			"statusCode": param.StatusCode,
			"latency":    param.Latency.Milliseconds(),
			"path":       param.Path,
			"span_id":    currentSpan.SpanContext().SpanID().String(),
			"trace_id":   currentTrace.String(),
		}
		logJSON, err := json.Marshal(logData)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling log data: %v\n", err)

			return ""
		}

		return string(logJSON) + "\n"
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
	r.GET("/status", httpmetrics.Handler("/status", metricsMiddleware), sr.GetStatus)

	return r
}

func NewOpsEngine() *gin.Engine {
	r := gin.New()
	pprof.Register(r)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", healthCheck)

	return r
}

func API(dbUser, dbPassword, dbHost, dbName string,
	port int, opsPort int, allowedOrigins []string, otlpEndpoint string) error {
	wsHub := ws.NewHub()
	go wsHub.Run()

	db := initDB(dbUser, dbPassword, dbHost, dbName, otlpEndpoint)
	r := NewWallAPIEngine(db, wsHub, allowedOrigins)

	opsRouter := NewOpsEngine()

	//nolint:errcheck // Can't check for return error in Go routine
	go opsRouter.Run(fmt.Sprintf("0.0.0.0:%d", opsPort))

	return r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

func initDB(dbUser, dbPassword, dbHost, dbName string, otlpEndpoint string) *gorm.DB {
	dbURI := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	if otlpEndpoint != "" {
		if err = db.Use(tracing.NewPlugin(tracing.WithTracerProvider(otel.GetTracerProvider()))); err != nil {
			log.Fatal().Err(err)
		}
	}
	// Auto-migrate the schema
	err = db.AutoMigrate(&models.WallMessage{})

	if err != nil {
		log.Fatal().Err(err)
	}

	// Limit connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err)
	}

	sqlDB.SetMaxOpenConns(SQLMaxConnections)

	return db
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func status(c *gin.Context) {

}
