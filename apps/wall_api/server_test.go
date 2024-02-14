package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/chazapp/o11y/apps/wall_api/ws"
	"github.com/glebarez/sqlite"
	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func createTestingDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(&models.WallMessage{}); err != nil {
		log.Panic().Err(err)
	}

	return db
}

func TestHealthcheck(t *testing.T) {
	router := NewOpsEngine()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestNewWallAPIEngine(t *testing.T) {
	db := createTestingDB()
	wsHub := ws.NewHub()
	engine := NewWallAPIEngine(db, wsHub, []string{
		"http://localhost:3000",
	})
	assert.Equal(t, len(engine.Routes()), 4)
}
