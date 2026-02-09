//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/employee/create", pushEmployeeData)
	r.GET("/employee/search/all", fetchALLEmployeeData)
	r.GET("/employee/healthz", healthCheck)
	return r
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/employee/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}

func TestCreateEmployee(t *testing.T) {
	r := setupRouter()

	body := EmployeeInfo{
		ID:      "test123",
		Name:    "Test User",
		JobRole: "DevOps",
	}

	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/employee/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}

