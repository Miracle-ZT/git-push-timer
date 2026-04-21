package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger 日志记录器
type Logger struct {
	mu          sync.Mutex
	logDir      string
	logFile     *os.File
	currentDate string
	now         func() time.Time
	out         io.Writer
	errOut      io.Writer
	closed      bool
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

	return newWithOptions(logDir, time.Now, os.Stdout, os.Stderr)
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
	now := l.now()
	timestamp := now.Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s %s\n", timestamp, level, message)

	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.closed {
		if err := l.rotateIfNeeded(now); err != nil {
			fmt.Fprintf(l.errOut, "日志写入失败：%v\n", err)
		} else {
			if _, err := fmt.Fprint(l.logFile, logLine); err != nil {
				fmt.Fprintf(l.errOut, "日志写入失败：%v\n", err)
			} else if err := l.logFile.Sync(); err != nil {
				fmt.Fprintf(l.errOut, "日志写入失败：%v\n", err)
			}
		}
	}

	// 同时输出到控制台
	fmt.Fprint(l.out, logLine)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}
	l.closed = true

	if l.logFile == nil {
		return nil
	}

	err := l.logFile.Close()
	l.logFile = nil
	return err
}

func newWithOptions(logDir string, now func() time.Time, out, errOut io.Writer) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}

	logger := &Logger{
		logDir: logDir,
		now:    now,
		out:    out,
		errOut: errOut,
	}

	if err := logger.rotateIfNeeded(now()); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) rotateIfNeeded(now time.Time) error {
	date := now.Format("2006-01-02")
	if l.logFile != nil && l.currentDate == date {
		return nil
	}

	logPath := filepath.Join(l.logDir, fmt.Sprintf("%s.log", date))
	newFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	oldFile := l.logFile
	l.logFile = newFile
	l.currentDate = date

	if oldFile != nil {
		_ = oldFile.Close()
	}

	return nil
}
