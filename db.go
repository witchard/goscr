package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Lock interface {
	Unlock()
}

// WriteLock is an exclusive lock (readers or writers) on a program hash
type WriteLock struct {
	hash string
	db   *sql.DB
}

// ReadLock allows multiple reads on a program hash, but only when no writers
type ReadLock struct {
	hash string
	db   *sql.DB
}

func openDb() (*sql.DB, error) {
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
		readers INTEGER NOT NULL,
		writers INTEGER NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func LockWrite(hash string) (Lock, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO programs VALUES (?, ?, ?, ?)",
		hash, time.Now().UTC(), 0, 1)
	if err != nil {
		// Already exists try and lock
		res, err := db.Exec(`UPDATE programs SET writers = ?, accessed = ? 
							WHERE hash = ? AND writers = ? AND readers = ?`,
			1, time.Now().UTC(), hash, 0, 0)
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
			db.Close()
			return nil, fmt.Errorf("attempt to lock currently locked program")
		}
	}
	return &WriteLock{hash, db}, nil
}

func (l *WriteLock) Unlock() {
	if l != nil && l.db != nil {
		l.db.Exec(`UPDATE programs SET writers = ?, accessed = ? 
					WHERE hash = ? AND writers = ? AND readers = ?`,
			0, time.Now().UTC(), l.hash, 1, 0)
		l.db.Close()
		l.db = nil
	}
}

func LockRead(hash string) (Lock, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO programs VALUES (?, ?, ?, ?)",
		hash, time.Now().UTC(), 1, 0)
	if err != nil {
		// Already exists try and lock
		res, err := db.Exec(`UPDATE programs SET readers = readers + 1, accessed = ? 
						WHERE hash = ? AND writers = ? AND readers >= 0`,
			time.Now().UTC(), hash, 0)
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
			db.Close()
			return nil, fmt.Errorf("attempt to lock currently locked program")
		}
	}
	return &ReadLock{hash, db}, nil
}

func (l *ReadLock) Unlock() {
	if l != nil && l.db != nil {
		l.db.Exec(`UPDATE programs SET readers = readers - 1, accessed = ? 
					WHERE hash = ? AND readers > 0 AND writers = ?`,
			time.Now().UTC(), l.hash, 0)
		l.db.Close()
		l.db = nil
	}
}
