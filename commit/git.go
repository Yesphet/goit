package commit

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type gitHelper struct{}

func newGitHelper() (*gitHelper, error) {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %s", cleanGitOutput(out, err))
	}
	if strings.TrimSpace(string(out)) != "true" {
		return nil, fmt.Errorf("not inside a git work tree")
	}

	return &gitHelper{}, nil
}

func (gh *gitHelper) gitCommitPreCheck() error {
	out, err := exec.Command("git", "diff", "--cached", "--quiet", "--exit-code").CombinedOutput()
	if err == nil {
		return fmt.Errorf("nothing to commit, no changes added to commit")
	}

	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return nil
	}

	return fmt.Errorf("git diff --cached failed: %s", cleanGitOutput(out, err))
}

func (gh *gitHelper) gitCommit(msg string) error {
	file, err := os.CreateTemp("", "goit-commit-message-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := io.WriteString(file, msg); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	cmd := exec.Command("git", "commit", "-F", file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cleanGitOutput(out []byte, err error) string {
	text := strings.TrimSpace(string(out))
	if text != "" {
		return text
	}
	return err.Error()
}
