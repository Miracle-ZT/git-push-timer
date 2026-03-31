package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger 日志记录器
type Logger struct {
	logFile *os.File
}

// New 创建日志记录器，日志输出到可执行文件所在目录的 logs 子目录
func New() (*Logger, error) {
	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	// 日志目录
	logDir := filepath.Join(execDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 日志文件（按日期命名）
	today := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", today))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{logFile: logFile}, nil
}

// Info 记录 info 级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Error 记录 error 级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

// Warn 记录 warn 级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log("WARN", format, args...)
}

// log 内部日志记录方法
func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s %s\n", timestamp, level, message)

	// 写入文件
	fmt.Fprint(l.logFile, logLine)
	l.logFile.Sync()

	// 同时输出到控制台
	fmt.Print(logLine)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	return l.logFile.Close()
}
