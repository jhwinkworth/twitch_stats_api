package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fourthfloor/internal/service"

	"github.com/gorilla/mux"
)

type VideoHandler struct {
	Service service.VideoServiceInterface
}

func (h *VideoHandler) GetStreamerVideosHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params["channel_id"]
	nStr := r.URL.Query().Get("n")

	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		http.Error(w, "Invalid query parameter 'n'", http.StatusBadRequest)
		return
	}

	stats, err := h.Service.GetVideoStats(channelID, n)
	if err != nil {
		// map service errors to HTTP codes
		if err.Error() == "no videos found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get video stats: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
