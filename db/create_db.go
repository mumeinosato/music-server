package db

import (
	"log"
)

func Create_DB() {
	if table_exists("playlist") {
		return
	}

	_, err := get_db().Exec(`
		CREATE TABLE IF NOT EXISTS playlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			music_ids TEXT NOT NULL,
			hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}

// table_exists は指定したテーブルが既に存在するかを返す
func table_exists(name string) bool {
	var count int
	err := get_db().QueryRow(`
		SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?
	`, name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
