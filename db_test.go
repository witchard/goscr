package main

import (
	"os"
	"testing"
)

func TestDBLock(t *testing.T) {
	// Use temp dir for goscr files
	tmp := t.TempDir()
	os.Setenv("GOSCR_PATH", tmp)

	// Can lock new hash
	lockA, err := NewProgramLock("abc")
	if err != nil {
		t.Error("Couldn't get lock A", err)
	}
	if lockA == nil {
		t.Error("Lock A is nil")
	}
	defer lockA.Unlock()

	// Fail to lock already locked hash
	lockB, err := NewProgramLock("abc")
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

	// Sucessfully lock when other lock is released
	lockA.Unlock()
	lockC, err := NewProgramLock("abc")
	if err != nil {
		t.Error("Couldn't get lock C", err)
	}
	if lockC == nil {
		t.Error("Lock C is nil")
	}
	defer lockC.Unlock()
}
