package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"download-anything/internal/services"
)

type TikTokRequest struct {
	URL  string `json:"url"`
	Type string `json:"type"` // "audio" or "video"
}

// DownloadTikTok handles TikTok download requests
func DownloadTikTok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TikTokRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if !isValidTikTokURL(req.URL) {
		writeError(w, "Invalid TikTok URL", http.StatusBadRequest)
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
	result, err := services.DownloadTikTok(req.URL, downloadType)
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
	filename := extractTikTokFilename(req.URL, req.Type)

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.Write(fileData)
}

func isValidTikTokURL(url string) bool {
	// Match tiktok.com URLs
	pattern := `^https?://(www\.|vm\.)?tiktok\.com/`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func extractTikTokFilename(url, downloadType string) string {
	// Extract the video ID from URL
	pattern := regexp.MustCompile(`tiktok\.com/@[^/]+/video/(\d+)`)
	matches := pattern.FindStringSubmatch(url)

	var ext string
	if downloadType == "audio" {
		ext = "mp3"
	} else {
		ext = "mp4"
	}

	if len(matches) > 1 {
		return fmt.Sprintf("tiktok_%s_%s.%s", matches[1], downloadType, ext)
	}
	return fmt.Sprintf("tiktok_%s.%s", downloadType, ext)
}
