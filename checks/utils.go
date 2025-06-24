package checks

import (
	"os"
	"os/exec"
)

// CommandExists checks if a command exists in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// IsGitURL checks if the given path is a git URL
func IsGitURL(path string) bool {
	return len(path) > 4 && ( // quick length check
		// HTTPS, SSH, git, file, or ends with .git
		(strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "git://") || 
		 strings.HasPrefix(path, "git@") || strings.HasPrefix(path, "ssh://") || strings.HasPrefix(path, "file://")) ||
		strings.HasSuffix(path, ".git"))
}

// CloneGitRepo clones the git repo at url to a temp dir, returns the dir path and a cleanup func
func CloneGitRepo(url string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "check-git-")
	if err != nil {
		return "", nil, err
	}
	cmd := exec.Command("git", "clone", "--depth=1", url, tmpDir)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return tmpDir, cleanup, nil
} 