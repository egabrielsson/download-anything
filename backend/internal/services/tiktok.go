package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DownloadTikTok downloads content from TikTok using yt-dlp
func DownloadTikTok(url string, downloadType DownloadType) (*DownloadResult, error) {
	tempDir, err := os.MkdirTemp("", "tiktok-download-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	var cmd *exec.Cmd
	var outputTemplate string
	var contentType string
	var fileExt string

	switch downloadType {
	case DownloadTypeAudio:
		outputTemplate = filepath.Join(tempDir, "audio.%(ext)s")
		cmd = exec.Command("yt-dlp",
			"-x",                    // Extract audio
			"--audio-format", "mp3", // Convert to mp3
			"--audio-quality", "0",  // Best quality
			"-o", outputTemplate,
			url,
		)
		contentType = "audio/mpeg"
		fileExt = "mp3"

	case DownloadTypeVideo:
		outputTemplate = filepath.Join(tempDir, "video.%(ext)s")
		cmd = exec.Command("yt-dlp",
			"-f", "best[ext=mp4]/best", // Best quality mp4
			"-o", outputTemplate,
			url,
		)
		contentType = "video/mp4"
		fileExt = "mp4"

	default:
		return nil, fmt.Errorf("invalid download type: %s", downloadType)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("yt-dlp failed: %s", string(output))
	}

	// Find the downloaded file
	var pattern string
	if downloadType == DownloadTypeAudio {
		pattern = filepath.Join(tempDir, "audio.*")
	} else {
		pattern = filepath.Join(tempDir, "video.*")
	}

	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to find downloaded file")
	}

	return &DownloadResult{
		FilePath:    files[0],
		Filename:    fmt.Sprintf("tiktok_%s.%s", downloadType, fileExt),
		ContentType: contentType,
	}, nil
}
