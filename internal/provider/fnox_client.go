package provider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// fnoxClientConfig holds the resolved provider-level configuration used to
// invoke the fnox CLI.
type fnoxClientConfig struct {
	BinaryPath string
	ConfigPath string
	WorkingDir string
	Profile    string
}

// getSecretParams are the per-call parameters for fetching a single secret.
type getSecretParams struct {
	Key          string
	Profile      string
	Base64Decode bool
}

// getSecret runs `fnox get` and returns the trimmed secret value.
func (c fnoxClientConfig) getSecret(ctx context.Context, params getSecretParams) (string, error) {
	binary := c.BinaryPath
	if binary == "" {
		binary = "fnox"
	}

	args := []string{"get", "--non-interactive"}
	if c.ConfigPath != "" {
		args = append(args, "--config", c.ConfigPath)
	}

	profile := params.Profile
	if profile == "" {
		profile = c.Profile
	}
	if profile != "" {
		args = append(args, "--profile", profile)
	}

	if params.Base64Decode {
		args = append(args, "--base64-decode")
	}

	args = append(args, params.Key)

	cmd := exec.CommandContext(ctx, binary, args...)
	if c.WorkingDir != "" {
		cmd.Dir = c.WorkingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrText := strings.TrimSpace(stderr.String())
		if stderrText != "" {
			return "", fmt.Errorf("fnox get failed: %s", stderrText)
		}
		return "", fmt.Errorf("fnox get failed: %w", err)
	}

	return strings.TrimRight(stdout.String(), "\n"), nil
}
