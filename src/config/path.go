package config

import (
	"log"
	"os"
	"path/filepath"
)

// DataDir は自動生成ファイルや認証情報を置くディレクトリ
const DataDir = "data"

// DataPath は data ディレクトリ内のファイルパスを返す（ディレクトリがなければ作成する）
func DataPath(name string) string {
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		log.Fatalf("data ディレクトリの作成に失敗しました: %v", err)
	}
	return filepath.Join(DataDir, name)
}
