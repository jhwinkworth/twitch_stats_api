//go:build integration

package service_test

import (
	"fourthfloor/internal/config"
	"fourthfloor/internal/service"
	"fourthfloor/internal/twitch"
	"testing"
)

func TestVideoService_Integration(t *testing.T) {
	cfg := config.LoadEnv("../../.env")

	// create a real Twitch client
	client := twitch.NewTwitchAPIClient(cfg.ClientID, cfg.ClientSecret)

	videoService := &service.VideoService{TwitchClient: client}

	stats, err := videoService.GetVideoStats(cfg.ChannelID, 2)
	if err != nil {
		t.Fatalf("FetchVideos failed: %v", err)
	}

	// no concrete assertion against stats so just log
	t.Logf("Integration test returned stats: %+v", stats)
}
