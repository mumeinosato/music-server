package db

import (
	"log"
	"strings"
)

func Update_List(ids []string, hash string) {
	_, err := get_db().Exec(`INSERT INTO playlist (music_ids, hash) VALUES (?, ?)`,
		strings.Join(ids, ","), hash,
	)
	if err != nil {
		log.Fatal(err)
	}
}
