package executor

import (
	"fmt"
	"os"
	"os/exec"
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

	// 切换到仓库目录
	oldDir, _ := os.Getwd()
	if err := os.Chdir(path); err != nil {
		e.logger.Error("无法进入目录：%s", path)
		return err
	}
	defer os.Chdir(oldDir)

	// 检查是否是 Git 仓库
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		e.logger.Error("不是 Git 仓库：%s", repo.Path)
		return fmt.Errorf("not a git repository: %s", repo.Path)
	}

	// 检查是否有变更
	hasChanges, err := e.hasRepositoryChanges()
	if err != nil {
		e.logger.Error("检查仓库变更失败：%v", err)
		return err
	}
	if !hasChanges {
		// 没有变更，跳过
		e.logger.Info("仓库 %s 没有变更，跳过", repo.Name)
		return nil
	}

	// 有变更，执行 commit 和 push
	// git add .
	if err := e.runCommand("git", "add", "."); err != nil {
		e.logger.Error("git add 失败：%v", err)
		return err
	}

	// git commit
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	commitMsg := fmt.Sprintf("auto: %s", timestamp)
	if err := e.runCommand("git", "commit", "-m", commitMsg); err != nil {
		e.logger.Error("git commit 失败：%v", err)
		return err
	}

	// git push
	if err := e.runCommand("git", "push", "origin", repo.Branch); err != nil {
		e.logger.Error("git push 失败：%v", err)
		return err
	}

	e.logger.Info("仓库 %s 推送成功", repo.Name)
	return nil
}

func (e *Executor) hasRepositoryChanges() (bool, error) {
	output, err := e.runCommandOutput("git", "status", "--porcelain=v1", "--untracked-files=normal")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) != "", nil
}

// runCommand 执行命令
func (e *Executor) runCommand(name string, args ...string) error {
	_, err := e.runCommandOutput(name, args...)
	return err
}

func (e *Executor) runCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %v, output: %s", name, strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}
