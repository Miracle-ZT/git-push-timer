package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHasRepositoryChangesReturnsFalseForCleanRepoWithoutCommits(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	otherDir := t.TempDir()
	restore := chdir(t, otherDir)
	defer restore()

	exec := &Executor{}
	hasChanges, err := exec.hasRepositoryChanges(repoDir)
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

	otherDir := t.TempDir()
	restore := chdir(t, otherDir)
	defer restore()

	exec := &Executor{}
	hasChanges, err := exec.hasRepositoryChanges(repoDir)
	if err != nil {
		t.Fatalf("hasRepositoryChanges() error = %v", err)
	}

	if !hasChanges {
		t.Fatal("hasRepositoryChanges() = false, want true for untracked file")
	}
}

func TestRunCommandOutputUsesSpecifiedWorkingDirectory(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	otherDir := t.TempDir()
	restore := chdir(t, otherDir)
	defer restore()

	exec := &Executor{}
	output, err := exec.runCommandOutput(repoDir, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatalf("runCommandOutput() error = %v", err)
	}

	gotPath, err := filepath.EvalSymlinks(strings.TrimSpace(output))
	if err != nil {
		t.Fatalf("EvalSymlinks(output) error = %v", err)
	}

	wantPath, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("EvalSymlinks(repoDir) error = %v", err)
	}

	if gotPath != wantPath {
		t.Fatalf("runCommandOutput() = %q, want %q", gotPath, wantPath)
	}

	if cwd, err := os.Getwd(); err != nil {
		t.Fatalf("Getwd() error = %v", err)
	} else {
		gotCWD, err := filepath.EvalSymlinks(cwd)
		if err != nil {
			t.Fatalf("EvalSymlinks(cwd) error = %v", err)
		}

		wantCWD, err := filepath.EvalSymlinks(otherDir)
		if err != nil {
			t.Fatalf("EvalSymlinks(otherDir) error = %v", err)
		}

		if gotCWD != wantCWD {
			t.Fatalf("Getwd() = %q, want %q", gotCWD, wantCWD)
		}
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
