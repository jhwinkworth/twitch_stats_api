package model

// Video model for single video
type Video struct {
	Title     string `json:"title"`
	ViewCount int    `json:"view_count"`
	Duration  string `json:"duration"`
}

// VideoResponse response model for call to Twitch API
type VideoResponse struct {
	Data []Video `json:"data"`
}

// VideoStatsResponse response model for video stats
type VideoStatsResponse struct {
	TotalViews           int     `json:"total_views"`
	AverageViews         float64 `json:"average_views"`
	TotalDurationMinutes float64 `json:"total_duration_minutes"`
	AvgViewsPerMinute    float64 `json:"views_per_minute"`
	MostViewedTitle      string  `json:"most_viewed_title"`
	MostViewedViewCount  int     `json:"most_viewed_view_count"`
}
