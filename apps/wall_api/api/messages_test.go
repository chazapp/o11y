package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chazapp/o11y/apps/wall_api/models"
	"github.com/chazapp/o11y/apps/wall_api/ws"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	assert "github.com/go-playground/assert/v2"
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

// toJSONReader converts a JSON payload into an io.Reader interface.
func toJSONReader(payload map[string]interface{}) io.Reader {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(jsonBytes)
}

func setupTestEnvironment() (*gorm.DB, *ws.Hub, *MessageRouter, *gin.Engine) {
	db := createTestingDB()
	wsHub := ws.NewHub()
	go wsHub.Run()
	mr := NewMessageRouter(db, wsHub)
	r := gin.New()
	r.POST("/message", mr.CreateMessage)
	r.GET("/message/:id", mr.GetMessage)
	r.GET("/messages", mr.GetMessages)
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

func TestGetMessagesWithLimitAndOffset(t *testing.T) {
	_, _, _, r := setupTestEnvironment()

	payloads := []map[string]interface{}{
		{"message": "Message 1", "username": "user1"},
		{"message": "Message 2", "username": "user2"},
		{"message": "Message 3", "username": "user3"},
	}

	for _, payload := range payloads {
		p, _ := json.Marshal(payload)
		sendMessage(t, r, p)
	}
	limit := 2
	offset := 1

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/messages?limit=%d&offset=%d", limit, offset), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response []map[string]interface{}
	err := json.NewDecoder(w.Result().Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, limit, len(response))

	expectedMessages := payloads[1:]
	for i, expected := range expectedMessages {
		assert.Equal(t, expected["message"], response[i]["message"])
		assert.Equal(t, expected["username"], response[i]["username"])
	}
}

func TestBadRequests(t *testing.T) {
	_, _, _, r := setupTestEnvironment()

	// Send a POST request with missing required fields
	payload := map[string]interface{}{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/message", toJSONReader(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Send a GET request with non-numeric id parameter
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/message/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Send a GET request to a non-existent endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Send a GET request with non-numeric limit
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/messages?limit=aaaa&offset=0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Send a GET request with non-numeric offset
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/messages?limit=10&offset=aaaaa", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
