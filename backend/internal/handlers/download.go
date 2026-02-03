package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type DownloadRequest struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DownloadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate Instagram URL
	if !isValidInstagramURL(req.URL) {
		writeError(w, "Invalid Instagram URL. Please provide a valid Instagram reel URL.", http.StatusBadRequest)
		return
	}

	// Create temp directory for download
	tempDir, err := os.MkdirTemp("", "insta-audio-*")
	if err != nil {
		writeError(w, "Failed to create temp directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	outputTemplate := filepath.Join(tempDir, "audio.%(ext)s")

	// Use yt-dlp to extract audio
	cmd := exec.Command("yt-dlp",
		"-x",                    // Extract audio
		"--audio-format", "mp3", // Convert to mp3
		"--audio-quality", "0", // Best quality
		"-o", outputTemplate, // Output template
		req.URL,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to download audio: %s", string(output)), http.StatusInternalServerError)
		return
	}

	// Find the downloaded file
	files, err := filepath.Glob(filepath.Join(tempDir, "audio.*"))
	if err != nil || len(files) == 0 {
		writeError(w, "Failed to find downloaded audio file", http.StatusInternalServerError)
		return
	}

	audioFile := files[0]

	// Read the file
	fileData, err := os.ReadFile(audioFile)
	if err != nil {
		writeError(w, "Failed to read audio file", http.StatusInternalServerError)
		return
	}

	// Get filename from URL or use default
	filename := extractFilename(req.URL)

	// Send the file
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.mp3\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))

	io.Copy(w, strings.NewReader(string(fileData)))
}

func isValidInstagramURL(url string) bool {
	// Match instagram.com/reel/ or instagram.com/reels/ or instagram.com/p/
	pattern := `^https?://(www\.)?instagram\.com/(reel|reels|p)/[A-Za-z0-9_-]+/?`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func extractFilename(url string) string {
	// Extract the reel ID from URL
	pattern := regexp.MustCompile(`instagram\.com/(?:reel|reels|p)/([A-Za-z0-9_-]+)`)
	matches := pattern.FindStringSubmatch(url)
	if len(matches) > 1 {
		return "instagram_audio_" + matches[1]
	}
	return "instagram_audio"
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
