package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"download-anything/internal/services"
)

type InstagramRequest struct {
	URL  string `json:"url"`
	Type string `json:"type"` // "audio" or "video"
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Health returns server health status
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DownloadAudio is a legacy endpoint for backward compatibility
func DownloadAudio(w http.ResponseWriter, r *http.Request) {
	// Redirect to new handler with audio type
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service directly
	if !isValidInstagramURL(req.URL) {
		writeError(w, "Invalid Instagram URL", http.StatusBadRequest)
		return
	}

	result, err := services.DownloadInstagram(req.URL, services.DownloadTypeAudio)
	if err != nil {
		writeError(w, fmt.Sprintf("Download failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer services.CleanupDownload(result.FilePath)

	fileData, err := os.ReadFile(result.FilePath)
	if err != nil {
		writeError(w, "Failed to read downloaded file", http.StatusInternalServerError)
		return
	}

	filename := extractInstagramFilename(req.URL, "audio")
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.Write(fileData)
}

// DownloadInstagram handles Instagram download requests
func DownloadInstagram(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InstagramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if !isValidInstagramURL(req.URL) {
		writeError(w, "Invalid Instagram URL", http.StatusBadRequest)
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
	result, err := services.DownloadInstagram(req.URL, downloadType)
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
	filename := extractInstagramFilename(req.URL, req.Type)

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.Write(fileData)
}

func isValidInstagramURL(url string) bool {
	pattern := `^https?://(www\.)?instagram\.com/(reel|reels|p)/[A-Za-z0-9_-]+/?`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func extractInstagramFilename(url, downloadType string) string {
	pattern := regexp.MustCompile(`instagram\.com/(?:reel|reels|p)/([A-Za-z0-9_-]+)`)
	matches := pattern.FindStringSubmatch(url)

	var ext string
	if downloadType == "audio" {
		ext = "mp3"
	} else {
		ext = "mp4"
	}

	if len(matches) > 1 {
		return fmt.Sprintf("instagram_%s_%s.%s", matches[1], downloadType, ext)
	}
	return fmt.Sprintf("instagram_%s.%s", downloadType, ext)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
