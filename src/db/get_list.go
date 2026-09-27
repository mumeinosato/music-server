package db

import (
	"database/sql"
	"log"
	"strings"
)

type mode int

const (
	mode_latest mode = iota
	mode_from_hash
)

var Mode = struct {
	Latest    mode
	From_Hash mode
}{
	Latest:    mode_latest,
	From_Hash: mode_from_hash,
}

func Get_List(m mode, hash ...string) ([]string, string) {
	switch m {
	case mode_latest:
		return get_latest()
	case mode_from_hash:
		if len(hash) == 0 {
			log.Fatal("hash is required for From_Hash mode")
		}
		return from_hash(hash[0])
	default:
		log.Fatal("invalid mode")
		return nil, ""
	}
}

func get_latest() ([]string, string) {
	row := get_db().QueryRow(`
		SELECT music_ids, hash FROM playlist ORDER BY id DESC LIMIT 1
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
