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
	Path      string
	Ref       string
	Branch    string
	localPath string
	localDir  string
}

func (r *Repo) Init() error {
	now := time.Now()
	ts := now.Format("20060102T150405")

	dir, err := os.MkdirTemp("", "gitdb")
	if err != nil {
		return err
	}
	r.localDir = dir
	r.localPath = fmt.Sprintf("%s/%s", dir, r.Path)
	cmd := exec.Command("git", "clone", "-b", r.Ref, r.URL, r.localDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	split := strings.Split(r.localPath, "/")
	folderPath := strings.Join(split[:len(split)-1], "/")
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}
	file, err := os.OpenFile(r.localPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if r.Branch != "" {
		cmd := exec.Command("git", "checkout", "-b", fmt.Sprintf("%s_%s", r.Branch, ts))
		cmd.Dir = r.localDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) Update() error {
	cmd := exec.Command("git", "pull", "origin", r.Ref)
	cmd.Dir = r.localDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Repo) Get() ([]byte, error) {
	return os.ReadFile(r.localPath)
}

func (r *Repo) Post(dat []byte) error {
	return os.WriteFile(r.localPath, []byte(dat), 0666)
}

// func (r *Repo) PR(ref string) error {
// 	cmd := exec.Command("git", "pull", "origin", r.Ref)
// 	cmd.Dir = r.localDir
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }
