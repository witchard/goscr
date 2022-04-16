package main

import (
	"os"
	"testing"
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

func TestReadLock(t *testing.T) {
	t.Error("TODO: Test multiple read locks can be obtained, and can't get write lock when read lock in place")
}
