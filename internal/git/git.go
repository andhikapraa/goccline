// Package git is a thin wrapper over the git CLI.
// Each call shells out; we cache results in-process via the Context's lifecycle.
package git

import (
	"os/exec"
	"strconv"
	"strings"
)

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Branch returns the current branch (or empty if detached / not a repo).
func Branch(dir string) string {
	out, err := run(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil || out == "HEAD" {
		return ""
	}
	return out
}

// IsRepo reports whether dir is inside a git working tree.
func IsRepo(dir string) bool {
	_, err := run(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// CommitsToday returns the number of commits authored today by the
// configured user.email.
func CommitsToday(dir string) int {
	if !IsRepo(dir) {
		return 0
	}
	out, err := run(dir, "log", "--since=midnight", "--author=$(git config user.email)", "--oneline")
	if err != nil {
		return 0
	}
	if out == "" {
		return 0
	}
	return strings.Count(out, "\n") + 1
}

// SubmoduleCount returns the number of submodules in dir.
func SubmoduleCount(dir string) int {
	out, err := run(dir, "submodule", "status")
	if err != nil || out == "" {
		return 0
	}
	return strings.Count(out, "\n") + 1
}

// HasWorktree returns true if dir is a git worktree (not the main repo).
func HasWorktree(dir string) bool {
	out, err := run(dir, "rev-parse", "--git-dir")
	if err != nil {
		return false
	}
	return strings.Contains(out, "/worktrees/")
}

// MustInt parses an int or returns 0 (silent).
func MustInt(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
