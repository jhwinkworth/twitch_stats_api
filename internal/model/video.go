package model

type Video struct {
	Title     string `json:"title"`
	ViewCount int    `json:"view_count"`
	Duration  string `json:"duration"`
}

type VideoResponse struct {
	Data []Video `json:"data"`
}

type VideoStatsResponse struct {
	TotalViews           int     `json:"total_views"`
	AverageViews         float64 `json:"average_views"`
	TotalDurationMinutes float64 `json:"total_duration_minutes"`
	AvgViewsPerMinute    float64 `json:"views_per_minute"`
	MostViewedTitle      string  `json:"most_viewed_title"`
	MostViewedViewCount  int     `json:"most_viewed_view_count"`
}
