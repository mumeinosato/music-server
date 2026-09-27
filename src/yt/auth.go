package yt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"music-server/src/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// redirect_url は Google 認証後に戻ってくるURL（gin の /callback ルートに対応）
func redirect_url() string {
	return config.Get().ServerURL + "/callback"
}

// AuthRequiredError は認証が必要なときに返され、認証URLを持つ
type AuthRequiredError struct {
	URL string
}

func (e *AuthRequiredError) Error() string {
	return "Google アカウントの認証が必要です: " + e.URL
}

var (
	pending_state string
	state_mu      sync.Mutex
)

func get_oauth_config() (*oauth2.Config, error) {
	b, err := GetClientSecret()
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, err
	}
	config.RedirectURL = redirect_url()
	return config, nil
}

// GetClient は保存済みトークンでクライアントを返す。トークンがなければ AuthRequiredError を返す
func GetClient(config *oauth2.Config) (*http.Client, error) {
	file, err := get_file_path(token_cache_file)
	if err != nil {
		return nil, err
	}
	token, err := get_token(file)
	if err != nil {
		return nil, new_auth_required(config)
	}
	return config.Client(context.Background(), token), nil
}

// new_auth_required は新しい state を発行し、認証URL付きのエラーを返す
func new_auth_required(config *oauth2.Config) error {
	buf := make([]byte, 16)
	rand.Read(buf)
	state := hex.EncodeToString(buf)

	state_mu.Lock()
	pending_state = state
	state_mu.Unlock()

	// prompt=consent で毎回 refresh token を受け取る
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	return &AuthRequiredError{URL: url}
}

// HandleCallback は認証後のコールバックを処理し、トークンを保存する
func HandleCallback(code string, state string) error {
	state_mu.Lock()
	expected := pending_state
	pending_state = ""
	state_mu.Unlock()

	if expected == "" || state != expected {
		return fmt.Errorf("state が一致しません。もう一度 /sync から認証してください")
	}
	if code == "" {
		return fmt.Errorf("認証コードが見つかりませんでした")
	}

	config, err := get_oauth_config()
	if err != nil {
		return err
	}
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("トークンの取得に失敗しました: %w", err)
	}

	file, err := get_file_path(token_cache_file)
	if err != nil {
		return err
	}
	save_token(file, token)
	return nil
}

// check_token_error はトークンが失効している場合、トークンを削除して AuthRequiredError に変換する
func check_token_error(config *oauth2.Config, err error) error {
	var retrieve_err *oauth2.RetrieveError
	if errors.As(err, &retrieve_err) {
		if file, ferr := get_file_path(token_cache_file); ferr == nil {
			os.Remove(file)
		}
		return new_auth_required(config)
	}
	return err
}
