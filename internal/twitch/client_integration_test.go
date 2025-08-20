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

	videos, err := client.FetchVideos(cfg.ChannelID, 2)
	if err != nil {
		t.Fatalf("FetchVideos failed: %v", err)
	}

	if len(videos) == 0 {
		t.Errorf("Expected at least 1 video, got %d", len(videos))
	}
}
