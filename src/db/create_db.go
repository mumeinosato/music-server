package db

import (
	"log"
)

// テーブルごとに存在チェックして、ないものだけ作る
var tables = []struct {
	name   string
	schema string
}{
	{"playlist", `
		CREATE TABLE playlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			music_ids TEXT NOT NULL,
			hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`},
	{"music_name", `
		CREATE TABLE music_name (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			music_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL
		)
	`},
}

func Create_DB() {
	for _, t := range tables {
		if table_exists(t.name) {
			continue
		}
		if _, err := get_db().Exec(t.schema); err != nil {
			log.Fatal(err)
		}
	}
}

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
