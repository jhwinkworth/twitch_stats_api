package twitch_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"fourthfloor/internal/model"
	"fourthfloor/internal/twitch"
)

// ---- Mocks ----

// mock token server response
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": "mock-token",
		"expires_in":   60,
		"token_type":   "bearer",
	})
}

// mock videos server response
func videosHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "missing auth", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.VideoResponse{
		Data: []model.Video{{Title: "Test Video"}},
	})
}

// ---- Tests ----

func TestFetchVideos(t *testing.T) {
	// spin up mock servers
	tokenSrv := httptest.NewServer(http.HandlerFunc(tokenHandler))
	defer tokenSrv.Close()

	videosSrv := httptest.NewServer(http.HandlerFunc(videosHandler))
	defer videosSrv.Close()

	// Custom refresh func hitting token server
	refresh := func() (string, time.Time, error) {
		req, _ := http.NewRequest("POST", tokenSrv.URL, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", time.Time{}, err
		}
		defer resp.Body.Close()

		var tr struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
			return "", time.Time{}, err
		}
		return tr.AccessToken, time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second), nil
	}

	// Client with mock video base URL and refresh func
	client := twitch.NewTwitchAPIClient("fake-client-id", "fake-secret",
		twitch.WithBaseURL(videosSrv.URL),
		twitch.WithRefreshFunc(refresh),
	)

	videos, err := client.FetchVideos("fake-channel", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(videos) != 1 || videos[0].Title != "Test Video" {
		t.Errorf("expected 1 video with Title=Test Video, got %+v", videos)
	}
}

func TestFetchVideosWithExpiredToken(t *testing.T) {
	// mock servers
	tokenSrv := httptest.NewServer(http.HandlerFunc(tokenHandler))
	defer tokenSrv.Close()

	videosSrv := httptest.NewServer(http.HandlerFunc(videosHandler))
	defer videosSrv.Close()

	// custom refresh func
	refresh := func() (string, time.Time, error) {
		req, _ := http.NewRequest("POST", tokenSrv.URL, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", time.Time{}, err
		}
		defer resp.Body.Close()

		var tr struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
			return "", time.Time{}, err
		}
		return tr.AccessToken, time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second), nil
	}

	// client starts with expired token
	client := twitch.NewTwitchAPIClient("fake-id", "fake-secret",
		twitch.WithBaseURL(videosSrv.URL),
		twitch.WithRefreshFunc(refresh),
		twitch.WithExpires(time.Now().Add(-1*time.Minute)), // force expired
	)
	client.Token = "stale-token"

	videos, err := client.FetchVideos("fake-channel", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(videos) != 1 || videos[0].Title != "Test Video" {
		t.Errorf("expected 1 refreshed video, got %+v", videos)
	}
}

func TestConcurrentTokenRefresh(t *testing.T) {
	// track refresh calls
	var refreshCalls int
	var mu sync.Mutex

	refresh := func() (string, time.Time, error) {
		mu.Lock()
		defer mu.Unlock()
		refreshCalls++
		return "new-token", time.Now().Add(time.Minute), nil
	}

	// mock video server
	videosSrv := httptest.NewServer(http.HandlerFunc(videosHandler))
	defer videosSrv.Close()

	client := twitch.NewTwitchAPIClient("id", "secret",
		twitch.WithBaseURL(videosSrv.URL),
		twitch.WithRefreshFunc(refresh),
		twitch.WithExpires(time.Now().Add(-1*time.Second)), // expired
	)

	// run multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.FetchVideos("chan", 1)
			if err != nil {
				t.Errorf("FetchVideos error: %v", err)
			}
		}()
	}
	wg.Wait()

	// assert refresh only called once
	mu.Lock()
	defer mu.Unlock()
	if refreshCalls != 1 {
		t.Errorf("expected refresh to be called once, got %d", refreshCalls)
	}
}
