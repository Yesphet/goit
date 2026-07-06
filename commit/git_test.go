package commit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitCommitRunsPreCommitHook(t *testing.T) {
	repoDir := initTestRepo(t)

	writeFile(t, "change.txt", "change\n")
	runGit(t, "add", "change.txt")

	hookPath := filepath.Join(repoDir, ".git", "hooks", "pre-commit")
	writeFile(t, hookPath, "#!/bin/sh\nprintf ran > .git/pre-commit-ran\n")
	if err := os.Chmod(hookPath, 0755); err != nil {
		t.Fatalf("chmod hook: %v", err)
	}

	gh, err := newGitHelper()
	if err != nil {
		t.Fatalf("newGitHelper: %v", err)
	}
	if err := gh.gitCommitPreCheck(); err != nil {
		t.Fatalf("gitCommitPreCheck: %v", err)
	}
	if err := gh.gitCommit("test: run hook\n"); err != nil {
		t.Fatalf("gitCommit: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repoDir, ".git", "pre-commit-ran")); err != nil {
		t.Fatalf("pre-commit hook did not run: %v", err)
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	chdir(t, repoDir)
	runGit(t, "init")
	runGit(t, "config", "user.name", "Goit Test")
	runGit(t, "config", "user.email", "goit@example.com")

	return repoDir
}

func chdir(t *testing.T, dir string) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func runGit(t *testing.T, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}
