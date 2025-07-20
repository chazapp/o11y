package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chazapp/o11y/apps/auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *AuthRouter) {
	// Use SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create temporary JWT key files
	privateKeyJSON := `{
		"use": "sig",
		"kty": "RSA",
		"kid": "yRlm6EGCiJBexDzolS6gwb6JsypHyh6088haxKnKv-8=",
		"alg": "RS256",
		"n": "yKUppaa1eD5lf9nZti5Jubakbn5QiUclELkldPbr2bEMRre0z6CJyEz6ScqbV4adrjvc4sLKBCB8pBZPBp-2Is1PWcnvALPRkJaIm8vSHy6hc91aMm-vFF9fnl3WUR7Y2Hph-mrBtFKssjKWqzUYOFom8VYd0wJThkgRDgVRsbVuc4mv1dS3Eh-dRky3BE3s3N3vm6L2OKLHTYFbe9XsQ2IBlqXEDv9DE16Js1mRHFkiWZuPg1Ay-b6zRyn_WTssDZDC2DIVoETrWLt61LXofIC_rn7kKrcsDhrUDVQ5yzaNpzD2yN8oXPpmtNnEOq3VnfDORvypTSDVv0Q4aO4TaQ",
		"e": "AQAB",
		"d": "4YFQQ8LvDd4pSQQvhtvKYhf2WvDGZAkpmGvMUDACnfS-yBrlQiJR8HsOMTH7qv3zLHH3SPWyOyB6LTnJqz0KRCaKICGN0dxIDxOmcvk7Wrz74lo9axDa8VzK52N2qk6YvxkbHn_rvZm45db4bGXw1CIwxrFi_XiedTW-3wou3r14ILDq3fb9PBotzGMI3WSXmMApeuRTz_-p5d3qzA8_wwM4fmhQ5hXG1969UmOh7Mc8qd0kjtFyCz0dNerfqiiyVSPcIO29N5GOVUKoAZLjjw1JbyWkowU4nomigZPXXOOhcBnOJ9niKDbUelsagWw6JIDDD1fKvj4Vgftpp40B",
		"p": "1mWDZEtXAalrCF7_0eO-QEw2E54LuN5doR4WG97LG96ylVBBK8YxM7BiN1VdCUHK2CO8z4BeqhzLV7r1gnkY1UrSaCKkC0FlY_I7taZcxXzdrO7nrszEuTkBeolmCGv6biUy9HxtryhJiqEVwuop2TEwKNhA_ay2YRO0LD-m8wE",
		"q": "75SHRxMKiEZug38vaUWvfKTkZxGzyT9iZEikpFoiP-kgQxqlxgW6EsILEbHXdBUCnYzTV1Z71XSfqIlPN2NGWNtxpA2STCPd1JiPjanlPAGZh54CJqHD-rff6fkSZaspO0Ici_2XNZZq_0YMvyWet82Xx1iyKPmgKdqd6ga8aGk",
		"dp": "MYF9RoJzE8ooEeXF0pRxEO3IKt16vXRzUEnfNw5J_iR9Picq9U2tfM8EztwiQIn1qdUOuydcNJGzjo14NWl7B31genVNReAS5nI_wWEp2NxNX6pGy0EzR8XXBpGFgvpT-G0UiAnXVfPKj31Exh5GDYXjJftRfoIMFvxyxSjphwE",
		"dq": "cJ-J2EeDM_yoBvjK-NnhXN7G4AzaT0iUoD_l5bzZTOHyYJkuRKB8kETXn0HS3qbhy95fmzb0j6t7QmcJ7iek8jB2g4A9vb0-kqoFEbtEH6lA2xfUOwTgdAPsJrkHhPOpNkol0Uksw-wp8Ealx1WP_yeOqg9v8QMn34T9pv3zo8E",
		"qi": "q3h5jVvsVg6s_v5BhuM3xtATmJ1h5WG8jH28i-x2aCZ1Rx3P4GqxBDeCWgQFDDVL2krSodIhiT0yrfpYcUPTkBM0nufcctdJfZ9wkEz4SFL2drDvB336im-k76tqGibRfpQnl9dwA2aDCkS3w8BJBy_8Aqj6q0L4YTOkBMakf4g"
	}`
	publicKeyJSON := `{
		"use": "sig",
		"kty": "RSA",
		"kid": "yRlm6EGCiJBexDzolS6gwb6JsypHyh6088haxKnKv-8=",
		"alg": "RS256",
		"n": "yKUppaa1eD5lf9nZti5Jubakbn5QiUclELkldPbr2bEMRre0z6CJyEz6ScqbV4adrjvc4sLKBCB8pBZPBp-2Is1PWcnvALPRkJaIm8vSHy6hc91aMm-vFF9fnl3WUR7Y2Hph-mrBtFKssjKWqzUYOFom8VYd0wJThkgRDgVRsbVuc4mv1dS3Eh-dRky3BE3s3N3vm6L2OKLHTYFbe9XsQ2IBlqXEDv9DE16Js1mRHFkiWZuPg1Ay-b6zRyn_WTssDZDC2DIVoETrWLt61LXofIC_rn7kKrcsDhrUDVQ5yzaNpzD2yN8oXPpmtNnEOq3VnfDORvypTSDVv0Q4aO4TaQ",
		"e": "AQAB"
	}`

	privKeyPath := "/tmp/test_private_key.json"
	pubKeyPath := "/tmp/test_public_key.json"

	err = os.WriteFile(privKeyPath, []byte(privateKeyJSON), 0600)
	if err != nil {
		t.Fatalf("Failed to write private key file: %v", err)
	}
	defer os.Remove(privKeyPath)

	err = os.WriteFile(pubKeyPath, []byte(publicKeyJSON), 0600)
	if err != nil {
		t.Fatalf("Failed to write public key file: %v", err)
	}
	defer os.Remove(pubKeyPath)

	authRouter := NewAuthRouter(":memory:", privKeyPath, pubKeyPath, true)
	authRouter.db = db // Use the in-memory database for testing
	db.AutoMigrate(&models.User{})
	engine := gin.Default()

	engine.POST("/register", authRouter.Register)
	engine.POST("/login", authRouter.Login)
	engine.GET("/.well-known/jwks.json", func(c *gin.Context) {
		jwks := JWKS{
			Keys: []jose.JSONWebKey{authRouter.jwkPublic},
		}
		c.JSON(200, jwks)
	})

	return engine, authRouter
}

func TestRegisterEndpoint(t *testing.T) {
	engine, _ := setupTestRouter(t)

	// Test successful registration
	w := httptest.NewRecorder()
	reqBody := AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["token"] == nil {
		t.Error("Expected token in response")
	}
	if response["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response["email"])
	}

	// Test duplicate registration
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate registration, got %d", w.Code)
	}
}

func TestLoginEndpoint(t *testing.T) {
	engine, _ := setupTestRouter(t)

	// First register a user
	reqBody := AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// Test successful login
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["token"] == nil {
		t.Error("Expected token in response")
	}
	if response["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response["email"])
	}

	// Test invalid credentials
	reqBody.Password = "wrongpassword"
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid credentials, got %d", w.Code)
	}
}

func TestJWKSEndpoint(t *testing.T) {
	engine, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/.well-known/jwks.json", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var jwks JWKS
	if err := json.Unmarshal(w.Body.Bytes(), &jwks); err != nil {
		t.Fatalf("Failed to parse JWKS response: %v", err)
	}

	if len(jwks.Keys) != 1 {
		t.Errorf("Expected 1 key in JWKS, got %d", len(jwks.Keys))
	}
}
