// Package logger 日志管理。支持终端输出和按天拆分文件。

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	logDir string
)

func init() {
	// 默认日志目录
	logDir = filepath.Join("storage", "logs")
	Init(logDir)
}

// Init 初始化日志，可指定目录
func Init(dir string) {
	logDir = dir

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建日志目录失败: %v，回退到终端输出", err)
		initConsole()
		return
	}

	// 按天拆分日志文件
	date := time.Now().Format("2006-01-02")
	infoFile := filepath.Join(dir, fmt.Sprintf("info-%s.log", date))
	errorFile := filepath.Join(dir, fmt.Sprintf("error-%s.log", date))

	infoWriter, err := os.OpenFile(infoFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("打开 info 日志文件失败: %v，回退到终端输出", err)
		initConsole()
		return
	}

	errorWriter, err := os.OpenFile(errorFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("打开 error 日志文件失败: %v，回退到终端输出", err)
		initConsole()
		return
	}

	// Info、Warn → info 文件 + 终端
	Info = log.New(io.MultiWriter(os.Stdout, infoWriter), "[INFO]  ", log.LstdFlags)
	Warn = log.New(io.MultiWriter(os.Stdout, infoWriter), "[WARN]  ", log.LstdFlags)
	// Error → error 文件 + 终端
	Error = log.New(io.MultiWriter(os.Stderr, errorWriter), "[ERROR] ", log.LstdFlags|log.Lshortfile)
}

func initConsole() {
	Info = log.New(os.Stdout, "[INFO]  ", log.LstdFlags)
	Warn = log.New(os.Stdout, "[WARN]  ", log.LstdFlags)
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}

// Infof 记录常规日志，附带请求上下文
func Infof(c *gin.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Info.Printf("[%s %s] %s", c.Request.Method, c.Request.URL.Path, msg)
}

// Errorf 记录错误日志，附带请求上下文
func Errorf(c *gin.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Error.Printf("[%s %s] %s", c.Request.Method, c.Request.URL.Path, msg)
}
