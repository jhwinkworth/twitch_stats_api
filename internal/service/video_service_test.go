package service_test

import (
	"errors"
	"fourthfloor/internal/model"
	"fourthfloor/internal/service"
	"testing"
	"time"
)

// ---- Mocks ----

// mockTwitchClient implements TwitchAPIClientInterface for testing.
type mockTwitchClient struct {
	videos []model.Video
	err    error
}

// FetchVideos mock return from FetchVideos function (client.go)
func (m *mockTwitchClient) FetchVideos(channelID string, limit int) ([]model.Video, error) {
	return m.videos, m.err
}

// ---- Tests ----

func TestVideoService_GetVideoStats(t *testing.T) {
	tests := []struct {
		name          string
		videos        []model.Video
		clientErr     error
		expectedErr   bool
		expectedTotal int
		expectedMost  string
	}{
		{
			name: "success case",
			videos: []model.Video{
				{Title: "Vid1", ViewCount: 100, Duration: "10m0s"},
				{Title: "Vid2", ViewCount: 200, Duration: "20m0s"},
			},
			expectedErr:   false,
			expectedTotal: 300,
			expectedMost:  "Vid2",
		},
		{
			name:          "no videos returned",
			videos:        []model.Video{},
			expectedErr:   true,
			expectedTotal: 0,
		},
		{
			name:        "client error",
			clientErr:   errors.New("fetch failed"),
			expectedErr: true,
		},
		{
			name: "invalid duration format",
			videos: []model.Video{
				{Title: "BadDuration", ViewCount: 50, Duration: "invalid"},
			},
			expectedErr:   false,
			expectedTotal: 50,
			expectedMost:  "BadDuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockTwitchClient{
				videos: tt.videos,
				err:    tt.clientErr,
			}

			svc := &service.VideoService{
				TwitchClient: mockClient,
			}

			stats, err := svc.GetVideoStats("channel1", 10)

			if tt.expectedErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectedErr {
				if stats.TotalViews != tt.expectedTotal {
					t.Errorf("expected total views %d, got %d", tt.expectedTotal, stats.TotalViews)
				}
				if stats.MostViewedTitle != tt.expectedMost {
					t.Errorf("expected most viewed %q, got %q", tt.expectedMost, stats.MostViewedTitle)
				}

				// only check total duration if at least one valid duration exists
				hasValidDuration := false
				for _, v := range tt.videos {
					if _, err := time.ParseDuration(v.Duration); err == nil {
						hasValidDuration = true
						break
					}
				}

				if hasValidDuration && stats.TotalDurationMinutes <= 0 {
					t.Errorf("wanted total duration > 0, got %f", stats.TotalDurationMinutes)
				}

				if stats.TotalDurationMinutes > 0 && stats.AvgViewsPerMinute <= 0 {
					t.Errorf("wanted AvgViewsPerMinute > 0, got %f", stats.AvgViewsPerMinute)
				}

				if stats.AverageViews <= 0 {
					t.Errorf("wanted AverageViews > 0, got %f", stats.AverageViews)
				}
			}
		})
	}
}
