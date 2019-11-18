package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetIntruction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/instructions/5db658632b93a04483a56922", nil)
	if err != nil {
		t.Errorf("Get hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/api/v1/instructions failed with error code %d.", resp.Code)
	}
}

func TestGetAllInstruction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/instructions", nil)
	if err != nil {
		t.Errorf("Get hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/api/v1/instructions failed with error code %d.", resp.Code)
	}
}

func TestPostInstruction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	body := bytes.NewBuffer([]byte("{\"event_status\": \"83\", \"event_name\": \"Test\"}"))

	req, err := http.NewRequest("POST", "/api/v1/instructions", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Post hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/api/v1/instructions failed insert with error code %d.", resp.Code)
	}
}

func TestPutInstruction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	body := bytes.NewBuffer([]byte("{\"event_status\": \"83\", \"event_name\": \"Test Update\"}"))

	req, err := http.NewRequest("PUT", "/api/v1/instructions/5dd21fe19f715e10d4993d79", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Put hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/api/v1/instructions failed update with error code %d.", resp.Code)
	}
}
