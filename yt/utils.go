package yt

import (
	"encoding/json"
	"log"
	"os"

	"music-server/config"

	"golang.org/x/oauth2"
)

const token_cache_file = "token.json"
const client_secret_file = "client_secret.json"

// get_file_path は data ディレクトリ内のファイルパスを返す
func get_file_path(name string) (string, error) {
	return config.DataPath(name), nil
}

func get_token(file string)(*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func save_token(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("トークンの保存に失敗しました: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetClientSecret() ([]byte, error) {
	path, err := get_file_path(client_secret_file)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}