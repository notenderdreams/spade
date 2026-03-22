package db

import (
	"database/sql"
	"encoding/json"
)

func GetAllChains() ([]Chain, error) {
	db, err := GetInstance()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, name FROM chains ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chains []Chain
	for rows.Next() {
		var c Chain
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		chains = append(chains, c)
	}

	for i, c := range chains {
		steps, err := getChainSteps(db, c.ID)
		if err != nil {
			return nil, err
		}
		chains[i].Steps = steps
	}

	return chains, nil
}

func GetChain(name string) (*Chain, error) {
	db, err := GetInstance()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT id, name FROM chains WHERE name = ?", name)
	var c Chain
	if err := row.Scan(&c.ID, &c.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.Steps, err = getChainSteps(db, c.ID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func AddChain(name string, scriptNames []string) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO chains (name) VALUES (?)", name)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()

	for i, s := range scriptNames {
		if _, err := tx.Exec(
			"INSERT INTO chain_steps (chain_id, script_name, seq) VALUES (?, ?, ?)",
			id, s, i,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func DeleteChain(name string) (bool, error) {
	db, err := GetInstance()
	if err != nil {
		return false, err
	}

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id FROM chains WHERE name = ?", name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// chain_steps deleted automatically via ON DELETE CASCADE
	if _, err := tx.Exec("DELETE FROM chains WHERE id = ?", id); err != nil {
		return false, err
	}

	return true, tx.Commit()
}

func getChainSteps(db *sql.DB, chainID int) ([]ChainStep, error) {
	rows, err := db.Query(`
		SELECT cs.seq, s.name, s.cmd, s.args, COALESCE(s.runner, '')
		FROM chain_steps cs
		JOIN scripts s ON cs.script_name = s.name
		WHERE cs.chain_id = ?
		ORDER BY cs.seq
	`, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []ChainStep
	for rows.Next() {
		var step ChainStep
		var args string
		if err := rows.Scan(&step.Seq, &step.Script.Name, &step.Script.Command, &args, &step.Script.Runner); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(args), &step.Script.Args); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func RenameChain(oldName, newName string) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE chains SET name = ? WHERE name = ?", newName, oldName)
	return err
}

func AppendChainStep(chainName, scriptName string) error {
	db, err := GetInstance()
	if err != nil {
		return err
	}

	row := db.QueryRow("SELECT id FROM chains WHERE name = ?", chainName)
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	var maxSeq int
	row = db.QueryRow("SELECT COALESCE(MAX(seq), -1) FROM chain_steps WHERE chain_id = ?", id)
	if err := row.Scan(&maxSeq); err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO chain_steps (chain_id, script_name, seq) VALUES (?, ?, ?)",
		id, scriptName, maxSeq+1,
	)
	return err
}
