package commit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type gitHelper struct {
	repo      *git.Repository
	userName  string
	userEmail string
}

func newGitHelper() (*gitHelper, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}
	gh := &gitHelper{
		repo: repo,
	}

	gh.userName, gh.userEmail = getUserConfig(repo)
	if gh.userName == "" || gh.userEmail == "" {
		return nil, fmt.Errorf("can't get user or email of git")
	}
	return gh, nil
}

// get user and email from repository config, if not exist, read from ~/.gitconfig
func getUserConfig(repo *git.Repository) (name, email string) {
	name, email = readUserFromGlobalGitConfig()

	cfg, err := repo.Config()
	if err != nil {
		return name, email
	}
	userSec := cfg.Raw.Section("user")
	if userSec.Option("name") != "" {
		name = userSec.Option("name")
	}
	if userSec.Option("email") != "" {
		email = userSec.Option("email")
	}
	return name, email
}

func readUserFromGlobalGitConfig() (name, email string) {
	globalCfgPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	cfg, err := readConfig(globalCfgPath)
	if err != nil {
		return "", ""
	}
	readAndMergeIncludeIfConfig(cfg)
	userSec := cfg.Raw.Section("user")
	return userSec.Option("name"), userSec.Option("email")
}

func readAndMergeIncludeIfConfig(cfg *config.Config) {
	includeSec := cfg.Raw.Section("includeIf")
	if includeSec == nil {
		return
	}

	pwd, _ := os.Getwd()
	for _, sub := range includeSec.Subsections {
		if !strings.HasPrefix(sub.Name, "gitdir:") {
			continue
		}
		dir := strings.TrimLeft(sub.Name, "gitdir:")
		if !pathMatch(pwd, dir) {
			continue
		}
		path := sub.Option("path")
		if !filepath.IsAbs(path) {
			path = filepath.Join(os.Getenv("HOME"), path)
		}
		include, err := readConfig(path)
		if err == nil {
			cfg.Raw.Sections = append(cfg.Raw.Sections, include.Raw.Sections...)
		}
	}
}

func readConfig(path string) (*config.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rawGlobalCfg, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := config.NewConfig()
	err = cfg.Unmarshal(rawGlobalCfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func pathMatch(path string, exp string) bool {
	exp = strings.TrimRight(exp, "**")
	if strings.HasPrefix(exp, "**/") {
		exp = strings.TrimLeft(exp, "**/")
		return strings.Contains(strings.ToLower(path), strings.ToLower(exp))
	} else {
		return strings.HasPrefix(strings.ToLower(path), strings.ToLower(exp))
	}
}

func (gh *gitHelper) gitCommitPreCheck() error {
	wt, err := gh.repo.Worktree()
	if err != nil {
		return err
	}

	status, err := wt.Status()
	if err != nil {
		return err
	}
	if status.IsClean() {
		return fmt.Errorf("nothing to commit, working tree clean")
	}

	if isGitStatusNothingStaged(status) {
		return fmt.Errorf("nothing to commit, no changes added to commit")
	}

	return nil
}

func isGitStatusNothingStaged(s git.Status) bool {
	for _, status := range s {
		if status.Staging != git.Unmodified && status.Staging != git.Untracked {
			return false
		}
	}
	return true
}

func (gh *gitHelper) gitCommit(msg string) error {
	wt, err := gh.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Commit(msg, &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			Name:  gh.userName,
			Email: gh.userEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
