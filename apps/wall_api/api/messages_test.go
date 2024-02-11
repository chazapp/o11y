package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/chazapp/o11y/apps/wall_api/ws"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-playground/assert"
	"gorm.io/gorm"
)

func createTestingDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.WallMessage{})
	return db
}

// func TestCreateMessage(t *testing.T) {
// 	db := createTestingDB()
// 	wsHub := ws.NewHub()
// 	go wsHub.Run()
// 	mr := NewMessageRouter(db, wsHub)
// 	r := gin.New()
// 	r.POST("/message", mr.CreateMessage)
// 	w := httptest.NewRecorder()
// 	payload, err := json.Marshal(map[string]interface{}{
// 		"message":  "Hello world !",
// 		"username": "foo",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(payload))
// 	req.Header.Set("Content-Type", "application/json")
// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, 201, w.Code)
// }

// func TestGetMessage(t *testing.T) {
// 	db := createTestingDB()
// 	wsHub := ws.NewHub()
// 	go wsHub.Run()
// 	mr := NewMessageRouter(db, wsHub)
// 	r := gin.New()
// 	r.POST("/message", mr.CreateMessage)
// 	w := httptest.NewRecorder()
// 	payload, err := json.Marshal(map[string]interface{}{
// 		"message":  "Hello world !",
// 		"username": "foo",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(payload))
// 	req.Header.Set("Content-Type", "application/json")
// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, 201, w.Code)
// 	var m map[string]interface{}
// 	err = json.NewDecoder(w.Result().Body).Decode(&m)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req2, err := http.NewRequest("GET", fmt.Sprintf("/message/%d", int(m["id"].(float64))), bytes.NewBuffer(payload))
// 	if err != nil {
// 		panic(err)
// 	}
// 	r.ServeHTTP(w, req2)
// 	assert.Equal(t, 200, w.Code)
// 	err = json.NewDecoder(w.Result().Body).Decode(&m)
// 	if err != nil {
// 		panic(err)
// 	}
// 	assert.Equal(t, "Hello world !", m["message"])

// }

func setupTestEnvironment() (*gorm.DB, *ws.Hub, *MessageRouter, *gin.Engine) {
	db := createTestingDB()
	wsHub := ws.NewHub()
	go wsHub.Run()
	mr := NewMessageRouter(db, wsHub)
	r := gin.New()
	r.POST("/message", mr.CreateMessage)
	r.GET("/message/:id", mr.GetMessage)
	return db, wsHub, mr, r
}

func sendMessage(t *testing.T, r *gin.Engine, payload []byte) map[string]interface{} {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var m map[string]interface{}
	err := json.NewDecoder(w.Result().Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	return m
}

func TestCreateMessage(t *testing.T) {
	_, _, _, r := setupTestEnvironment()

	payload, err := json.Marshal(map[string]interface{}{
		"message":  "Hello world !",
		"username": "foo",
	})
	if err != nil {
		panic(err)
	}

	sendMessage(t, r, payload)
}

func TestGetMessage(t *testing.T) {
	_, _, _, r := setupTestEnvironment()

	payload, err := json.Marshal(map[string]interface{}{
		"message":  "Hello world !",
		"username": "foo",
	})
	if err != nil {
		panic(err)
	}

	m := sendMessage(t, r, payload)

	req, err := http.NewRequest("GET", fmt.Sprintf("/message/%d", int(m["id"].(float64))), nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err = json.NewDecoder(w.Result().Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello world !", response["message"])
}
