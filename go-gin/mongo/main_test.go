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

func TestGetInstructions(t *testing.T) {
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

	body := bytes.NewBuffer([]byte("{\"EventStatus\": \"83\", \"EventName\": \"Test\"}"))

	req, err := http.NewRequest("POST", "/api/v1/instructions", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Post hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 201 {
		t.Errorf("/api/v1/instructions failed with error code %d.", resp.Code)
	}
}

func TestPutInstruction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	body := bytes.NewBuffer([]byte("{\"EventStatus\": \"83\", \"EventName\": \"Test Update\"}"))

	req, err := http.NewRequest("PUT", "/api/v1/instructions/5db689da9f715e3068813c21", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Put hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/api/v1/instructions failed with error code %d.", resp.Code)
	}
}
