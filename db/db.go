package db

import (
	"database/sql"
	"sync"
	"spd/utils"	
	_ "modernc.org/sqlite"
)

var (
	once 	sync.Once
	instance *sql.DB
	initErr	error 
)

func GetInstance() (*sql.DB, error) {
	once.Do(func() {
		db, err := sql.Open("sqlite", utils.GetDBPath())
		if err != nil {
			initErr = err
			return
		}
		db.SetMaxOpenConns(1)
		if err := schema(db); err != nil {
			initErr = err
			return
		}
		instance = db
	})
	return instance, initErr
}
		