package yt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var temp_dir = "tmp"

func ytdlp_path() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join("yt-dlp", "yt-dlp.exe"), nil
	case "linux":
		return filepath.Join("yt-dlp", "yt-dlp_linux"), nil
	default:
		return "", fmt.Errorf("Failed to get yt-dlp path: unsupported OS: %s", runtime.GOOS)
	}
}

func Download(ids []string) (map[string]string, error) {
	bin, err := ytdlp_path()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(temp_dir, 0755); err != nil {
		return nil, fmt.Errorf("Failed to create tmp directory: %w", err)
	}

	files := make(map[string]string)
	var errs []error

	for _, id := range ids {
		cmd := exec.Command(bin,
			"-f", "bestaudio",
			"-x", "--audio-format", "opus",
			"--no-playlist",
			"-o", filepath.Join(temp_dir, "%(id)s.%(ext)s"),
			"https://www.youtube.com/watch?v="+id,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			errs = append(errs, fmt.Errorf("Failed to download %s: %w\n%s", id, err, out))
			continue
		}

		path := filepath.Join(temp_dir, id+".opus")
		if _, err := os.Stat(path); err != nil {
			errs = append(errs, fmt.Errorf("Cannot find downloaded file for %s: %w", id, err))
			continue
		}
		files[id] = path
	}

	return files, errors.Join(errs...)
}

func Clean_Temp() {
	if err := os.RemoveAll(temp_dir); err != nil {
		log.Printf("Failed to remove tmp directory: %v", err)
	}
}
