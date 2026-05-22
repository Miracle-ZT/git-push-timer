package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-push-timer/internal/config"
	"git-push-timer/internal/logger"
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

func TestHasUnpushedCommitsDetectsAheadBranch(t *testing.T) {
	repoDir := t.TempDir()
	remoteDir := t.TempDir()
	initGitRepo(t, repoDir)
	initBareGitRepo(t, remoteDir)
	configureGitUser(t, repoDir)
	addGitRemote(t, repoDir, remoteDir)
	commitFile(t, repoDir, "tracked.txt", "initial\n")
	pushBranch(t, repoDir, "master")

	exec := &Executor{}
	hasUnpushedCommits, err := exec.hasUnpushedCommits(repoDir, "master")
	if err != nil {
		t.Fatalf("hasUnpushedCommits() clean error = %v", err)
	}
	if hasUnpushedCommits {
		t.Fatal("hasUnpushedCommits() clean = true, want false")
	}

	commitFile(t, repoDir, "tracked.txt", "updated\n")

	hasUnpushedCommits, err = exec.hasUnpushedCommits(repoDir, "master")
	if err != nil {
		t.Fatalf("hasUnpushedCommits() ahead error = %v", err)
	}
	if !hasUnpushedCommits {
		t.Fatal("hasUnpushedCommits() ahead = false, want true")
	}
}

func TestHasUnpushedCommitsReturnsFalseForRepoWithoutCommits(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	exec := &Executor{}
	hasUnpushedCommits, err := exec.hasUnpushedCommits(repoDir, "master")
	if err != nil {
		t.Fatalf("hasUnpushedCommits() error = %v", err)
	}
	if hasUnpushedCommits {
		t.Fatal("hasUnpushedCommits() = true, want false for repo without commits")
	}
}

func TestExecuteRepositoryPushesCleanRepoWithUnpushedCommit(t *testing.T) {
	repoDir := t.TempDir()
	remoteDir := t.TempDir()
	initGitRepo(t, repoDir)
	initBareGitRepo(t, remoteDir)
	configureGitUser(t, repoDir)
	addGitRemote(t, repoDir, remoteDir)
	commitFile(t, repoDir, "tracked.txt", "initial\n")
	pushBranch(t, repoDir, "master")
	commitFile(t, repoDir, "tracked.txt", "updated\n")

	log, err := logger.New()
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}
	defer log.Close()

	exec := New(log)
	repo := config.Repository{
		Name:   "test-repo",
		Path:   repoDir,
		Branch: "master",
	}

	if err := exec.ExecuteRepository(repo); err != nil {
		t.Fatalf("ExecuteRepository() error = %v", err)
	}

	hasUnpushedCommits, err := exec.hasUnpushedCommits(repoDir, "master")
	if err != nil {
		t.Fatalf("hasUnpushedCommits() error = %v", err)
	}
	if hasUnpushedCommits {
		t.Fatal("hasUnpushedCommits() = true, want false after ExecuteRepository push")
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init", "-q", dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init error = %v, output = %s", err, string(output))
	}
}

func initBareGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init", "--bare", "-q", dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init --bare error = %v, output = %s", err, string(output))
	}
}

func configureGitUser(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "config", "user.name", "Git Push Timer Test")
	runGit(t, dir, "config", "user.email", "git-push-timer@example.invalid")
}

func addGitRemote(t *testing.T, dir, remote string) {
	t.Helper()

	runGit(t, dir, "remote", "add", "origin", remote)
}

func commitFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-q", "-m", "test commit")
}

func pushBranch(t *testing.T, dir, branch string) {
	t.Helper()

	runGit(t, dir, "push", "-q", "-u", "origin", branch)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s error = %v, output = %s", strings.Join(args, " "), err, string(output))
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
