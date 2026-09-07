package provider

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeFakeFnox writes a shell script named "fnox" to a temp dir that
// echoes its received arguments (and optionally fails), then returns that
// dir so it can be prepended to PATH.
func writeFakeFnox(t *testing.T, script string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake fnox script uses a shell shebang, unsupported on windows")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "fnox")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write fake fnox script: %v", err)
	}
	return dir
}

func withPath(t *testing.T, dir string) {
	t.Helper()
	original := os.Getenv("PATH")
	if err := os.Setenv("PATH", dir+string(os.PathListSeparator)+original); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PATH", original)
	})
}

func TestGetSecret_Success(t *testing.T) {
	script := `#!/bin/sh
echo "args: $@"
echo "hello-world"
`
	dir := writeFakeFnox(t, script)
	withPath(t, dir)

	client := fnoxClientConfig{}
	value, err := client.getSecret(context.Background(), getSecretParams{Key: "MY_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(value, "hello-world") {
		t.Fatalf("expected output to end with hello-world, got %q", value)
	}
}

func TestGetSecret_TrimsTrailingNewlineOnly(t *testing.T) {
	script := `#!/bin/sh
printf 'multi\nline\nvalue\n'
`
	dir := writeFakeFnox(t, script)
	withPath(t, dir)

	client := fnoxClientConfig{}
	value, err := client.getSecret(context.Background(), getSecretParams{Key: "MY_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "multi\nline\nvalue" {
		t.Fatalf("expected internal newlines preserved, got %q", value)
	}
}

func TestGetSecret_ErrorSurfacesStderr(t *testing.T) {
	script := `#!/bin/sh
echo "boom: secret not found" >&2
exit 1
`
	dir := writeFakeFnox(t, script)
	withPath(t, dir)

	client := fnoxClientConfig{}
	_, err := client.getSecret(context.Background(), getSecretParams{Key: "MISSING"})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "boom: secret not found") {
		t.Fatalf("expected error to contain stderr text, got %q", err.Error())
	}
}

func TestGetSecret_BuildsExpectedArgs(t *testing.T) {
	script := `#!/bin/sh
echo "$@" > "$FAKE_FNOX_ARGS_FILE"
echo "value"
`
	dir := writeFakeFnox(t, script)
	withPath(t, dir)

	argsFile := filepath.Join(t.TempDir(), "args.txt")
	t.Setenv("FAKE_FNOX_ARGS_FILE", argsFile)

	client := fnoxClientConfig{
		ConfigPath: "/tmp/fnox.toml",
		Profile:    "default-profile",
	}
	_, err := client.getSecret(context.Background(), getSecretParams{
		Key:          "MY_KEY",
		Profile:      "override-profile",
		Base64Decode: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("failed to read args file: %v", err)
	}

	gotStr := strings.TrimSpace(string(got))
	want := "get --non-interactive --config /tmp/fnox.toml --profile override-profile --base64-decode MY_KEY"
	if gotStr != want {
		t.Fatalf("unexpected args\n got: %q\nwant: %q", gotStr, want)
	}
}
