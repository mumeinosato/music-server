package db

import (
	"database/sql"
	"log"
	"strings"
)

func get_latest() ([]string, string) {
	row := get_db().QueryRow(`
		SELECT music_ids, hash FROM playlist ORDER BY created_at DESC LIMIT 1
	`)
	return scan_playlist(row)
}

func from_hash(hash string) ([]string, string) {
	row := get_db().QueryRow(`
		SELECT music_ids, hash FROM playlist WHERE hash = ? LIMIT 1
	`, hash)
	return scan_playlist(row)
}

// scan_playlist は1行を読み取り、該当なしなら nil, "" を返す
func scan_playlist(row *sql.Row) ([]string, string) {
	var music_ids string
	var hash string
	err := row.Scan(&music_ids, &hash)
	if err == sql.ErrNoRows {
		return nil, ""
	}
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(music_ids, ","), hash
}
