package ctrl

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite" // 纯 Go，不需要 CGO

	. "sleeply-alive/internal/model"
)

var Path string
var DB *sql.DB

func OpenDB(path string) {
	Path = filepath.Join(path, "device.db")

	db, err := sql.Open("sqlite", Path)
	if err != nil {
		log.Fatal(err)
	}

	// 创建表
	db.Exec(`CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT,
		status INTEGER,
		lastmsg TEXT
	)`)
	DB = db
}

func SaveDB(dev *Device) {
	var (
		id, lastmsg, name string
		status            int
	)

	if dev == nil {
		log.Println("db:device is nil,cant save")
		return
	}

	if dev.LastMsg.UUID != nil {
		id = dev.LastMsg.UUID.String()
		m, _ := dev.LastMsg.ToJSON()
		lastmsg = string(m)
		name = dev.Name
		if dev.Status {
			status = 1
		} else {
			status = 0
		}

		_, err := DB.Exec(
			"INSERT INTO devices (id, name, status, lastmsg) VALUES (?, ?, ?, ?)",
			id, name, status, lastmsg,
		)
		if err != nil {
			log.Printf("db:%v", err)
			return
		}
	}
}

func CloseDB() {
	DB.Close()
}
