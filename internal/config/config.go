package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Repository 代表一个需要监控的 Git 仓库
type Repository struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Branch  string `json:"branch"`
	Enabled bool   `json:"enabled"`
}

// Config 配置文件结构
type Config struct {
	Repositories []Repository `json:"repositories"`
}

// Load 从可执行文件所在目录加载配置文件
func Load() (*Config, error) {
	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	// 配置文件路径
	configPath := filepath.Join(execDir, "config", "repos.json")

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
