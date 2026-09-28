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
		path, err := download_one(bin, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		files[id] = path
	}

	return files, errors.Join(errs...)
}

// 試すフォーマットの順番（opus を優先し、それ以外は ffmpeg で opus に変換される）
var audio_formats = []string{
	"251", "250", "249", // webm / opus
	"140", "139", // m4a
	"234", "233", // m3u8
	"bestaudio",
}

func download_one(bin string, id string) (string, error) {
	path := filepath.Join(temp_dir, id+".opus")

	var last_err error
	for _, format := range audio_formats {
		args := append(js_runtime_args(),
			"-f", format,
			"-x", "--audio-format", "opus",
			"--no-playlist",
			"-o", filepath.Join(temp_dir, "%(id)s.%(ext)s"),
			"https://www.youtube.com/watch?v="+id,
		)
		cmd := exec.Command(bin, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			last_err = fmt.Errorf("Failed to download %s (format %s): %w\n%s", id, format, err, out)
			continue
		}

		if _, err := os.Stat(path); err != nil {
			last_err = fmt.Errorf("Cannot find downloaded file for %s (format %s): %w", id, format, err)
			continue
		}
		return path, nil
	}
	return "", last_err
}

func Clean_Temp() {
	if err := os.RemoveAll(temp_dir); err != nil {
		log.Printf("Failed to remove tmp directory: %v", err)
	}
}

func js_runtime_args() []string {
	if _, err := exec.LookPath("deno"); err == nil {
		return nil
	}
	if _, err := exec.LookPath("node"); err == nil {
		return []string{"--js-runtimes", "node"}
	}
	return nil
}
