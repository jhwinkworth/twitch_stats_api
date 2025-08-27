//go:build integration

package handlers_test

import (
	"encoding/json"
	"fourthfloor/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"fourthfloor/internal/handlers"
	"fourthfloor/internal/model"
	"fourthfloor/internal/service"
	"fourthfloor/internal/twitch"

	"github.com/gorilla/mux"
)

func TestVideoHandler_Integration(t *testing.T) {
	cfg := config.LoadEnv("../../.env")

	// create a real Twitch client
	client := twitch.NewTwitchAPIClient(cfg.ClientID, cfg.ClientSecret)

	videoService := &service.VideoService{
		TwitchClient: client,
	}

	handler := &handlers.VideoHandler{
		Service: videoService,
	}

	router := mux.NewRouter()
	router.HandleFunc("/streamers/{channel_id}/videos", handler.GetStreamerVideosHandler).Methods("GET")

	req := httptest.NewRequest("GET", "/streamers/"+cfg.ChannelID+"/videos?n=2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %d", rec.Code)
	}

	var stats model.VideoStatsResponse
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// No concrete assertion against stats so just log
	t.Logf("Integration test returned stats: %+v", stats)
}
