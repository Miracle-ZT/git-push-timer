package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"git-push-timer/internal/config"
	"git-push-timer/internal/logger"
)

// Executor Git 执行器
type Executor struct {
	logger *logger.Logger
}

// New 创建执行器
func New(logger *logger.Logger) *Executor {
	return &Executor{
		logger: logger,
	}
}

// ExecuteRepository 执行单个仓库的 commit 和 push
func (e *Executor) ExecuteRepository(repo config.Repository) error {
	// 展开路径中的 ~
	path, err := config.ExpandTilde(repo.Path)
	if err != nil {
		e.logger.Error("路径展开失败：%s", repo.Path)
		return err
	}

	e.logger.Info("开始处理仓库：%s (%s)", repo.Name, path)

	// 检查目录是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e.logger.Error("目录不存在：%s", path)
		return fmt.Errorf("directory not found: %s", path)
	}

	// 检查是否是 Git 仓库
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		e.logger.Error("不是 Git 仓库：%s", repo.Path)
		return fmt.Errorf("not a git repository: %s", repo.Path)
	}

	// 检查是否有变更
	hasChanges, err := e.hasRepositoryChanges(path)
	if err != nil {
		e.logger.Error("检查仓库变更失败：%v", err)
		return err
	}

	if hasChanges {
		// 有变更，先执行 commit
		if err := e.runCommand(path, "git", "add", "."); err != nil {
			e.logger.Error("git add 失败：%v", err)
			return err
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		commitMsg := fmt.Sprintf("auto: %s", timestamp)
		if err := e.runCommand(path, "git", "commit", "-m", commitMsg); err != nil {
			e.logger.Error("git commit 失败：%v", err)
			return err
		}
	}

	hasUnpushedCommits, err := e.hasUnpushedCommits(path, repo.Branch)
	if err != nil {
		e.logger.Error("检查未推送提交失败：%v", err)
		return err
	}
	if !hasUnpushedCommits {
		if !hasChanges {
			e.logger.Info("仓库 %s 没有变更，跳过", repo.Name)
		} else {
			// 正常 commit 成功后应该存在未推送提交；这里主要覆盖分支配置不一致、git hook 或其他进程在检查前已完成推送等边界情况。
			e.logger.Info("仓库 %s 没有未推送提交，跳过推送", repo.Name)
		}
		return nil
	}

	// git push
	if err := e.runCommand(path, "git", "push", "origin", repo.Branch); err != nil {
		e.logger.Error("git push 失败：%v", err)
		return err
	}

	e.logger.Info("仓库 %s 推送成功", repo.Name)
	return nil
}

func (e *Executor) hasRepositoryChanges(path string) (bool, error) {
	output, err := e.runCommandOutput(path, "git", "status", "--porcelain=v1", "--untracked-files=normal")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) != "", nil
}

func (e *Executor) hasUnpushedCommits(path, branch string) (bool, error) {
	if err := e.runCommand(path, "git", "rev-parse", "--verify", "--quiet", "HEAD"); err != nil {
		return false, nil
	}

	remoteRef := fmt.Sprintf("origin/%s", branch)
	if err := e.runCommand(path, "git", "rev-parse", "--verify", "--quiet", remoteRef); err != nil {
		if err := e.runCommand(path, "git", "rev-parse", "--verify", "--quiet", branch); err != nil {
			return false, err
		}
		return true, nil
	}

	output, err := e.runCommandOutput(path, "git", "rev-list", "--count", fmt.Sprintf("origin/%s..%s", branch, branch))
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) != "0", nil
}

// runCommand 执行命令
func (e *Executor) runCommand(workDir, name string, args ...string) error {
	_, err := e.runCommandOutput(workDir, name, args...)
	return err
}

func (e *Executor) runCommandOutput(workDir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %v, output: %s", name, strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}
