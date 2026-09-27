package db

import (
	"database/sql"
	"log"
	"sync"

	"music-server/config"

	_ "github.com/mattn/go-sqlite3"
)

var (
	conn     *sql.DB
	connOnce sync.Once
)

// get_db は初回呼び出し時にだけ接続を開き、以降は同じ *sql.DB を返す
func get_db() *sql.DB {
	connOnce.Do(func() {
		var err error
		conn, err = sql.Open("sqlite3", config.DataPath("music.db"))
		if err != nil {
			log.Fatal(err)
		}
	})
	return conn
}

// Close はアプリ終了時に呼ぶ
func Close() {
	if conn != nil {
		conn.Close()
	}
}
