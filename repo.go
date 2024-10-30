package gitdb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Repo struct {
	URL       string
	Repo      string
	Path      string
	Ref       string
	Branch    string
	LocalRef  string
	LocalPath string
	LocalDir  string
}

func (r *Repo) Init() error {
	now := time.Now()
	ts := now.Format("20060102T150405")
	gitParts := strings.Split(r.URL, "@")
	urlParts := strings.Split(gitParts[1], ":")
	r.Repo = fmt.Sprintf("%s/%s", urlParts[0], urlParts[1])

	dir, err := os.MkdirTemp("", "gitdb")
	if err != nil {
		return err
	}
	r.LocalRef = r.Ref
	r.LocalDir = dir
	r.LocalPath = fmt.Sprintf("%s/%s", dir, r.Path)
	cmd := exec.Command("git", "clone", "-b", r.Ref, r.URL, r.LocalDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	split := strings.Split(r.LocalPath, "/")
	folderPath := strings.Join(split[:len(split)-1], "/")
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}
	file, err := os.OpenFile(r.LocalPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if r.Branch != "" {
		localRef := fmt.Sprintf("%s_%s", r.Branch, ts)
		cmd := exec.Command("git", "checkout", "-b", localRef)
		cmd.Dir = r.LocalDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		r.LocalRef = localRef
	}
	return nil
}

func (r *Repo) Pull() error {
	cmd := exec.Command("git", "pull", "origin", r.Ref)
	cmd.Dir = r.LocalDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Repo) Push(message string) error {
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = r.LocalDir
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		return err
	}
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = r.LocalDir
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		return err
	}
	pushCmd := exec.Command("git", "push", "origin", r.LocalRef)
	pushCmd.Dir = r.LocalDir
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return err
	}
	return nil
}

func (r *Repo) Get() ([]byte, error) {
	return os.ReadFile(r.LocalPath)
}

func (r *Repo) Post(dat []byte) error {
	return os.WriteFile(r.LocalPath, []byte(dat), 0666)
}

func (r *Repo) PR(title, body string) error {
	cmd := exec.Command("gh", "pr", "create", "--repo", r.Repo, "--title", title, "--body", body, "--body", r.Ref)
	cmd.Dir = r.LocalDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
