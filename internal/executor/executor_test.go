package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHasRepositoryChangesReturnsFalseForCleanRepoWithoutCommits(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	restore := chdir(t, repoDir)
	defer restore()

	exec := &Executor{}
	hasChanges, err := exec.hasRepositoryChanges()
	if err != nil {
		t.Fatalf("hasRepositoryChanges() error = %v", err)
	}

	if hasChanges {
		t.Fatal("hasRepositoryChanges() = true, want false for clean repo")
	}
}

func TestHasRepositoryChangesDetectsUntrackedFiles(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	if err := os.WriteFile(filepath.Join(repoDir, "new-file.txt"), []byte("hello\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	restore := chdir(t, repoDir)
	defer restore()

	exec := &Executor{}
	hasChanges, err := exec.hasRepositoryChanges()
	if err != nil {
		t.Fatalf("hasRepositoryChanges() error = %v", err)
	}

	if !hasChanges {
		t.Fatal("hasRepositoryChanges() = false, want true for untracked file")
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init", "-q", dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init error = %v, output = %s", err, string(output))
	}
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", dir, err)
	}

	return func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("restore Chdir(%q) error = %v", oldDir, err)
		}
	}
}
