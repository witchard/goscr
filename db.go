package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// WriteLock is an exclusive lock (readers or writers) on a program hash
type WriteLock struct {
	hash     string
	db       *sql.DB
	accessed time.Time // Last accessed time
}

// ReadLock allows multiple reads on a program hash, but only when no writers
type ReadLock struct {
	hash string
	db   *sql.DB
}

// Entry represents an entry in the program database
type Entry struct {
	hash     string
	accessed time.Time
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

// ListEntries lists all known programs
func ListEntries() ([]Entry, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT hash, accessed FROM programs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.hash, &e.accessed); err != nil {
			return nil, err
		}
		ret = append(ret, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

// LockWrite locks a compilation hash for writing into its directory
func LockWrite(hash string) (*WriteLock, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}

	accessed := time.Now().UTC()
	_, err = db.Exec("INSERT INTO programs VALUES (?, ?, ?, ?)",
		hash, accessed, 0, 1)
	if err != nil {
		// Retrieve last accessed time
		err := db.QueryRow(`SELECT accessed FROM programs
							WHERE hash = ?`,
			hash).Scan(&accessed)
		if err != nil {
			db.Close()
			return nil, err
		}

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
	return &WriteLock{hash, db, accessed}, nil
}

// Unlock releases the lock
func (l *WriteLock) Unlock() {
	if l != nil && l.db != nil {
		l.db.Exec(`UPDATE programs SET writers = ?, accessed = ? 
					WHERE hash = ? AND writers = ? AND readers = ?`,
			0, time.Now().UTC(), l.hash, 1, 0)
		l.db.Close()
		l.db = nil
	}
}

// Accessed retrieves the last time at which this element was accessed (prior to obtaining the lock)
func (l *WriteLock) Accessed() time.Time {
	return l.accessed
}

// Delete removes this entry (and therefore releases lock)
func (l *WriteLock) Delete() {
	if l != nil && l.db != nil {
		l.db.Exec(`DELETE FROM programs
					WHERE hash = ? AND writers = ? AND readers = ?`,
			l.hash, 1, 0)
		l.db.Close()
		l.db = nil
	}
}

// LockRead locks a compilation hash for reading from its directory
func LockRead(hash string) (*ReadLock, error) {
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

// Unlock releases the lock
func (l *ReadLock) Unlock() {
	if l != nil && l.db != nil {
		l.db.Exec(`UPDATE programs SET readers = readers - 1, accessed = ? 
					WHERE hash = ? AND readers > 0 AND writers = ?`,
			time.Now().UTC(), l.hash, 0)
		l.db.Close()
		l.db = nil
	}
}
