package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chazapp/o11/apps/wall_api/models"
	"github.com/chazapp/o11/apps/wall_api/ws"
	"github.com/glebarez/sqlite"
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
)

func TestHealthcheck(t *testing.T) {
	router := newOpsEngine()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func createTestingDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.WallMessage{})
	return db
}

func TestCreateMessage(t *testing.T) {
	db := createTestingDB()
	wsHub := ws.NewHub()
	go wsHub.Run()

	router := newWallAPIEngine(db, wsHub, []string{"http://localhost"})
	w := httptest.NewRecorder()
	payload, err := json.Marshal(map[string]interface{}{
		"message":  "Hello world !",
		"username": "foo",
	})
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
