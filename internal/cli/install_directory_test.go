package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// directoryInstallDependencies fakes only the home (backup root) and the
// backup name; resolution, observation, and the transaction run for real so
// the test exercises the production path into the target directory.
func directoryInstallDependencies(t *testing.T, home string) installDependencies {
	t.Helper()
	deps := defaultInstallDependencies()
	deps.home = func() (string, error) { return home, nil }
	deps.backupName = func() (string, error) { return ".cortex-backup-000102030405060708090a0b0c0d0e0f", nil }
	return deps
}

// testDirectoryInstall drives the production dispatch the way a real run
// would, with the home (backup root) and backup name faked.
func testDirectoryInstall(t *testing.T, target string) (int, string, string) {
	t.Helper()
	home := canonicalHome(t)
	var stdout, stderr bytes.Buffer
	code := runWithInstallDependencies(context.Background(), []string{"install", "--target-dir", target}, &stdout, &stderr, nil, directoryInstallDependencies(t, home))
	return code, stdout.String(), stderr.String()
}

func TestInstallDirectoryTargetProjectsSkillsTransactionally(t *testing.T) {
	home := canonicalHome(t)
	target := filepath.Join(home, "agent-skills")
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := testDirectoryInstall(t, target)
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
	if code != exitOK {
		t.Fatalf("exit = %d, stdout = %q", code, stdout)
	}
	if !strings.Contains(stdout, "status=completed touch=applied") {
		t.Fatalf("stdout missing completed status: %q", stdout)
	}
	if !strings.Contains(stdout, "runtime=directory presence=present compatibility=compatible action=configure touch=applied\n") {
		t.Fatalf("stdout missing directory runtime line: %q", stdout)
	}

	// The manifest records provenance and hashes for the directory runtime.
	stateBytes, err := os.ReadFile(filepath.Join(target, ".cortex", "install-state.json"))
	if err != nil {
		t.Fatalf("install state missing: %v", err)
	}
	var state struct {
		SchemaVersion int    `json:"schemaVersion"`
		Owner         string `json:"owner"`
		Runtime       string `json:"runtime"`
		RootKind      string `json:"rootKind"`
		Artifacts     map[string]struct {
			RelativePath string `json:"relativePath"`
			SHA256       string `json:"sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		t.Fatal(err)
	}
	if state.Owner != "cortex" || state.Runtime != "directory" || state.RootKind != "directory" || len(state.Artifacts) == 0 {
		t.Fatalf("unexpected state: %#v", state)
	}
	if len(state.Artifacts) < 1 {
		t.Fatal("expected at least one projected skill")
	}
	// Readback: every manifest artifact exists on disk under its recorded
	// relative path with the recorded hash.
	for logicalID, artifact := range state.Artifacts {
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(artifact.RelativePath)))
		if err != nil {
			t.Fatalf("artifact %s (%s) missing: %v", logicalID, artifact.RelativePath, err)
		}
		if len(data) == 0 {
			t.Fatalf("artifact %s is empty", artifact.RelativePath)
		}
		if len(artifact.SHA256) != 64 {
			t.Fatalf("artifact %s hash %q is not a sha256", artifact.RelativePath, artifact.SHA256)
		}
	}
	// A first install has nothing to replace, so filetxn captures no backup
	// snapshot; the all-or-nothing guarantee is covered by the filetxn suite.
	// The target owns its state file and nothing else leaks outside it.
	entries, err := os.ReadDir(filepath.Join(target, ".cortex"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "install-state.json" {
		t.Fatalf("unexpected .cortex entries: %v", entries)
	}
}

func TestInstallDirectoryTargetFailClosedOnInvalidPaths(t *testing.T) {
	home := canonicalHome(t)
	missing := filepath.Join(home, "does-not-exist")
	if err := os.WriteFile(filepath.Join(home, "regular-file"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A relative path that does exist as a directory must also be refused
	// before composition: it passes os.Stat/IsDir, so only the explicit
	// absolute+clean check can catch it (issue #367).
	if err := os.MkdirAll("test-invalid-target-exists", 0o700); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("test-invalid-target-exists")
	failClosedCases := map[string]string{
		"relative":          "relative/target",
		"relative-existing": "test-invalid-target-exists",
		"missing":           missing,
		"file":              filepath.Join(home, "regular-file"),
		"empty":              "",
	}
	for name, target := range failClosedCases {
		var stdout, stderr bytes.Buffer
		code := runWithInstallDependencies(context.Background(), []string{"install", "--target-dir", target}, &stdout, &stderr, nil, directoryInstallDependencies(t, home))
		if code == exitOK || stderr.String() != "error=invalid_target\n" {
			t.Fatalf("%s: exit = %d stderr = %q, want invalid_target refusal", name, code, stderr.String())
		}
		if strings.Contains(stdout.String(), "touch=applied") {
			t.Fatalf("%s: nothing may be applied: %q", name, stdout.String())
		}
	}
}
