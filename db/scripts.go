package db

import (
	"database/sql"
	"encoding/json"
)

func GetAllScripts() ([]Script, error) {
	db, err := GetInstance()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name, cmd, args, COALESCE(runner, '') FROM scripts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []Script
	for rows.Next() {
		var s Script
		var args string
		if err := rows.Scan(&s.Name, &s.Command, &args, &s.Runner); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(args), &s.Args); err != nil {
			return nil, err
		}
		scripts = append(scripts, s)
	}
	return scripts, nil
}

func GetScript(name string) (*Script, error) {
	db, err := GetInstance()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT name, cmd, args, COALESCE(runner, '') FROM scripts WHERE name = ?", name)
	var s Script
	var args string
	if err := row.Scan(&s.Name, &s.Command, &args, &s.Runner); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal([]byte(args), &s.Args); err != nil {
		return nil, err
	}
	return &s, nil
}

func AddScript(s Script) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}

	args, err := json.Marshal(s.Args)
	if err != nil {
		return err
	}

	var runner *string
	if s.Runner != "" {
		runner = &s.Runner
	}

	_, err = db.Exec(
		"INSERT INTO scripts (name, cmd, args, runner) VALUES (?, ?, ?, ?)",
		s.Name, s.Command, string(args), runner,
	)
	return err
}

func UpdateScript(s Script) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}

	args, err := json.Marshal(s.Args)
	if err != nil {
		return err
	}

	var runner *string
	if s.Runner != "" {
		runner = &s.Runner
	}

	_, err = db.Exec(
		"UPDATE scripts SET cmd = ?, args = ?, runner = ? WHERE name = ?",
		s.Command, string(args), runner, s.Name,
	)
	return err
}

func DeleteScript(name string) (bool, error) {
	db, err := GetInstance()
	if err != nil {
		return false, err
	}

	res, err := db.Exec("DELETE FROM scripts WHERE name = ?", name)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

func RenameScript(oldName, newName string) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE scripts SET name = ? WHERE name = ?", newName, oldName)
	return err
}
