package main

// func TestMain(t *testing.T) {
// 	// Clear all SSH environment variables first
// 	envVars := []string{"SSH_USER", "SSH_KEY", "SSH_HOST", "SSH_PORT", "SCP_FILE"}
// 	for _, env := range envVars {
// 		t.Setenv(env, "")
// 	}

// 	// Save original args
// 	origArgs := os.Args
// 	defer func() { os.Args = origArgs }()

// 	// Test missing required env vars
// 	os.Args = []string{"cmd"}
// 	exitCode := runWithExit(t, func() {
// 		main()
// 	})
// 	if exitCode != 1 {
// 		t.Error("Expected exit code 1 for missing env vars")
// 	}

// 	// Test with invalid SSH key
// 	validKey := `-----BEGIN OPENSSH PRIVATE KEY-----
// b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
// NhAAAAAwEAAQAAAQEAvRQk2oQqLB01iCnJuv0J6qEgMrLFPYChZZmykYgNQcxxjBVqFHn6
// -----END OPENSSH PRIVATE KEY-----`

// 	t.Setenv("SSH_USER", "testuser")
// 	t.Setenv("SSH_KEY", validKey)
// 	t.Setenv("SSH_HOST", "localhost")
// 	t.Setenv("SSH_PORT", "22")

// 	exitCode = runWithExit(t, func() {
// 		main()
// 	})
// 	if exitCode != 1 {
// 		t.Error("Expected exit code 1 for invalid SSH key")
// 	}
// }

// func TestMainWithInvalidPort(t *testing.T) {
// 	// Clear all SSH environment variables first
// 	envVars := []string{"SSH_USER", "SSH_KEY", "SSH_HOST", "SSH_PORT", "SCP_FILE"}
// 	for _, env := range envVars {
// 		t.Setenv(env, "")
// 	}

// 	// Save original args
// 	origArgs := os.Args
// 	defer func() { os.Args = origArgs }()

// 	// Test invalid port
// 	validKey := `-----BEGIN OPENSSH PRIVATE KEY-----
// b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
// NhAAAAAwEAAQAAAQEAvRQk2oQqLB01iCnJuv0J6qEgMrLFPYChZZmykYgNQcxxjBVqFHn6
// -----END OPENSSH PRIVATE KEY-----`

// 	t.Setenv("SSH_USER", "testuser")
// 	t.Setenv("SSH_KEY", validKey)
// 	t.Setenv("SSH_HOST", "localhost")
// 	t.Setenv("SSH_PORT", "invalid")

// 	exitCode := runWithExit(t, func() {
// 		main()
// 	})
// 	if exitCode != 1 {
// 		t.Error("Expected exit code 1 for invalid port")
// 	}
// }
