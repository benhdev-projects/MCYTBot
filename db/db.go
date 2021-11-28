package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB
var err error

func Init() {
	DB, err = sqlx.Connect("sqlite3", "/home/container/db.db")
	if err != nil {
		fmt.Printf("Error connecting to sqlite database: %v", err)
	}
	err := DB.MustExec("CREATE TABLE IF NOT EXISTS guilds (id BIGINT PRIMARY KEY, bans_enabled BOOLEAN DEFAULT FALSE)")
	fmt.Printf("Result: %v", err)
}
