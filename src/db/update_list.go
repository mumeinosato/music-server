package db

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
)

func Update_List(ids []string) {
	_, err := get_db().Exec(`INSERT INTO playlist (music_ids, hash) VALUES (?, ?)`,
		strings.Join(ids, ","), hash_list(ids),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func hash_list(ids []string) string {
	hasher := sha256.New()
	for _, id := range ids {
		hasher.Write([]byte(id))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func Add_Music_Name(ids []string, names []string) {
	if len(ids) == 0 {
		return
	}

	args := make([]any, len(ids)*2)
	for i, id := range ids {
		args[i*2] = id
		args[i*2+1] = names[i]
	}

	_, err := get_db().Exec(`
		INSERT OR REPLACE INTO music_name (music_id, name) VALUES `+strings.Repeat("(?, ?),", len(ids)-1)+`(?, ?)
	`, args...)
	if err != nil {
		log.Fatal(err)
	}
}

// Hash_List は Update_List が保存する hash と同じ方式で ids の hash を返す
func Hash_List(ids []string) string {
	return hash_list(ids)
}
