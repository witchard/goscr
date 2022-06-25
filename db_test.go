package main

import (
	"os"
	"testing"
	"time"
)

func TestWriteLock(t *testing.T) {
	// Use temp dir for goscr files
	tmp := t.TempDir()
	os.Setenv("GOSCR_PATH", tmp)

	// Can lock new hash
	lockA, err := LockWrite("abc")
	if err != nil {
		t.Error("Couldn't get lock A", err)
	}
	if lockA == nil {
		t.Error("Lock A is nil")
	}
	defer lockA.Unlock()

	// Fail to lock already locked hash
	lockB, err := LockWrite("abc")
	if lockB != nil {
		lockB.Unlock()
		t.Error("Got lock B and should not have done")
	}
	if err == nil {
		t.Error("Should have received error obtaining lock B")
	}
	if err.Error() != "attempt to lock currently locked program" {
		t.Error("Unexpected error obtaining lock B", err)
	}

	// Fail to read lock already locked hash
	read, err := LockRead("abc")
	if read != nil {
		read.Unlock()
		t.Error("Got read lock and should not have done")
	}
	if err == nil {
		t.Error("Should have received error obtaining read lock")
	}
	if err.Error() != "attempt to lock currently locked program" {
		t.Error("Unexpected error obtaining read lock", err)
	}

	// Sucessfully lock when other lock is released
	lockA.Unlock()
	lockC, err := LockWrite("abc")
	if err != nil {
		t.Error("Couldn't get lock C", err)
	}
	if lockC == nil {
		t.Error("Lock C is nil")
	}
	defer lockC.Unlock()
}

func TestWriteUpdated(t *testing.T) {
	// Use temp dir for goscr files
	tmp := t.TempDir()
	os.Setenv("GOSCR_PATH", tmp)

	// Lock new hash
	before := time.Now().UTC()
	lock, err := LockWrite("abc")
	if err != nil || lock == nil {
		t.Error("Couldn't get lock", err)
	}
	defer lock.Unlock()
	after := time.Now().UTC()

	// List database to confirm correct entries
	entries, err := ListEntries()
	if err != nil {
		t.Error("Failed to list entries", err)
	}
	if len(entries) != 1 {
		t.Error("Incorrect number of entries", entries)
	}
	if entries[0].hash != "abc" {
		t.Error("Wrong entry found", entries[0])
	}

	// Check accessed time is between before & after
	accessed := lock.Accessed()
	if before.After(accessed) || after.Before(accessed) {
		t.Errorf("Invalid updated times %s %s %s", before, accessed, after)
	}
	if accessed != entries[0].accessed {
		t.Error("Accessed time missmatch", accessed, entries[0].accessed)
	}

	// Re-obtain lock and verify accessed time unchanged
	before = time.Now().UTC()
	lock.Unlock() // Releasing lock updates access time
	after = time.Now().UTC()
	lock, err = LockWrite("abc")
	if err != nil || lock == nil {
		t.Error("Couldn't get lock", err)
	}
	defer lock.Unlock()
	accessed = lock.Accessed()
	if before.After(accessed) || after.Before(accessed) {
		t.Errorf("Invalid updated times %s %s %s", before, accessed, after)
	}
}

func TestWriteDelete(t *testing.T) {
	// Use temp dir for goscr files
	tmp := t.TempDir()
	os.Setenv("GOSCR_PATH", tmp)

	// Lock new hash
	lock, err := LockWrite("abc")
	if err != nil || lock == nil {
		t.Error("Couldn't get lock", err)
	}
	defer lock.Unlock()

	// Delete it and make sure its gone
	lock.Delete()

	entries, err := ListEntries()
	if err != nil {
		t.Error("Failed to list entries", err)
	}
	if len(entries) != 0 {
		t.Error("Still entries in database", entries)
	}
}

func TestReadLock(t *testing.T) {
	// Use temp dir for goscr files
	tmp := t.TempDir()
	os.Setenv("GOSCR_PATH", tmp)

	// Can obtain multiple read locks
	lockA, err := LockRead("abc")
	if err != nil {
		t.Error("Couldn't get lock A", err)
	}
	if lockA == nil {
		t.Error("Lock A is nil")
	}
	defer lockA.Unlock()
	lockB, err := LockRead("abc")
	if err != nil {
		t.Error("Can not obtain read lock B")
	}
	if lockB == nil {
		t.Error("Lock B is nil")
	}
	defer lockB.Unlock()
	lockC, err := LockRead("abc")
	if err != nil {
		t.Error("Can not obtain read lock C")
	}
	if lockC == nil {
		t.Error("Lock C is nil")
	}
	defer lockC.Unlock()

	// Fail to write lock when holding read locks
	lockD, err := LockWrite("abc")
	if lockD != nil {
		lockD.Unlock()
		t.Error("Got lock D and should not have done")
	}
	if err == nil {
		t.Error("Should have received error obtaining lock D")
	}
	if err.Error() != "attempt to lock currently locked program" {
		t.Error("Unexpected error obtaining lock D", err)
	}

	// Can obtain write lock when read locks are released
	lockA.Unlock()
	lockB.Unlock()
	lockC.Unlock()
	lockE, err := LockWrite("abc")
	if err != nil {
		t.Error("Couldn't get lock E", err)
	}
	if lockE == nil {
		t.Error("Lock E is nil")
	}
	defer lockE.Unlock()
}
