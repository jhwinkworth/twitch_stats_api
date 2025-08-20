package handlers_test

import (
	"encoding/json"
	"errors"
	"fourthfloor/internal/handlers"
	"fourthfloor/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// ---- Mocks ----

// mockVideoService implements VideoServiceInterface for testing.
type mockVideoService struct {
	Response model.VideoStatsResponse
	Err      error
}

func (m *mockVideoService) GetVideoStats(channelID string, limit int) (model.VideoStatsResponse, error) {
	return m.Response, m.Err
}

// ---- Tests ----

func TestGetStreamerVideosHandler(t *testing.T) {
	tests := []struct {
		name           string
		channelID      string
		queryN         string
		serviceResp    model.VideoStatsResponse
		serviceErr     error
		expectedCode   int
		expectedInBody string
	}{
		{
			name:      "successful case",
			channelID: "123",
			queryN:    "5",
			serviceResp: model.VideoStatsResponse{
				TotalViews: 100,
			},
			serviceErr:     nil,
			expectedCode:   http.StatusOK,
			expectedInBody: `"total_views":100`,
		},
		{
			name:           "invalid n query",
			channelID:      "123",
			queryN:         "abc",
			expectedCode:   http.StatusBadRequest,
			expectedInBody: "Invalid query parameter 'n'",
		},
		{
			name:           "no videos found",
			channelID:      "123",
			queryN:         "5",
			serviceErr:     errors.New("no videos found"),
			expectedCode:   http.StatusNotFound,
			expectedInBody: "no videos found",
		},
		{
			name:           "internal server error",
			channelID:      "123",
			queryN:         "5",
			serviceErr:     errors.New("some failure"),
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: "Failed to get video stats: some failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockVideoService{
				Response: tt.serviceResp,
				Err:      tt.serviceErr,
			}

			handler := &handlers.VideoHandler{Service: mockSvc}

			req := httptest.NewRequest("GET", "/streamers/"+tt.channelID+"?n="+tt.queryN, nil)
			req = mux.SetURLVars(req, map[string]string{"channel_id": tt.channelID})

			rec := httptest.NewRecorder()
			handler.GetStreamerVideosHandler(rec, req)

			if rec.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rec.Code)
			}

			body := rec.Body.String()
			if !strings.Contains(body, tt.expectedInBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedInBody, body)
			}

			// For successful responses, also check JSON can be decoded
			if rec.Code == http.StatusOK {
				var resp model.VideoStatsResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("failed to decode JSON response: %v", err)
				}
			}
		})
	}
}
