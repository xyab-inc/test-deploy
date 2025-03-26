package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration structure
type Config struct {
	Hosts []Host `yaml:"hosts"`
}

// FindHost returns a pointer to the Host with the given name
// Returns nil and an error if no host is found
func (c *Config) FindHost(name string) (*Host, error) {
	for i := range c.Hosts {
		if c.Hosts[i].Name == name {
			return &c.Hosts[i], nil
		}
	}
	return nil, fmt.Errorf("host not found: %s", name)
}

// Host represents a single host configuration
type Host struct {
	Name          string        `yaml:"name"`
	IP            string        `yaml:"-"`       // Parsed from address
	Port          string        `yaml:"-"`       // Parsed from address
	address       string        `yaml:"address"` // Private field for raw address
	User          string        `yaml:"user"`
	DockerCompose DockerCompose `yaml:"docker_compose,omitempty"`
}

// UnmarshalYAML implements custom unmarshaling for Host to split address into IP and Port
func (h *Host) UnmarshalYAML(value *yaml.Node) error {
	// Create a temporary type without the custom unmarshaler
	type HostTemp struct {
		Name          string        `yaml:"name"`
		Address       string        `yaml:"address"`
		User          string        `yaml:"user"`
		DockerCompose DockerCompose `yaml:"docker_compose,omitempty"`
	}

	// Decode into the temporary struct
	var temp HostTemp
	if err := value.Decode(&temp); err != nil {
		return fmt.Errorf("decoding host: %w", err)
	}

	// Split address into IP and Port
	parts := strings.Split(temp.Address, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address format %q, expected IP:PORT", temp.Address)
	}

	// Copy all fields to the target struct
	h.Name = temp.Name
	h.User = temp.User
	h.DockerCompose = temp.DockerCompose
	h.IP = parts[0]
	h.Port = parts[1]
	h.address = temp.Address

	return nil
}

// MarshalJSON implements custom JSON marshaling to include the original address
func (h Host) MarshalJSON() ([]byte, error) {
	type HostAlias Host
	return json.Marshal(struct {
		HostAlias
		Address string `json:"address"`
	}{
		HostAlias: HostAlias(h),
		Address:   h.address,
	})
}

// DockerCompose represents Docker Compose configuration
type DockerCompose struct {
	Path string `yaml:"path"`
}

// GitHub represents the GitHub Actions environment
type GitHub struct {
	ConfigFile   string
	ValidateOnly bool
	HostName     string
}

func main() {
	gh := &GitHub{
		ConfigFile:   os.Getenv("CONFIG_FILE"),
		ValidateOnly: os.Getenv("VALIDATE_ONLY") == "true",
		HostName:     os.Getenv("HOST_NAME"),
	}

	if err := gh.run(); err != nil {
		fmt.Printf("::error::Failed to parse config: %v\n", err)
		os.Exit(1)
	}
}

// writeGitHubOutput writes the given string to the GITHUB_OUTPUT file
// Panics if unable to write or if GITHUB_OUTPUT env var is not set
func writeGitHubOutput(output string) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		panic("GITHUB_OUTPUT environment variable not set")
	}
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		panic(fmt.Sprintf("failed to write to GITHUB_OUTPUT: %v", err))
	}
}

func (gh *GitHub) run() error {
	// Get workspace path and ensure it exists
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		return fmt.Errorf("GITHUB_WORKSPACE environment variable not set")
	}

	// Clean the config file path and join with workspace
	configPath := filepath.Join(workspace, filepath.Clean(gh.ConfigFile))

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	// Parse YAML into Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}

	// Validate hosts
	if len(config.Hosts) == 0 {
		return fmt.Errorf("no hosts defined in configuration")
	}

	// Validate each host has required fields
	for _, host := range config.Hosts {
		if host.Name == "" {
			return fmt.Errorf("host missing required field: name - %v", host)
		}
		if host.address == "" {
			return fmt.Errorf("host %s missing required field: address", host.Name)
		}
		if host.User == "" {
			return fmt.Errorf("host %s missing required field: user", host.Name)
		}
	}

	// Calculate config hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// Find and output the requested host config
	host, err := config.FindHost(gh.HostName)
	if err != nil {
		writeGitHubOutput("host-found=false\n")
		return fmt.Errorf("finding host: %w", err)
	}

	// Validate docker compose path is present
	if host.DockerCompose.Path == "" {
		return fmt.Errorf("host %s is missing required docker compose path", host.Name)
	}

	// Convert host to JSON for consistent output format
	hostJSON, err := json.MarshalIndent(host, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling host config: %w", err)
	}

	// Write outputs
	output := fmt.Sprintf(
		"host-found=true\n"+
			"is-valid=true\n"+
			"config-hash=%s\n"+
			"ip=%s\n"+
			"port=%s\n"+
			"user=%s\n"+
			"compose-path=%s\n",
		hashStr,
		host.IP,
		host.Port,
		host.User,
		host.DockerCompose.Path,
	)
	writeGitHubOutput(output)

	fmt.Println("Found host configuration:")
	fmt.Println(string(hostJSON))

	return nil
}
