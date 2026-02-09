package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_main(t *testing.T) {
	r := SetUpRouter()
	r.GET("/employee/healthz", healthCheck)

	req, _ := http.NewRequest("GET", "/employee/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_pushEmployeeData(t *testing.T) {
	r := SetUpRouter()
	r.POST("/employee/create", pushEmployeeData)

	employee := EmployeeInfo{ID: "1", Name: "Test"}
	body, _ := json.Marshal(employee)

	req, _ := http.NewRequest("POST", "/employee/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_fetchALLEmployeeData(t *testing.T) {
	r := SetUpRouter()
	r.GET("/employee/search/all", fetchALLEmployeeData)

	req, _ := http.NewRequest("GET", "/employee/search/all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
