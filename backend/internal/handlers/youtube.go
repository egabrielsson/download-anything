package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"download-anything/internal/services"
)

type YouTubeRequest struct {
	URL  string `json:"url"`
	Type string `json:"type"` // "audio" or "video"
}

// DownloadYouTube handles YouTube download requests
func DownloadYouTube(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req YouTubeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if !isValidYouTubeURL(req.URL) {
		writeError(w, "Invalid YouTube URL", http.StatusBadRequest)
		return
	}

	// Validate download type
	var downloadType services.DownloadType
	switch req.Type {
	case "audio":
		downloadType = services.DownloadTypeAudio
	case "video":
		downloadType = services.DownloadTypeVideo
	default:
		writeError(w, "Invalid download type. Use 'audio' or 'video'", http.StatusBadRequest)
		return
	}

	// Download using service
	result, err := services.DownloadYouTube(req.URL, downloadType)
	if err != nil {
		writeError(w, fmt.Sprintf("Download failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer services.CleanupDownload(result.FilePath)

	// Read file and send response
	fileData, err := os.ReadFile(result.FilePath)
	if err != nil {
		writeError(w, "Failed to read downloaded file", http.StatusInternalServerError)
		return
	}

	// Extract ID for filename
	filename := extractYouTubeFilename(req.URL, req.Type)

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.Write(fileData)
}

func isValidYouTubeURL(url string) bool {
	// Match youtube.com/watch, youtu.be/, youtube.com/shorts/
	pattern := `^https?://(www\.)?(youtube\.com/(watch\?v=|shorts/)|youtu\.be/)[A-Za-z0-9_-]+`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func extractYouTubeFilename(url, downloadType string) string {
	// Extract the video ID from URL
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`youtube\.com/watch\?v=([A-Za-z0-9_-]+)`),
		regexp.MustCompile(`youtu\.be/([A-Za-z0-9_-]+)`),
		regexp.MustCompile(`youtube\.com/shorts/([A-Za-z0-9_-]+)`),
	}

	var ext string
	if downloadType == "audio" {
		ext = "mp3"
	} else {
		ext = "mp4"
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(url)
		if len(matches) > 1 {
			return fmt.Sprintf("youtube_%s_%s.%s", matches[1], downloadType, ext)
		}
	}
	return fmt.Sprintf("youtube_%s.%s", downloadType, ext)
}
