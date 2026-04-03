#!/bin/bash

# Git Push Timer - Release 打包脚本

set -e

echo "=== Git Push Timer Release ==="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误：Go 未安装"
    exit 1
fi

# 获取版本号（从 git tag）
VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null)}
if [ -z "$VERSION" ]; then
    echo "错误：未指定版本号，且当前不在 tag 上"
    echo "用法：./release.sh v1.0.0"
    echo "   或：git tag v1.0.0 && ./release.sh"
    exit 1
fi

echo "版本号：$VERSION"

# 下载依赖
echo "下载依赖..."
go mod download
go mod tidy

# 创建输出目录
rm -rf dist
mkdir -p dist

# 编译 macOS 版本
echo "编译 macOS 版本..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o dist/git-push-timer ./cmd/git-push-timer

# 编译 Windows 版本
echo "编译 Windows 版本..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o dist/git-push-timer.exe ./cmd/git-push-timer

# 打包
echo "打包..."
cd dist
zip "git-push-timer_${VERSION}_darwin_amd64.zip" git-push-timer
zip "git-push-timer_${VERSION}_windows_amd64.zip" git-push-timer.exe
cd ..

# 复制配置文件
cp config/repos.json.example dist/

echo ""
echo "=== 打包完成 ==="
echo "输出文件:"
ls -la dist/*.zip
echo ""
echo "发布步骤:"
echo "1. git tag $VERSION"
echo "2. git push origin $VERSION"
echo "3. 在 GitHub Releases 页面上传 dist/*.zip 文件"
