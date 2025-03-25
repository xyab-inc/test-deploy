package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// CommandRunner defines the interface for running commands
type CommandRunner func(cmd string) (string, error)

// SSHClient handles SSH connections and operations
type SSHClient struct {
	client     *ssh.Client
	RunCommand CommandRunner
}

// CreateSSHClient creates a new SSH client with the given credentials
func CreateSSHClient(user, key, host string, port int) (*SSHClient, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	sshClient := &SSHClient{client: client}
	sshClient.RunCommand = sshClient.runCommand // Set default implementation
	return sshClient, nil
}

// runCommand executes a command on the remote server (internal implementation)
func (s *SSHClient) runCommand(cmd string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("SSH client is nil")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return string(output), nil
}

// TransferFile copies a local file to the remote server using SCP
// The remote path will be the base name of the local file
func (s *SSHClient) TransferFile(localPath string) error {
	return s.TransferFileWithRemotePath(localPath, filepath.Base(localPath))
}

// TransferFileWithRemotePath copies a local file to the remote server using SCP
// with a specified remote path
func (s *SSHClient) TransferFileWithRemotePath(localPath, remotePath string) error {
	if s.client == nil {
		return fmt.Errorf("SSH client is nil")
	}

	// Create a new session for file transfer
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Open local file
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer f.Close()

	// Get file info for size
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	// Start scp on the remote side
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C0644 %d %s\n", fi.Size(), remotePath)
		io.Copy(w, f)
		fmt.Fprint(w, "\x00")
	}()

	// Run scp command to receive file
	if err := session.Run("/usr/bin/scp -t ."); err != nil {
		return fmt.Errorf("failed to transfer file: %v", err)
	}

	return nil
}

// Close closes the SSH client connection
func (s *SSHClient) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
