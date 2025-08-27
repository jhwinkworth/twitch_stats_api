package service

import (
	"errors"
	"fourthfloor/internal/model"
	"fourthfloor/internal/twitch"
	"time"
)

// VideoServiceInterface defines the interface for fetching video stats.
type VideoServiceInterface interface {
	GetVideoStats(channelID string, limit int) (model.VideoStatsResponse, error)
}

// VideoService implements VideoServiceInterface
type VideoService struct {
	TwitchClient twitch.TwitchAPIClientInterface
}

// GetVideoStats fetches videos from TwitchClient and computes stats.
func (s *VideoService) GetVideoStats(channelID string, limit int) (model.VideoStatsResponse, error) {
	videos, err := s.TwitchClient.FetchVideos(channelID, limit)
	if err != nil {
		return model.VideoStatsResponse{}, err
	}

	if len(videos) == 0 {
		return model.VideoStatsResponse{}, errors.New("no videos found")
	}

	var totalViews int
	var totalDur float64
	var mostViewed model.Video

	for _, v := range videos {
		totalViews += v.ViewCount

		// parse duration
		if dur, err := time.ParseDuration(v.Duration); err == nil {
			totalDur += dur.Minutes()
		}

		// track most viewed video
		if v.ViewCount > mostViewed.ViewCount {
			mostViewed = v
		}
	}

	avgViews := float64(totalViews) / float64(len(videos))

	var avgViewsPerMinute float64
	if totalDur > 0 {
		avgViewsPerMinute = float64(totalViews) / totalDur
	}

	return model.VideoStatsResponse{
		TotalViews:           totalViews,
		AverageViews:         avgViews,
		TotalDurationMinutes: totalDur,
		AvgViewsPerMinute:    avgViewsPerMinute,
		MostViewedTitle:      mostViewed.Title,
		MostViewedViewCount:  mostViewed.ViewCount,
	}, nil
}
