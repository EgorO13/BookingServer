package Handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyLoginHandler(t *testing.T) {
	handler := DummyLoginHandler("testsecret")
	reqBody := `{"role":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader([]byte(reqBody)))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var response tokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Token)
}

func TestDummyLoginHandler_InvalidRole(t *testing.T) {
	handler := DummyLoginHandler("testsecret")
	reqBody := `{"role":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader([]byte(reqBody)))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
