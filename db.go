package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type ProgramLock struct {
	hash string
	db   *sql.DB
}

func (p *ProgramLock) Unlock() {
	if p != nil && p.db != nil {
		p.db.Exec(`UPDATE programs SET locked = ?, accessed = ? 
					WHERE hash = ? AND locked = ?`,
			false, time.Now().UTC(), p.hash, true)
		p.db.Close()
		p.db = nil
	}
}

func NewProgramLock(hash string) (*ProgramLock, error) {
	dbPath, err := Workdir()
	if err != nil {
		return nil, err
	}
	dbPath = filepath.Join(dbPath, "db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS programs (
		hash TEXT PRIMARY KEY,
		accessed DATETIME NOT NULL,
		locked BOOLEAN NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("INSERT INTO programs VALUES (?, ?, ?)",
		hash, time.Now().UTC(), true)
	if err != nil {
		// Already exists try and lock
		res, err := db.Exec(`UPDATE programs SET locked = ?, accessed = ? 
							WHERE hash = ? AND locked = ?`,
			true, time.Now().UTC(), hash, false)
		if err != nil {
			db.Close()
			return nil, err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			db.Close()
			return nil, err
		}
		if rows != 1 {
			return nil, fmt.Errorf("attempt to lock currently locked program")
		}
	}
	return &ProgramLock{hash, db}, nil
}
