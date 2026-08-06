package ctrl

import (
	"database/sql"
	"log"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

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

func ReadDB() error {
	var (
		id, lastmsg, name string
		status            int
	)
	resDev := make(map[uuid.UUID]*Device)

	// 查所有
	rows, err := DB.Query(
		"SELECT id, name, status, lastmsg FROM devices",
	)
	if err != nil {
		return err
	}
	log.Println("db:loading devices")

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &name, &status, &lastmsg)
		fmt.Printf("%v:%v\n", id, name)

		uuid, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		msg, err := FromMessage([]byte(lastmsg))
		if err != nil {
			return err
		}

		resDev[uuid] = &Device{Name:name, LastMsg: msg, Status: false}
	}
	devices = resDev
	return nil
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
		log.Println("db:saving")

		_, err := DB.Exec(
			"INSERT OR REPLACE INTO devices (id, name, status, lastmsg) VALUES (?, ?, ?, ?)",
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
