package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fourthfloor/internal/model"
	"net/http"
	"sync"
	"time"
)

type TwitchAPIClientInterface interface {
	FetchVideos(channelID string, limit int) ([]model.Video, error)
}

// TwitchAPIClient represents a Twitch API client with token management.
type TwitchAPIClient struct {
	ClientID     string
	ClientSecret string
	Token        string
	BaseURL      string

	expires          time.Time
	now              func() time.Time
	refreshTokenFunc func() (string, time.Time, error)
	httpClient       *http.Client

	mu sync.Mutex // protects token refresh
}

// NewTwitchAPIClient creates a TwitchAPIClient with default Twitch API URL.
func NewTwitchAPIClient(clientID, clientSecret string, options ...func(*TwitchAPIClient)) *TwitchAPIClient {
	c := &TwitchAPIClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      "https://api.twitch.tv/helix/videos",
		httpClient:   http.DefaultClient,
		now:          time.Now,
	}

	c.refreshTokenFunc = c.defaultRefreshTokenFunc

	for _, opt := range options {
		opt(c)
	}
	return c
}

// WithBaseURL allows overriding the base URL (useful for tests)
func WithBaseURL(url string) func(*TwitchAPIClient) {
	return func(c *TwitchAPIClient) { c.BaseURL = url }
}

// WithRefreshFunc allows injecting a custom token refresh function (useful for tests)
func WithRefreshFunc(fn func() (string, time.Time, error)) func(*TwitchAPIClient) {
	return func(c *TwitchAPIClient) { c.refreshTokenFunc = fn }
}

// WithExpires allows setting a custom expiry (useful for tests).
func WithExpires(t time.Time) func(*TwitchAPIClient) {
	return func(c *TwitchAPIClient) { c.expires = t }
}

// EnsureTokenValid refreshes the token if expired or near expiry.
func (c *TwitchAPIClient) EnsureTokenValid() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.now().After(c.expires) {
		newToken, newExpires, err := c.refreshTokenFunc()
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		c.Token = newToken
		c.expires = newExpires
	}
	return nil
}

func (c *TwitchAPIClient) defaultRefreshTokenFunc() (string, time.Time, error) {
	token, expiresIn, err := c.fetchToken()
	if err != nil {
		return "", time.Time{}, err
	}
	return token, time.Now().Add(time.Duration(expiresIn) * time.Second), nil
}

// fetchToken obtains a new OAuth token from Twitch.
func (c *TwitchAPIClient) fetchToken() (string, int, error) {
	type tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	data := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=client_credentials",
		c.ClientID, c.ClientSecret,
	)
	req, _ := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return "", 0, fmt.Errorf("failed to get token: status %d, body=%s", resp.StatusCode, buf.String())
	}

	var t tokenResp
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", 0, err
	}

	return t.AccessToken, t.ExpiresIn, nil
}

// FetchVideos fetches videos for a channel, ensuring a valid token first.
func (c *TwitchAPIClient) FetchVideos(channelID string, limit int) ([]model.Video, error) {
	if err := c.EnsureTokenValid(); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?user_id=%s&first=%d", c.BaseURL, channelID, limit), nil)
	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API returned %d", resp.StatusCode)
	}

	var result model.VideoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
