package main

import (
	"fourthfloor/internal/config"
	"fourthfloor/internal/handlers"
	"fourthfloor/internal/service"
	"fourthfloor/internal/twitch"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadEnv(".env")
	
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		log.Fatal("twitch credentials missing")
	}

	twitchClient := twitch.NewTwitchAPIClient(cfg.ClientID, cfg.ClientSecret)

	videoService := &service.VideoService{
		TwitchClient: twitchClient,
	}

	handler := &handlers.VideoHandler{Service: videoService}

	r := mux.NewRouter()
	r.HandleFunc("/streamers/{channel_id}/videos", handler.GetStreamerVideosHandler).Methods("GET")

	log.Printf("Server running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
