package botnet

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	_db, err := sql.Open("sqlite3", "sessions.db")
	if err != nil {
		return err
	}
	db = _db
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Main (Login TEXT, Password TEXT, NativeDeviceID TEXT, UID TEXT, SID TEXT, Nickname TEXT)")
	if err != nil {
		return err
	}
	return nil
}

func Read() ([]*Bot, error) {
	rows, err := db.Query("SELECT * FROM Main")
	if err != nil {
		return nil, err
	}
	var bots []*Bot
	for rows.Next() {
		bot := Bot{}
		err = rows.Scan(&bot.Login, &bot.Password, &bot.NativeDeviceID, &bot.UID, &bot.SID, &bot.Nickname)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &bot)
	}
	return bots, nil
}

func Write(account Bot) error {
	_, err := db.Exec("INSERT INTO Main VALUES ($1, $2, $3, $4, $5, $6)",
		account.Login, account.Password, account.NativeDeviceID, account.UID, account.SID, account.Nickname)
	if err != nil {
		return err
	}
	return nil
}
