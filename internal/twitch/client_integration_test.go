//go:build integration

package twitch_test

import (
	"fourthfloor/internal/config"
	"testing"

	"fourthfloor/internal/twitch"
)

func TestFetchVideos_Integration(t *testing.T) {
	cfg := config.LoadEnv("../../.env")

	client := twitch.NewTwitchAPIClient(cfg.ClientID, cfg.ClientSecret)

	// Number of videos to return
	limit := 10

	videos, err := client.FetchVideos(cfg.ChannelID, limit)
	if err != nil {
		t.Fatalf("FetchVideos failed: %v", err)
	}

	if len(videos) != limit {
		t.Errorf("Wanted %d videos, got %d", limit, len(videos))
	}
}
