package db

import (
	"database/sql"
	"encoding/json"

	"spd/utils"

	_ "modernc.org/sqlite"
)

func getConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", utils.GetDBPath())
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scripts (
		name 	TEXT PRIMARY KEY,
		cmd		TEXT ,
		args 	TEXT	
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetAllScripts() ([]Script, error) {
	db, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM scripts")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []Script
	for rows.Next() {
		var s Script
		var args string
		err := rows.Scan(&s.Name, &s.Command, &args)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(args), &s.Args)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, s)
	}
	return scripts, nil
}

func GetScript(name string) (*Script, error) {
	db, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scripts WHERE name = ?", name)
	var s Script
	var args string
	err = row.Scan(&s.Name, &s.Command, &args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(args), &s.Args)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func AddScript(s Script) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	args, err := json.Marshal(s.Args)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO scripts (name, cmd, args) VALUES (?, ?, ?)", s.Name, s.Command, string(args))
	return err
}

func UpdateScript(s Script) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	args, err := json.Marshal(s.Args)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE scripts SET cmd = ?, args = ? WHERE name = ?", s.Command, string(args), s.Name)
	return err
}

func DeleteScript(name string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM scripts WHERE name = ?", name)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func RenameScript(oldName, newName string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE scripts SET name = ? WHERE name = ?", newName, oldName)
	return err
}