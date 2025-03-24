package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSSHClientCreation(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		key     string
		host    string
		port    int
		wantErr bool
	}{
		{
			name: "invalid_key",
			user: "testuser",
			key:  "invalid-key",
			host: "localhost",
			port: 22,
			wantErr: true,
		},
		{
			name: "empty_credentials",
			user: "",
			key:  "",
			host: "",
			port: 0,
			wantErr: true,
		},
		{
			name: "valid_key_but_invalid_host",
			user: "testuser",
			key: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAvRQk2oQqLB01iCnJuv0J6qEgMrLFPYChZZmykYgNQcxxjBVqFHn6
-----END OPENSSH PRIVATE KEY-----`,
			host: "nonexistent.host",
			port: 22,
			wantErr: true,
		},
		{
			name: "invalid_port",
			user: "testuser",
			key: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAvRQk2oQqLB01iCnJuv0J6qEgMrLFPYChZZmykYgNQcxxjBVqFHn6
-----END OPENSSH PRIVATE KEY-----`,
			host: "localhost",
			port: -1,
			wantErr: true,
		},
		{
			name: "malformed_key",
			user: "testuser",
			key: `-----BEGIN OPENSSH PRIVATE KEY-----
malformed
-----END OPENSSH PRIVATE KEY-----`,
			host: "localhost",
			port: 22,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := CreateSSHClient(tt.user, tt.key, tt.host, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSSHClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if client != nil {
				if err := client.Close(); err != nil {
					t.Errorf("Close() returned unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSSHClientMethods(t *testing.T) {
	// Test with nil client
	client := &SSHClient{client: nil}

	// Test RunCommand with nil client
	if _, err := client.RunCommand("test"); err == nil {
		t.Error("RunCommand() with nil client should return error")
	}

	// Test Close with nil client
	if err := client.Close(); err != nil {
		t.Error("Close() with nil client should not return error")
	}
}

func TestSSHClientTransferFile(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test with nil client
	client := &SSHClient{client: nil}
	if err := client.TransferFile(tmpFile); err == nil {
		t.Error("TransferFile() with nil client should return error")
	}

	// Test with nonexistent file
	validKey := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAvRQk2oQqLB01iCnJuv0J6qEgMrLFPYChZZmykYgNQcxxjBVqFHn6
-----END OPENSSH PRIVATE KEY-----`

	client, err := CreateSSHClient("testuser", validKey, "nonexistent.host", 22)
	if err == nil {
		t.Error("Expected error for nonexistent host")
		if err := client.Close(); err != nil {
			t.Errorf("Close() returned unexpected error: %v", err)
		}
	}

	// Test with invalid file path
	if client != nil {
		if err := client.TransferFile("/nonexistent/file"); err == nil {
			t.Error("Expected error for nonexistent file")
		}
	}
}

func TestValidateFiles(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test with nil client
	client := &SSHClient{client: nil}
	if err := validateFiles(client, ".env", "test.txt"); err == nil {
		t.Error("validateFiles() with nil client should return error")
	}

	// Test with mock client that simulates missing files
	mockClient := &SSHClient{
		client: nil,
		RunCommand: func(cmd string) (string, error) {
			// Simulate ls output with only one file
			return "-rw-r--r-- 1 user group 123 Mar 21 17:38 .env", nil
		},
	}

	if err := validateFiles(mockClient, ".env", "test.txt"); err == nil {
		t.Error("validateFiles() should return error when file is missing")
	} else if !strings.Contains(err.Error(), "test.txt") {
		t.Errorf("validateFiles() error should mention missing file, got: %v", err)
	}

	// Test with mock client that simulates all files present
	mockClient = &SSHClient{
		client: nil,
		RunCommand: func(cmd string) (string, error) {
			// Simulate ls output with both files
			return `-rw-r--r-- 1 user group 123 Mar 21 17:38 .env
-rw-r--r-- 1 user group 456 Mar 21 17:38 test.txt`, nil
		},
	}

	if err := validateFiles(mockClient, ".env", "test.txt"); err != nil {
		t.Errorf("validateFiles() returned unexpected error: %v", err)
	}
}
