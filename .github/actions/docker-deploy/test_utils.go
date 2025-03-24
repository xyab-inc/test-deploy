package main

import (
	"os"
	"testing"
)

// For testing
var osExit = os.Exit

// Helper function to run code that might call os.Exit
func runWithExit(t testing.TB, f func()) (code int) {
	// Save current osExit and restore it after test
	origExit := osExit
	defer func() { osExit = origExit }()

	// Mock osExit
	code = 0
	osExit = func(c int) {
		code = c
		panic("exit")
	}

	// Run the function and catch any exit panic
	defer func() {
		if r := recover(); r != nil {
			if r != "exit" {
				panic(r)
			}
		}
	}()

	f()
	return code
}

// Helper to set up test environment
func setupTestEnv(t testing.TB) (cleanup func()) {
	// Save original environment
	origEnv := map[string]string{
		"GITHUB_OUTPUT": os.Getenv("GITHUB_OUTPUT"),
		"SSH_USER":     os.Getenv("SSH_USER"),
		"SSH_KEY":      os.Getenv("SSH_KEY"),
		"SSH_HOST":     os.Getenv("SSH_HOST"),
		"SSH_PORT":     os.Getenv("SSH_PORT"),
	}

	// Return cleanup function
	return func() {
		for k, v := range origEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}
}
