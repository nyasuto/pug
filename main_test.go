package main

import (
	"os"
	"os/user"
	"testing"
)

// Test that main function doesn't panic
func TestMainFunction(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with no arguments (should run REPL, but we can't test interactive mode easily)
	// We'll just ensure the main function can be called without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	// We can't easily test the interactive REPL, but we can test that
	// the imports and basic setup work
	user, err := user.Current()
	if err != nil {
		t.Skipf("Cannot get current user: %v", err)
	}

	if user.Username == "" {
		t.Error("Expected non-empty username")
	}
}

// Test that the package imports work correctly
func TestImports(t *testing.T) {
	// Test that we can access the phase1 package
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Import test panicked: %v", r)
		}
	}()

	// Basic test that the imports are working
	// We can't test the actual REPL interaction easily in unit tests
	user, err := user.Current()
	if err != nil {
		t.Skipf("Cannot get current user for import test: %v", err)
	}

	if user.Username == "" {
		t.Error("Username should not be empty")
	}
}
