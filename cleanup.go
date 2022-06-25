package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Cleanup cleans up old compiled programs
func Cleanup() error {
	window, ok := os.LookupEnv("GOSCR_CLEANUP")
	if !ok {
		window = "90"
	}
	parsed, err := strconv.Atoi(window)
	if err != nil {
		return fmt.Errorf("GOSCR_CLEANUP (%s) invalid: %e", window, err)
	}
	cleanup := time.Duration(parsed*24) * time.Hour

	dbg.Println("Cleaning up programs older than", cleanup)

	entries, err := ListEntries()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, e := range entries {
		if now.Sub(e.accessed) > cleanup {
			dbg.Println("Cleaning up program", e.hash)
			// TODO delete it! Don't forget to check accessed again after lock is obtained
		}
	}
	return nil
}
