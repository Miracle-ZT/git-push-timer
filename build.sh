#!/bin/bash

# Git Push Timer - 构建脚本

echo "=== Git Push Timer Build ==="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误：Go 未安装，请先安装 Go 1.21+"
    echo "下载地址：https://go.dev/dl/"
    exit 1
fi

echo "Go 版本：$(go version)"

# 下载依赖
echo "下载依赖..."
go mod download
go mod tidy

# 创建输出目录
mkdir -p dist

# macOS 编译
echo "编译 macOS 版本..."
GOOS=darwin GOARCH=amd64 go build -o dist/git-push-timer-darwin-amd64 ./cmd/git-push-timer
echo "  -> dist/git-push-timer-darwin-amd64"

# Windows 编译
echo "编译 Windows 版本..."
GOOS=windows GOARCH=amd64 go build -o dist/git-push-timer-windows-amd64.exe ./cmd/git-push-timer
echo "  -> dist/git-push-timer-windows-amd64.exe"

echo "=== 编译完成 ==="
echo ""
echo "使用方式："
echo "1. 在 dist/ 目录中找到对应平台的二进制文件"
echo "2. 在同级目录下创建 config/repos.json 配置文件"
echo "3. 运行二进制文件"
echo ""
echo "示例配置已生成在 config/repos.json.example"
