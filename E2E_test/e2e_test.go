package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

type tokenResp struct {
	Token string `json:"token"`
}
type roomResp struct {
	Room struct {
		ID string `json:"id"`
	} `json:"room"`
}
type slotListResp struct {
	Slots []struct {
		ID string `json:"id"`
	} `json:"slots"`
}
type bookingResp struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"booking"`
}

func getToken(t *testing.T, role string) string {
	body := []byte(`{"role":"` + role + `"}`)
	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	var tr tokenResp
	err = json.NewDecoder(resp.Body).Decode(&tr)
	require.NoError(t, err)
	return tr.Token
}

func doRequest(t *testing.T, method, url, token string, reqBody interface{}) *http.Response {
	var body []byte
	if reqBody != nil {
		var err error
		body, err = json.Marshal(reqBody)
		require.NoError(t, err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

func TestE2E_CreateBookingAndCancel(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	roomBody := map[string]interface{}{
		"name":        "E2E Room",
		"description": "Test",
		"capacity":    5,
	}
	resp := doRequest(t, "POST", baseURL+"/rooms/create", adminToken, roomBody)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var rr roomResp
	err := json.NewDecoder(resp.Body).Decode(&rr)
	require.NoError(t, err)
	roomID := rr.Room.ID
	require.NotEmpty(t, roomID)

	scheduleBody := map[string]interface{}{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "18:00",
	}
	resp = doRequest(t, "POST", baseURL+"/rooms/"+roomID+"/schedule/create", adminToken, scheduleBody)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	tomorrow := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	resp = doRequest(t, "GET", baseURL+"/rooms/"+roomID+"/slots/list?date="+tomorrow, userToken, nil)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var sl slotListResp
	err = json.NewDecoder(resp.Body).Decode(&sl)
	require.NoError(t, err)
	require.NotEmpty(t, sl.Slots, "No slots generated for tomorrow")
	slotID := sl.Slots[0].ID

	bookingBody := map[string]interface{}{"slotId": slotID}
	resp = doRequest(t, "POST", baseURL+"/bookings/create", userToken, bookingBody)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var br bookingResp
	err = json.NewDecoder(resp.Body).Decode(&br)
	require.NoError(t, err)
	bookingID := br.Booking.ID
	require.NotEmpty(t, bookingID)
	assert.Equal(t, "active", br.Booking.Status)

	resp = doRequest(t, "POST", baseURL+"/bookings/"+bookingID+"/cancel", userToken, nil)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&br)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", br.Booking.Status)

	resp = doRequest(t, "POST", baseURL+"/bookings/"+bookingID+"/cancel", userToken, nil)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
