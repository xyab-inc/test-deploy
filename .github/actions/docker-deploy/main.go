package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// log prints a message to stdout with GitHub Actions format
func log(msg string) {
	fmt.Printf("::notice::%s\n", msg)
}

// logError prints an error message to stdout with GitHub Actions format
func logError(msg string) {
	fmt.Printf("::error::%s\n", msg)
}

func createEnvFile(dockerTag string) (string, error) {
	// Create a temporary file
	tmpDir := os.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// Write the DOCKER_TAG to the .env file
	content := fmt.Sprintf("DOCKER_TAG=%s\n", dockerTag)
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create .env file: %v", err)
	}

	return envFile, nil
}

func validateFiles(client *SSHClient, files ...string) error {
	// Build ls command for all files
	cmd := fmt.Sprintf("ls -l %s", strings.Join(files, " "))
	output, err := client.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to validate files: %v", err)
	}

	// Check if any files are missing from the output
	for _, file := range files {
		if !strings.Contains(output, file) {
			return fmt.Errorf("file %s not found in remote directory", file)
		}
	}

	return nil
}

func main() {
	// Get SSH configuration from environment variables
	sshUser := os.Getenv("SSH_USER")
	sshKey := os.Getenv("SSH_KEY")
	sshHost := os.Getenv("SSH_HOST")
	sshPortStr := os.Getenv("SSH_PORT")
	composeFile := os.Getenv("COMPOSE_FILE")
	dockerTag := os.Getenv("DOCKER_TAG")

	if sshUser == "" || sshKey == "" || sshHost == "" || dockerTag == "" || composeFile == "" || sshPortStr == "" {
		logError("Missing required environment variables")
		os.Exit(1)
	}

	sshPort, err := strconv.Atoi(sshPortStr)
	if err != nil {
		logError(fmt.Sprintf("Invalid SSH port: %v", err))
		os.Exit(1)
	}

	// Create SSH client
	client, err := CreateSSHClient(sshUser, sshKey, sshHost, sshPort)
	if err != nil {
		logError(fmt.Sprintf("Failed to create SSH client: %v", err))
		os.Exit(1)
	}
	defer client.Close()

	// Create and transfer .env file with DOCKER_TAG
	envFile, err := createEnvFile(dockerTag)
	if err != nil {
		logError(fmt.Sprintf("Failed to create .env file: %v", err))
		os.Exit(1)
	}
	defer os.Remove(envFile) // Clean up temporary file

	if err := client.TransferFileWithRemotePath(envFile, ".env"); err != nil {
		logError(fmt.Sprintf("Failed to transfer .env file: %v", err))
		os.Exit(1)
	}
	log("Successfully transferred .env file")

	// Transfer docker-compose file
	remoteComposeFile := filepath.Base(composeFile)
	if err := client.TransferFileWithRemotePath(composeFile, remoteComposeFile); err != nil {
		logError(fmt.Sprintf("Failed to transfer docker-compose file: %v", err))
		os.Exit(1)
	}
	log(fmt.Sprintf("Successfully transferred docker-compose file: %s", remoteComposeFile))

	// Validate transferred files
	filesToValidate := []string{".env", remoteComposeFile}
	if err := validateFiles(client, filesToValidate...); err != nil {
		logError(fmt.Sprintf("File validation failed: %v", err))
		os.Exit(1)
	}
	log(fmt.Sprintf("Successfully validated files: %s", strings.Join(filesToValidate, ", ")))

	// Run docker compose pull
	log("Running docker compose pull...")
	pullCmd := fmt.Sprintf("docker compose -f %s pull", remoteComposeFile)
	pullOutput, err := client.RunCommand(pullCmd)
	if err != nil {
		logError(fmt.Sprintf("Failed to run docker compose pull: %v\nOutput: %s", err, pullOutput))
		os.Exit(1)
	}
	log(fmt.Sprintf("Successfully pulled Docker images:\n%s", pullOutput))

	// Run docker compose up -d
	log("Running docker compose up -d...")
	upCmd := fmt.Sprintf("docker compose -f %s up -d", remoteComposeFile)
	upOutput, err := client.RunCommand(upCmd)
	if err != nil {
		logError(fmt.Sprintf("Failed to run docker compose up: %v\nOutput: %s", err, upOutput))
		os.Exit(1)
	}
	log(fmt.Sprintf("Successfully started Docker containers:\n%s", upOutput))
}
