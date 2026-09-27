package config

import (
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerURL string
	Port      string

	YouTubePlaylistID string

	R2Token           string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2Endpoint        string
}

var (
	cfg     *Config
	cfgOnce sync.Once
)

// Get は初回呼び出し時にだけ .env を読み込み、以降は同じ *Config を返す
func Get() *Config {
	cfgOnce.Do(func() {
		godotenv.Load()

		server_url := strings.TrimRight(os.Getenv("SERVER_URL"), "/")
		if server_url == "" {
			server_url = "http://localhost:8080"
		}

		cfg = &Config{
			ServerURL: server_url,
			Port:      port_from_url(server_url),

			YouTubePlaylistID: os.Getenv("YOUTUBE_PLAYLIST_ID"),

			R2Token:           os.Getenv("CLOUDFLARE_R2_TOKEN"),
			R2AccessKeyID:     os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID"),
			R2SecretAccessKey: os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY"),
			R2BucketName:      os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
			R2Endpoint:        os.Getenv("CLOUDFLARE_R2_ENDPOINT"),
		}
	})
	return cfg
}

// port_from_url は SERVER_URL からポート番号を取り出す（省略時は 8080）
func port_from_url(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Port() == "" {
		return "8080"
	}
	return u.Port()
}
